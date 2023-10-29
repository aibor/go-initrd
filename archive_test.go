package initramfs

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/aibor/initramfs/internal/archive"
	"github.com/aibor/initramfs/internal/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveNew(t *testing.T) {
	archive := New("first")
	assert.Equal(t, os.DirFS("/"), archive.sourceFS)
	entry, err := archive.fileTree.GetEntry("/init")
	require.NoError(t, err)
	assert.Equal(t, "first", entry.RelatedPath)
	assert.Equal(t, files.TypeRegular, entry.Type)
}

func TestArchiveAddFile(t *testing.T) {
	archive := New("first")

	require.NoError(t, archive.AddFiles("second", "rel/third", "/abs/fourth"))
	require.NoError(t, archive.AddFiles("fifth"))
	require.NoError(t, archive.AddFiles())

	expected := map[string]string{
		"second": "second",
		"third":  "rel/third",
		"fourth": "/abs/fourth",
		"fifth":  "fifth",
	}

	for file, relPath := range expected {
		path := filepath.Join("files", file)
		e, err := archive.fileTree.GetEntry(path)
		require.NoError(t, err, path)
		assert.Equal(t, files.TypeRegular, e.Type)
		assert.Equal(t, relPath, e.RelatedPath)
	}
}

func TestArchiveWriteTo(t *testing.T) {
	testFS := fstest.MapFS{
		"input": &fstest.MapFile{},
	}
	testFile, err := testFS.Open("input")
	require.NoError(t, err)

	test := func(entry *files.Entry, w *archive.MockWriter) error {
		i := Archive{sourceFS: testFS}
		_, err := i.fileTree.GetRoot().AddEntry("init", entry)
		require.NoError(t, err)
		return i.writeTo(w)
	}

	t.Run("unknown file type", func(t *testing.T) {
		err := test(&files.Entry{Type: files.Type(99)}, &archive.MockWriter{})
		assert.ErrorContains(t, err, "unknown file type 99")
	})

	t.Run("nonexisting source", func(t *testing.T) {
		entry := &files.Entry{
			Type:        files.TypeRegular,
			RelatedPath: "nonexisting",
		}
		err := test(entry, &archive.MockWriter{})
		assert.ErrorContains(t, err, "open nonexisting: file does not exist")
	})

	t.Run("existing files", func(t *testing.T) {
		tests := []struct {
			name  string
			entry files.Entry
			mock  archive.MockWriter
		}{
			{
				name: "regular",
				entry: files.Entry{
					Type:        files.TypeRegular,
					RelatedPath: "input",
				},
				mock: archive.MockWriter{
					Path:   "/init",
					Source: testFile,
					Mode:   0755,
				},
			},
			{
				name: "directory",
				entry: files.Entry{
					Type: files.TypeDirectory,
				},
				mock: archive.MockWriter{
					Path: "/init",
				},
			},
			{
				name: "link",
				entry: files.Entry{
					Type:        files.TypeLink,
					RelatedPath: "/lib",
				},
				mock: archive.MockWriter{
					Path:        "/init",
					RelatedPath: "/lib",
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Run("works", func(t *testing.T) {
					i := Archive{sourceFS: testFS}
					_, err := i.fileTree.GetRoot().AddEntry("init", &tt.entry)
					require.NoError(t, err)
					mock := archive.MockWriter{}
					err = i.writeTo(&mock)
					assert.NoError(t, err)
					assert.Equal(t, tt.mock, mock)
				})
				t.Run("fails", func(t *testing.T) {
					i := Archive{sourceFS: testFS}
					_, err := i.fileTree.GetRoot().AddEntry("init", &tt.entry)
					require.NoError(t, err)
					mock := archive.MockWriter{Err: assert.AnError}
					err = i.writeTo(&mock)
					assert.Error(t, err, assert.AnError)
				})
			})
		}
	})
}