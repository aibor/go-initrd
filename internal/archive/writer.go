package archive

import "io/fs"

// Writer defines initramfs archive writer interface.
type Writer interface {
	WriteRegular(string, string, fs.FileMode) error
	WriteDirectory(string) error
	WriteLink(string, string) error
}
