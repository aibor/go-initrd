package archive

import "io/fs"

type MockWriter struct {
	Path        string
	RelatedPath string
	Mode        fs.FileMode
	Err         error
}

func (m *MockWriter) WriteRegular(path, source string, mode fs.FileMode) error {
	m.Path = path
	m.RelatedPath = source
	m.Mode = mode
	return m.Err
}

func (m *MockWriter) WriteDirectory(path string) error {
	m.Path = path
	return m.Err
}

func (m *MockWriter) WriteLink(path, target string) error {
	m.Path = path
	m.RelatedPath = target
	return m.Err
}