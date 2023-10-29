// Package initramfs provides a simple library for creating Linux initramfs
// archives.
//
// A simple program creating an initramfs archive and writing it to stdout can
// be found in "cmd/mkinitramfs".
//
// Only regular files are copied from the local file system. Mode is always set
// to 0755. For all added ELF file, the linked libraries can be resolved and
// added to the archive by calling [Archive.ResolveLinkedLibs].
package initramfs
