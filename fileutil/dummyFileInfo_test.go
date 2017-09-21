package fileutil

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDummyFileInfo(t *testing.T) {
	dfi := dummyfileinfo{"file", false}
	assert.Equal(t, dfi.Name(), "file")
	assert.Equal(t, dfi.Size(), int64(0))
	assert.Equal(t, dfi.Mode(), os.ModeTemporary)
	assert.Equal(t, dfi.ModTime(), time.Time{})
	assert.Equal(t, dfi.IsDir(), false)
	assert.Equal(t, dfi.Sys(), nil)
}
