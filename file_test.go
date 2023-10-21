package initrd_test

import (
	"testing"

	"github.com/aibor/go-initrd"
	"github.com/stretchr/testify/assert"
)

func TestFileSpecWriteTo(t *testing.T) {
	t.Run("unknown file type", func(t *testing.T) {
		fs := initrd.FileSpec{FileType: initrd.FileType(99)}
		err := fs.WriteTo(&initrd.MockWriter{})
		assert.ErrorContains(t, err, "unknown file type 99")
	})

	tests := []struct {
		name     string
		fileSpec initrd.FileSpec
		mock     initrd.MockWriter
		err      error
	}{
		{
			name: "regular",
			fileSpec: initrd.FileSpec{
				ArchivePath: "init",
				RelatedPath: "input",
				FileType:    initrd.FileTypeRegular,
				Mode:        0755,
			},
			mock: initrd.MockWriter{
				Path:        "/init",
				RelatedPath: "input",
				Mode:        0755,
			},
		},
		{
			name: "directory",
			fileSpec: initrd.FileSpec{
				ArchivePath: "lib",
				FileType:    initrd.FileTypeDirectory,
			},
			mock: initrd.MockWriter{
				Path: "/lib",
			},
		},
		{
			name: "link",
			fileSpec: initrd.FileSpec{
				ArchivePath: "usr/lib64",
				RelatedPath: "/lib",
				FileType:    initrd.FileTypeLink,
			},
			mock: initrd.MockWriter{
				Path:        "/usr/lib64",
				RelatedPath: "/lib",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("works", func(t *testing.T) {
				mock := initrd.MockWriter{}
				err := tt.fileSpec.WriteTo(&mock)
				assert.NoError(t, err)
				assert.Equal(t, tt.mock, mock)
			})
			t.Run("fails", func(t *testing.T) {
				mock := initrd.MockWriter{Err: assert.AnError}
				err := tt.fileSpec.WriteTo(&mock)
				assert.Error(t, err, assert.AnError)
			})
		})
	}
}
