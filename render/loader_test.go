package render

import (
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestBatchLoad(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	assert.Nil(t, BatchLoad(filepath.Join("assets", "images")))
	sh := GetSheet(filepath.Join("16", "jeremy.png"))
	assert.Equal(t, len(sh), 8)
	_, err := loadImage("dir", "dummy.jpg")
	assert.NotNil(t, err)
	sp := LoadSprite("dummy.gif")
	assert.Nil(t, sp)
	sp = LoadSprite(filepath.Join("16", "jeremy.png"))
	assert.NotNil(t, sp)
}
