package render

import (
	"image"
	"path/filepath"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/oakerr"
)

// LoadSheet loads a file in some directory with sheets of (w,h) sized sprites.
// This will blow away any cached sheet with the same fileName.
func (c *Cache) LoadSheet(file string, cellSize intgeom.Point2) (*Sheet, error) {
	var rgba *image.RGBA
	var ok bool
	var err error

	if !ok {
		rgba, err = c.loadSprite(file, 0)
		if err != nil {
			return nil, err
		}
	}

	sheet, err := MakeSheet(rgba, cellSize)
	if err != nil {
		return nil, err
	}

	c.sheetLock.Lock()
	c.loadedSheets[file] = sheet
	c.loadedSheets[filepath.Base(file)] = sheet
	c.sheetLock.Unlock()

	return sheet, nil
}

// MakeSheet converts an image into a sheet with cellSize sized sprites
func MakeSheet(rgba *image.RGBA, cellSize intgeom.Point2) (*Sheet, error) {

	w := cellSize.X()
	h := cellSize.Y()

	if w <= 0 {
		return nil, oakerr.InvalidInput{InputName: "cellSize.X"}
	}
	if h <= 0 {
		return nil, oakerr.InvalidInput{InputName: "cellSize.Y"}
	}

	bounds := rgba.Bounds()

	sheetW := bounds.Max.X / w
	sheetH := bounds.Max.Y / h

	if sheetW < 1 || sheetH < 1 {
		return nil, oakerr.InvalidInput{InputName: "cellSize"}
	}

	if sheetW*w != sheetW {
		return nil, oakerr.InvalidInput{InputName: "nondivisibile dimensions x:"}
	}

	if sheetH*h != sheetH {
		return nil, oakerr.InvalidInput{InputName: "nondivisibile dimensions y:"}
	}

	sheet := make(Sheet, sheetW)
	i := 0
	for x := 0; x < bounds.Max.X; x += w {
		sheet[i] = make([]*image.RGBA, sheetH)
		j := 0
		for y := 0; y < bounds.Max.Y; y += h {
			sheet[i][j] = subImage(rgba, x, y, w, h)
			j++
		}
		i++
	}

	return &sheet, nil
}

// GetSheet tries to find the given file in the set of loaded sheets.
// If SheetIsLoaded(filename) is not true, this returns an error.
// Otherwise it will return the sheet as a 2d array of sprites
func (c *Cache) GetSheet(fileName string) (*Sheet, error) {
	c.sheetLock.RLock()
	sh, ok := c.loadedSheets[fileName]
	c.sheetLock.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: fileName}
	}
	return sh, nil
}

func subImage(rgba *image.RGBA, x, y, w, h int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			out.Set(i, j, rgba.At(x+i, y+j))
		}
	}
	return out
}
