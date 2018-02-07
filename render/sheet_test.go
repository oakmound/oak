package render

import (
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestSheetSequence(t *testing.T) {

	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset

	_, err := NewSheetSequence(nil, 10, 0)
	assert.NotNil(t, err)

	sheet, err := LoadSheet(dir, filepath.Join("16", "jeremy.png"), 16, 16, 0)
	assert.Nil(t, err)
	sq, err := NewSheetSequence(sheet, 10, 0, 1, 0, 2)
	assert.Nil(t, err)
	assert.NotNil(t, sq)
}
