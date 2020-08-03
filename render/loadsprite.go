package render

import (
	"errors"
	"image"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/fileutil"
	"github.com/oakmound/oak/v2/oakerr"
)

func loadSprite(directory, fileName string) (*image.RGBA, error) {

	imageLock.Lock()
	defer imageLock.Unlock()

	if _, ok := loadedImages[fileName]; !ok {
		imgFile, err := fileutil.Open(filepath.Join(directory, fileName))
		if err != nil {
			return nil, err
		}
		defer func() {
			dlog.ErrorCheck(imgFile.Close())
		}()

		ext := filepath.Ext(fileName)
		decoder, ok := fileDecoders[ext]
		if !ok {
			return nil, errors.New("No decoder available for file type: " + ext)
		}
		img, err := decoder(imgFile)

		if err != nil {
			return nil, err
		}

		// Todo: we internally just use *image.RGBA, but that choice
		// of image encoding was arbitrary. If using the image.Image
		// interface would not hurt performance considerably, we should
		// just use that.
		//
		// This converts the
		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		for x := 0; x < bounds.Max.X; x++ {
			for y := 0; y < bounds.Max.Y; y++ {
				rgba.Set(x, y, color.RGBAModel.Convert(img.At(x, y)))
			}
		}

		loadedImages[fileName] = rgba

		dlog.Verb("Loaded filename: ", fileName)
	}
	return loadedImages[fileName], nil
}

// SpriteIsLoaded returns whether, when LoadSprite is called, a cached sheet will
// be used, or if false that a new file will attempt to be loaded and stored
func SpriteIsLoaded(fileName string) bool {
	imageLock.RLock()
	_, ok := loadedImages[fileName]
	imageLock.RUnlock()
	return ok
}

// GetSprite tries to find the given file in a private set of
// loaded sprites. If that file isn't cached, it will return an error.
func GetSprite(fileName string) (*Sprite, error) {
	imageLock.RLock()
	r, ok := loadedImages[fileName]
	imageLock.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: fileName}
	}
	return NewSprite(0, 0, r), nil
}

// LoadSprite will load the given file as an image by combining directory and fileName.
// The resulting image, if found, will be cached under fileName for
// later access through GetSprite. If the empty string is passed in for directory,
// the directory defined by oak.SetupConfig.Assets.Images will be used.
func LoadSprite(directory, fileName string) (*Sprite, error) {
	if directory == "" {
		directory = dir
	}
	r, err := loadSprite(directory, fileName)
	if err != nil {
		dlog.Error(err)
		return nil, err
	}
	return NewSprite(0, 0, r), nil
}
