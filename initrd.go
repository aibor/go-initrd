package initrd

import (
	"fmt"
	"io"
	"path/filepath"
	"slices"

	"github.com/aibor/go-initrd/internal/archive"
	"github.com/aibor/go-initrd/internal/files"
)

const (
	// LibsDir is the archive's directory for all dynamically linked libraries.
	LibsDir = "lib"
	// AdditionalFilesDir is the archive's directory for all additional files
	// beside the init file.
	FilesDir = "files"
	// LibSearchPath defines the directories to lookup linked libraries.
	LibSearchPath = "/lib:/lib64:/usr/lib:/usr/lib64"
)

// InitRD represents a file tree that can be used as an initrd for the Linux
// kernel.
//
// Create a new instance using [New]. Additional files can be added with
// [InitRD.AddFiles]. Dynamically linked ELF libraries can be resolved and added
// for all already added files by calling [InitRD.ResolveLinkedLibs]. Once
// ready, write the [InitRD] with [InitRD.WriteCPIO].
type InitRD struct {
	fileTree files.Tree
}

// New creates a new [InitRD] with the given file added as "/init".
func New(initFile string) *InitRD {
	fileTree := files.Tree{}
	_, _ = fileTree.GetRoot().AddFile("init", initFile)
	return &InitRD{fileTree}
}

// AddFiles creates [FilesDir] and adds the given files to it.
func (i *InitRD) AddFiles(files ...string) error {
	if len(files) == 0 {
		return nil
	}

	dirEntry, err := i.fileTree.Mkdir(FilesDir)
	if err != nil {
		return fmt.Errorf("add dir: %v", err)
	}

	for _, file := range files {
		if _, err := dirEntry.AddFile(filepath.Base(file), file); err != nil {
			return fmt.Errorf("add file %s: %v", file, err)
		}
	}

	return nil
}

// ResolveLinkedLibs recursively resolves the dynamically linked libraries of
// all regular files in the [InitRD].
//
// If the given searchPath string is empty the default [LibSearchPath] is used.
// Resolved libraries are added to [LibsDir]. For each search path a symoblic
// link is added pointiong to [LibsDir].
func (i *InitRD) ResolveLinkedLibs(searchPath string) error {
	if searchPath == "" {
		searchPath = LibSearchPath
	}
	searchPaths := filepath.SplitList(searchPath)
	searchPaths = slices.DeleteFunc(searchPaths, func(e string) bool { return e == "" })

	resolver := files.ELFLibResolver{
		SearchPaths: searchPaths,
	}

	err := i.fileTree.Walk(func(path string, entry *files.Entry) error {
		if entry.Type != files.TypeRegular {
			return nil
		}
		return resolver.Resolve(entry.RelatedPath)
	})
	if err != nil {
		return fmt.Errorf("resolve: %v", err)
	}

	if len(resolver.Libs) == 0 {
		return nil
	}

	dirEntry, err := i.fileTree.Mkdir(LibsDir)
	if err != nil {
		return fmt.Errorf("add libs dir: %v", err)
	}
	for _, lib := range resolver.Libs {
		name := filepath.Base(lib)
		if _, err := dirEntry.AddFile(name, lib); err != nil {
			return fmt.Errorf("add lib %s: %v", name, err)
		}
	}

	absLibDir := filepath.Join(string(filepath.Separator), LibsDir)
	for _, searchPath := range searchPaths {
		dir, name := filepath.Split(searchPath)
		dirEntry, err := i.fileTree.Mkdir(dir)
		if err != nil {
			return fmt.Errorf("get dir: %v", err)
		}
		if _, err := dirEntry.AddLink(name, absLibDir); err != nil {
			if err != files.ErrFSEntryExists {
				return fmt.Errorf("add lib link %s: %v", searchPath, err)
			}
		}
	}

	return nil
}

// WriteCPIO writes the [InitRD] as CPIO archive to the given writer.
func (i *InitRD) WriteCPIO(writer io.Writer) error {
	w := archive.NewCPIOWriter(writer)
	defer w.Close()
	return i.fileTree.Walk(func(p string, e *files.Entry) error {
		return writeEntry(w, p, e)
	})
}

func writeEntry(writer archive.Writer, path string, entry *files.Entry) error {
	switch entry.Type {
	case files.TypeRegular:
		return writer.WriteRegular(path, entry.RelatedPath, 0755)
	case files.TypeDirectory:
		return writer.WriteDirectory(path)
	case files.TypeLink:
		return writer.WriteLink(path, entry.RelatedPath)
	default:
		return fmt.Errorf("unknown file type %d", entry.Type)
	}
}
