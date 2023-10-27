package initramfs

import (
	"path/filepath"
	"testing"

	"github.com/aibor/initramfs/internal/archive"
	"github.com/aibor/initramfs/internal/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveNew(t *testing.T) {
	archive := New("first")
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
	t.Run("unknown file type", func(t *testing.T) {
		i := Archive{}
		_, err := i.fileTree.GetRoot().AddEntry("init", &files.Entry{
			Type: files.Type(99),
		})
		require.NoError(t, err)
		err = i.writeTo(&archive.MockWriter{})
		assert.ErrorContains(t, err, "unknown file type 99")
	})

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
				Path:        "/init",
				RelatedPath: "input",
				Mode:        0755,
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
				i := Archive{}
				_, err := i.fileTree.GetRoot().AddEntry("init", &tt.entry)
				require.NoError(t, err)
				mock := archive.MockWriter{}
				err = i.writeTo(&mock)
				assert.NoError(t, err)
				assert.Equal(t, tt.mock, mock)
			})
			t.Run("fails", func(t *testing.T) {
				i := Archive{}
				_, err := i.fileTree.GetRoot().AddEntry("init", &tt.entry)
				require.NoError(t, err)
				mock := archive.MockWriter{Err: assert.AnError}
				err = i.writeTo(&mock)
				assert.Error(t, err, assert.AnError)
			})
		})
	}
}
