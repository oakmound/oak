package render

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/oakmound/oak/v3/oakerr"
)

// LoadSheet loads a file in some directory with sheets of (w,h) sized sprites,
// where there is pad pixels of vertical/horizontal empty space between each sprite.
// This will blow away any cached sheet with the same fileName.
func (c *Cache) LoadSheet(file string, cellW, cellH, padding int) (*Sheet, error) {
	var rgba *image.RGBA
	var ok bool
	var err error

	if !ok {
		rgba, err = c.loadSprite(file, 0)
		if err != nil {
			return nil, err
		}
	}

	sheet, err := MakeSheet(rgba, cellW, cellH, padding)
	if err != nil {
		return nil, err
	}

	c.sheetLock.Lock()
	c.loadedSheets[file] = sheet
	c.loadedSheets[filepath.Base(file)] = sheet
	c.sheetLock.Unlock()

	return sheet, nil
}

// MakeSheet converts an image into a sheet with (w,h) sized sprites,
// where there is pad pixels of vertical/horizontal empty space between each sprite.
func MakeSheet(rgba *image.RGBA, w, h, pad int) (*Sheet, error) {

	if w <= 0 {
		return nil, oakerr.InvalidInput{InputName: "w"}
	}
	if h <= 0 {
		return nil, oakerr.InvalidInput{InputName: "h"}
	}
	if pad < 0 {
		return nil, oakerr.InvalidInput{InputName: "pad"}
	}

	bounds := rgba.Bounds()

	sheetW := bounds.Max.X / w
	remainderW := bounds.Max.X % w
	sheetH := bounds.Max.Y / h
	remainderH := bounds.Max.Y % h

	var widthBuffers, heightBuffers int
	if pad != 0 {
		widthBuffers = remainderW / pad
		heightBuffers = remainderH / pad
	} else {
		widthBuffers = sheetW - 1
		heightBuffers = sheetH - 1
	}

	if sheetW < 1 || sheetH < 1 ||
		widthBuffers != sheetW-1 ||
		heightBuffers != sheetH-1 {
		return nil, oakerr.InvalidInput{InputName: "w,h"}
	}

	sheet := make(Sheet, sheetW)
	i := 0
	for x := 0; x < bounds.Max.X; x += (w + pad) {
		sheet[i] = make([]*image.RGBA, sheetH)
		j := 0
		for y := 0; y < bounds.Max.Y; y += (h + pad) {
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
	fmt.Println(c.loadedSheets)
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
