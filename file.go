package initrd

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// FileType defines the type of a file.
type FileType int

const (
	// A regular file is copied completely into the archive.
	FileTypeRegular FileType = iota
	// A directory is created in the archive. Parent directories are not created
	// automatically. Ensure to create the complete file tree yourself.
	FileTypeDirectory
	// A symbolic link in the archive.
	FileTypeLink
)

// FileSpec defines an archive file.
type FileSpec struct {
	// File path in the initrd archive.
	ArchivePath string
	// Depending on the file type, the RelatedPath has different meanings.
	// Foe [FileTypeDirectory] it is not used at all. For [FileTypeRegular] it
	// is the path to the real file that should be copied into the archive. For
	// [FileTypeLink] it is the path of the link target.
	RelatedPath string
	// FileType defines of which type the file is.
	FileType FileType
	// Mode defines the permissions of the archive file.
	Mode fs.FileMode
}

// WriteTo calls the appropriate write method of the given [Writer] according
// to its [FileType].
func (s *FileSpec) WriteTo(w Writer) error {
	switch s.FileType {
	case FileTypeRegular:
		return w.WriteRegular(absRootPath(s.ArchivePath), s.RelatedPath, s.Mode)
	case FileTypeDirectory:
		return w.WriteDirectory(absRootPath(s.ArchivePath))
	case FileTypeLink:
		return w.WriteLink(absRootPath(s.ArchivePath), s.RelatedPath)
	default:
		return fmt.Errorf("unknown file type %d", s.FileType)
	}
}

// absRootPath wraps [filepath.Join] and prefixes the path with "/".
func absRootPath(parts ...string) string {
	sep := string(filepath.Separator)
	parts = append([]string{sep}, parts...)
	path := filepath.Join(parts...)
	if path == sep {
		return ""
	}
	return path
}
