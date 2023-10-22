package files

import (
	"path/filepath"
)

type Entry struct {
	Type        Type
	RelatedPath string
	children    map[string]*Entry
}

func (e *Entry) IsDir() bool {
	return e.Type == TypeDirectory
}

func (e *Entry) IsLink() bool {
	return e.Type == TypeLink
}

func (e *Entry) IsRegular() bool {
	return e.Type == TypeRegular
}

func (e *Entry) AddFile(name, relatedPath string) (*Entry, error) {
	entry := &Entry{
		Type:        TypeRegular,
		RelatedPath: relatedPath,
	}
	return e.AddEntry(name, entry)
}

func (e *Entry) AddDirectory(name string) (*Entry, error) {
	entry := &Entry{
		Type: TypeDirectory,
	}
	return e.AddEntry(name, entry)
}

func (e *Entry) AddLink(name, relatedPath string) (*Entry, error) {
	entry := &Entry{
		Type:        TypeLink,
		RelatedPath: relatedPath,
	}
	return e.AddEntry(name, entry)
}

func (e *Entry) AddEntry(name string, entry *Entry) (*Entry, error) {
	if !e.IsDir() {
		return nil, ErrFSEntryNotDir
	}
	if ee, exists := e.children[name]; exists {
		return ee, ErrFSEntryExists
	}
	if e.children == nil {
		e.children = make(map[string]*Entry)
	}
	e.children[name] = entry
	return entry, nil
}

func (e *Entry) GetEntry(name string) (*Entry, error) {
	if !e.IsDir() {
		return nil, ErrFSEntryNotDir
	}
	entry, exists := e.children[name]
	if !exists {
		return nil, ErrFSEntryNotExists
	}
	return entry, nil
}

func (e *Entry) walk(base string, fn walkFunc) error {
	for name, entry := range e.children {
		path := filepath.Join(base, name)
		if err := fn(path, entry); err != nil {
			return err
		}
		if entry.IsDir() {
			if err := entry.walk(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
