package initrd

import (
	"testing"

	"github.com/aibor/go-initrd/internal/archive"
	"github.com/aibor/go-initrd/internal/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRDNew(t *testing.T) {
	initRD := New("first")
	entry, err := initRD.fileTree.GetEntry("/init")
	require.NoError(t, err)
	assert.Equal(t, "first", entry.RelatedPath)
	assert.Equal(t, files.TypeRegular, entry.Type)
}

func TestWriteEntry(t *testing.T) {
	t.Run("unknown file type", func(t *testing.T) {
		entry := &files.Entry{
			Type: files.Type(99),
		}
		err := writeEntry(&archive.MockWriter{}, "init", entry)
		assert.ErrorContains(t, err, "unknown file type 99")
	})

	tests := []struct {
		name  string
		entry files.Entry
		mock  archive.MockWriter
		err   error
	}{
		{
			name: "regular",
			entry: files.Entry{
				Type:        files.TypeRegular,
				RelatedPath: "input",
			},
			mock: archive.MockWriter{
				Path:        "init",
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
				Path: "init",
			},
		},
		{
			name: "link",
			entry: files.Entry{
				Type:        files.TypeLink,
				RelatedPath: "/lib",
			},
			mock: archive.MockWriter{
				Path:        "init",
				RelatedPath: "/lib",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("works", func(t *testing.T) {
				mock := archive.MockWriter{}
				err := writeEntry(&mock, "init", &tt.entry)
				assert.NoError(t, err)
				assert.Equal(t, tt.mock, mock)
			})
			t.Run("fails", func(t *testing.T) {
				mock := archive.MockWriter{Err: assert.AnError}
				err := writeEntry(&mock, "init", &tt.entry)
				assert.Error(t, err, assert.AnError)
			})
		})
	}
}
