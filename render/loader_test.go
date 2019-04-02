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

func TestSetAssetPath(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	_, err := LoadSheet(dir, filepath.Join("16", "jeremy.png"), 16, 16, 0)
	assert.Nil(t, err)
	UnloadAll()
	SetAssetPaths(wd)
	_, err = LoadSheet(dir, filepath.Join("16", "jeremy.png"), 16, 16, 0)
	assert.NotNil(t, err)
	UnloadAll()
	SetAssetPaths(
		filepath.Join(
			wd,
			"assets",
			"images"),
	)
	_, err = LoadSheet(dir, filepath.Join("16", "jeremy.png"), 16, 16, 0)
	assert.Nil(t, err)

}
