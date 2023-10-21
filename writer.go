package initrd

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/cavaliergopher/cpio"
)

// Writer wraps cpio.Writer.
type Writer struct {
	cpioWriter *cpio.Writer
}

// NewWriter creates a new archive writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{cpio.NewWriter(w)}
}

// Close closes the [Writer]. Flush is called by the underlying closer.
func (w *Writer) Close() error {
	return w.cpioWriter.Close()
}

// Flush writes the data to the underlying [io.Writer].
func (w *Writer) Flush() error {
	return w.cpioWriter.Flush()
}

// writeHeader writes the cpio header.
func (w *Writer) writeHeader(hdr *cpio.Header) error {
	if err := w.cpioWriter.WriteHeader(hdr); err != nil {
		return fmt.Errorf("write header for %s: %v", hdr.Name, err)
	}
	return nil
}

// WriteDirectory add a directory entry for the given path to the archive.
func (w *Writer) WriteDirectory(path string) error {
	header := &cpio.Header{
		Name:  path,
		Mode:  cpio.TypeDir | cpio.ModePerm,
		Links: 2,
	}
	return w.writeHeader(header)
}

// WriteLink adds a symbolic link for the given path pointing to the given
// target.
func (w *Writer) WriteLink(path, target string) error {
	header := &cpio.Header{
		Name: path,
		Mode: cpio.TypeSymlink | cpio.ModePerm,
		Size: int64(len(target)),
	}
	if err := w.writeHeader(header); err != nil {
		return err
	}

	// Body of a link is the path of the target file.
	if _, err := w.cpioWriter.Write([]byte(target)); err != nil {
		return fmt.Errorf("write body for %s: %v", path, err)
	}

	return nil
}

// WriteRegular copies the exisiting file from source into the archive.
func (w *Writer) WriteRegular(path, source string, mode fs.FileMode) error {
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open file %s: %v", source, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("read info: %v", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", source)
	}

	cpioHdr, err := cpio.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("create header: %v", err)
	}

	cpioHdr.Name = path
	if mode != 0 {
		cpioHdr.Mode = cpio.FileMode(mode)
	}

	if err := w.writeHeader(cpioHdr); err != nil {
		return err
	}

	if _, err := io.Copy(w.cpioWriter, file); err != nil {
		return fmt.Errorf("write body for %s: %v", path, err)
	}

	return nil
}
