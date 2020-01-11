package render

import (
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/v2/fileutil"
	"github.com/stretchr/testify/assert"
)

var (
	imgPath1    = filepath.Join("16", "jeremy.png")
	badImgPath1 = filepath.Join("16", "invalid.png")
)

func TestBatchLoad(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	assert.Nil(t, BatchLoad(filepath.Join("assets", "images")))
	sh, err := GetSheet(imgPath1)
	assert.Nil(t, err)
	assert.Equal(t, len(sh.ToSprites()), 8)
	_, err = loadSprite("dir", "dummy.jpg")
	assert.NotNil(t, err)
	sp, err := GetSprite("dummy.gif")
	assert.Nil(t, sp)
	assert.NotNil(t, err)
	sp, err = GetSprite(imgPath1)
	assert.NotNil(t, sp)
	assert.Nil(t, err)
	UnloadAll()
}

func TestSetAssetPath(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	_, err := LoadSheet(dir, imgPath1, 16, 16, 0)
	assert.Nil(t, err)
	UnloadAll()
	SetAssetPaths(wd)
	_, err = LoadSheet(dir, imgPath1, 16, 16, 0)
	assert.NotNil(t, err)
	UnloadAll()
	SetAssetPaths(
		filepath.Join(
			wd,
			"assets",
			"images"),
	)
	_, err = LoadSheet(dir, imgPath1, 16, 16, 0)
	assert.Nil(t, err)
	UnloadAll()

}

func TestBadSheetParams(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset

	_, err := LoadSheet(dir, imgPath1, 0, 16, 0)
	assert.NotNil(t, err)
	_, err = LoadSheet(dir, imgPath1, 16, 0, 0)
	assert.NotNil(t, err)
	_, err = LoadSheet(dir, imgPath1, 16, 16, -1)
	assert.NotNil(t, err)

	_, err = LoadSheet(dir, imgPath1, 16, 16, 1000)
	assert.NotNil(t, err)

}

func TestSheetStorage(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset

	assert.False(t, SheetIsLoaded(imgPath1))
	_, err := GetSheet(imgPath1)
	assert.NotNil(t, err)
	_, err = LoadSheet(dir, imgPath1, 16, 16, 0)
	assert.Nil(t, err)
	assert.True(t, SheetIsLoaded(imgPath1))
	_, err = GetSheet(imgPath1)
	assert.Nil(t, err)
	UnloadAll()
}

func TestSheetUtility(t *testing.T) {

	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	_, err := LoadSprites(dir, imgPath1, 16, 16, 0)
	assert.Nil(t, err)
	_, err = LoadSprites(dir, badImgPath1, 16, 16, 0)
	assert.NotNil(t, err)

	_, err = LoadSheetSequence(imgPath1, 16, 16, 0, 1, 0, 0)
	assert.Nil(t, err)
	_, err = LoadSheetSequence(badImgPath1, 16, 16, 0, 1, 0, 0)
	assert.NotNil(t, err)

}
