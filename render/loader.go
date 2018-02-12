package render

import (
	"image"
	"os"
	"path/filepath"
	"sync"
)

var (
	// Form ./assets/images,
	// the default image directory.
	wd, _ = os.Getwd()
	dir   = filepath.Join(
		wd,
		"assets",
		"images")
	loadedImages = make(map[string]*image.RGBA)
	loadedSheets = make(map[string]*Sheet)
	// move to some batch load settings
	defaultPad = 0
	imageLock  = sync.RWMutex{}
	sheetLock  = sync.RWMutex{}
)

func subImage(rgba *image.RGBA, x, y, w, h int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			out.Set(i, j, rgba.At(x+i, y+j))
		}
	}
	return out
}

// SetAssetPaths sets the directories that files are loaded from when using
// the LoadSprite utility (and others). Oak will call this with SetupConfig.Assets
// joined with SetupConfig.Images after Init.
func SetAssetPaths(imagedir string) {
	dir = imagedir
}

// UnloadAll resets the cached set of loaded sprites and sheets to empty.
func UnloadAll() {
	imageLock.Lock()
	sheetLock.Lock()
	loadedImages = make(map[string]*image.RGBA)
	loadedSheets = make(map[string]*Sheet)
	sheetLock.Unlock()
	imageLock.Unlock()
}
