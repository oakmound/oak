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
	sh, err := GetSheet(filepath.Join("16", "jeremy.png"))
	assert.Nil(t, err)
	assert.Equal(t, len(sh.ToSprites()), 8)
	_, err = loadSprite("dir", "dummy.jpg")
	assert.NotNil(t, err)
	sp, err := GetSprite("dummy.gif")
	assert.Nil(t, sp)
	assert.NotNil(t, err)
	sp, err = GetSprite(filepath.Join("16", "jeremy.png"))
	assert.NotNil(t, sp)
	assert.Nil(t, err)
}
