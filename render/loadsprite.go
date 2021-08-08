package render

import (
	"errors"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

func loadSpriteNoCache(file string, maxFileSize int64) (*image.RGBA, error) {
	imgFile, err := fileutil.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	ext := filepath.Ext(file)
	ext = strings.ToLower(ext)

	if maxFileSize != 0 {
		cfgDecoder, ok := cfgDecoders[ext]
		if ok {
			// This can't reasonably error as we already loaded the file above
			info, _ := os.Lstat(file)
			if info.Size() > maxFileSize {
				// construct a blank image of the correct dimensions
				cfg, err := cfgDecoder(imgFile)
				if err != nil {
					return nil, err
				}
				rgba := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
				return rgba, nil
			}
		}
	}

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
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba.Set(x, y, color.RGBAModel.Convert(img.At(x, y)))
		}
	}

	return rgba, nil
}

func (c *Cache) loadSprite(file string, maxFileSize int64) (*image.RGBA, error) {
	rgba, err := loadSpriteNoCache(file, maxFileSize)
	if err != nil {
		return nil, err
	}
	c.imageLock.Lock()
	c.loadedImages[file] = rgba
	c.loadedImages[filepath.Base(file)] = rgba
	c.imageLock.Unlock()
	return rgba, nil
}

// GetSprite tries to find the given file in a private set of
// loaded sprites. If that file isn't cached, it will return an error.
func (c *Cache) GetSprite(file string) (*Sprite, error) {
	c.imageLock.RLock()
	r, ok := c.loadedImages[file]
	c.imageLock.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: file}
	}
	return NewSprite(0, 0, r), nil
}

// LoadSprite will load the given file as an image by combining directory and fileName.
// The resulting image, if found, will be cached under its last path element for
// later access through GetSprite.
func (c *Cache) LoadSprite(file string) (*Sprite, error) {
	r, err := c.loadSprite(file, 0)
	if err != nil {
		return nil, err
	}
	return NewSprite(0, 0, r), nil
}
