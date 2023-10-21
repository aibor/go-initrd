package initrd_test

import (
	"path/filepath"
	"testing"

	"github.com/aibor/go-initrd"
	"github.com/stretchr/testify/assert"
)

func TestInitrdNew(t *testing.T) {
	i := initrd.New("first", "second", "rel/third", "/abs/fourth")
	expected := initrd.Initrd{
		initrd.FileSpec{
			ArchivePath: "init",
			RelatedPath: "first",
			FileType:    initrd.FileTypeRegular,
			Mode:        0755,
		},
		initrd.FileSpec{
			ArchivePath: initrd.AdditionalFilesDir,
			FileType:    initrd.FileTypeDirectory,
		},
		initrd.FileSpec{
			ArchivePath: filepath.Join(initrd.AdditionalFilesDir, "second"),
			RelatedPath: "second",
			FileType:    initrd.FileTypeRegular,
			Mode:        0755,
		},
		initrd.FileSpec{
			ArchivePath: filepath.Join(initrd.AdditionalFilesDir, "third"),
			RelatedPath: "rel/third",
			FileType:    initrd.FileTypeRegular,
			Mode:        0755,
		},
		initrd.FileSpec{
			ArchivePath: filepath.Join(initrd.AdditionalFilesDir, "fourth"),
			RelatedPath: "/abs/fourth",
			FileType:    initrd.FileTypeRegular,
			Mode:        0755,
		},
	}
	assert.Equal(t, expected, i)
}

func TestInitrdWriteTo(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		i := initrd.Initrd{
			initrd.FileSpec{
				FileType: initrd.FileTypeRegular,
			},
		}
		mock := initrd.MockWriter{}
		err := i.WriteTo(&mock)
		assert.NoError(t, err)
		assert.Equal(t, initrd.MockWriter{}, mock)
	})
	t.Run("fails", func(t *testing.T) {
		i := initrd.Initrd{
			initrd.FileSpec{
				FileType: initrd.FileTypeRegular,
			},
		}
		mock := initrd.MockWriter{Err: assert.AnError}
		err := i.WriteTo(&mock)
		assert.Error(t, err, assert.AnError)
	})

}
