package files

import (
	"fmt"
	"path/filepath"
)

type Tree struct {
	root *Entry
}

func isRoot(path string) bool {
	switch path {
	case "", ".", string(filepath.Separator):
		return true
	default:
		return false
	}
}

func (t *Tree) GetRoot() *Entry {
	if t.root == nil {
		t.root = &Entry{
			Type: TypeDirectory,
		}
	}
	return t.root
}

func (t *Tree) GetEntry(path string) (*Entry, error) {
	if isRoot(path) {
		return t.GetRoot(), nil
	}
	dir, name := filepath.Split(filepath.Clean(path))
	parent, err := t.GetEntry(dir)
	if err != nil {
		return nil, err
	}
	return parent.GetEntry(name)
}

func (t *Tree) Mkdir(path string) (*Entry, error) {
	if isRoot(path) {
		return t.GetRoot(), nil
	}
	dir, name := filepath.Split(filepath.Clean(path))
	parent, err := t.Mkdir(dir)
	if err != nil {
		return nil, fmt.Errorf("mkdir %s: %v", path, err)
	}
	entry, err := parent.AddDirectory(name)
	if err == ErrFSEntryExists {
		err = nil
	}
	return entry, err
}

type walkFunc func(path string, entry *Entry) error

func (f *Tree) Walk(fn walkFunc) error {
	return f.GetRoot().walk(string(filepath.Separator), fn)
}
