package initrd

import (
	"fmt"
	"path/filepath"
	"slices"
)

// LibDir is the archive's directory for all dynamically linked libraries.
const LibDir = "lib"

// Initrd is the collection of archive files.
type Initrd []FileSpec

// New creates a new [Initrd]. initFile is added as "/init". All additional
// files are added to the directory "files".
func New(initFile string, additionalFiles ...string) Initrd {
	initrd := Initrd{
		FileSpec{
			ArchivePath: "init",
			RelatedPath: initFile,
			FileType:    FileTypeRegular,
			Mode:        0755,
		},
		FileSpec{
			ArchivePath: "files",
			FileType:    FileTypeDirectory,
		},
	}
	for _, file := range additionalFiles {
		initrd = append(initrd, FileSpec{
			ArchivePath: absRootPath("files", filepath.Base(file)),
			RelatedPath: file,
			FileType:    FileTypeRegular,
			Mode:        0755,
		})
	}
	return initrd
}

// ResolveLinkedLibs resolves the dynamically linked libraries for each regular
// file in [Initrd]. The libraries are added to the archive's [LibDir]. Symlinks
// for the resolver's search paths are added, pointing to [LibDir].
func (i *Initrd) ResolveLinkedLibs(resolver *ELFLibResolver) error {
	for _, f := range *i {
		if f.FileType != FileTypeRegular {
			continue
		}
		if err := resolver.Resolve(f.RelatedPath); err != nil {
			return fmt.Errorf("resolve: %v", err)
		}
	}
	dirPaths := make([]string, 0)
	libs := resolver.Libs()
	newFiles := make([]FileSpec, 0)
	var addDir func(path string)
	addDir = func(dir string) {
		if dir == "." || dir == "/" || slices.Contains(dirPaths, dir) {
			return
		}
		addDir(filepath.Dir(dir))
		dirPaths = append(dirPaths, dir)
		newFiles = append(newFiles, FileSpec{
			ArchivePath: dir,
			FileType:    FileTypeDirectory,
		})
	}
	addDir(LibDir)
	for _, searchPath := range resolver.searchPaths {
		if searchPath == absRootPath(LibDir) {
			continue
		}

		addDir(filepath.Dir(searchPath))

		newFiles = append(newFiles, FileSpec{
			ArchivePath: searchPath,
			RelatedPath: absRootPath(LibDir),
			FileType:    FileTypeLink,
		})
	}
	for _, lib := range libs {
		newFiles = append(newFiles, FileSpec{
			ArchivePath: filepath.Join(LibDir, filepath.Base(lib)),
			RelatedPath: lib,
			FileType:    FileTypeRegular,
			Mode:        0755,
		})
	}
	*i = append(*i, newFiles...)

	return nil
}

// WriteTo writes the file to the given archive writer.
func (i *Initrd) WriteTo(w *Writer) error {
	for _, file := range *i {
		if err := file.WriteTo(w); err != nil {
			return err
		}
	}
	return nil
}
