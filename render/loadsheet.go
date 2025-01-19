package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/oakerr"
	"golang.org/x/image/colornames"
)

// LoadSheet loads a file in some directory with sheets of (w,h) sized sprites.
// This will blow away any cached sheet with the same fileName.
func (c *Cache) LoadSheet(file string, cellSize intgeom.Point2) (*Sheet, error) {
	parseSheet := func(rgba *image.RGBA) (*Sheet, error) {
		return MakeSheet(rgba, cellSize)
	}
	return c.loadSheet(file, parseSheet)
}

// LoadSheetWithOptions loads a file in some directory with sheets of sprites.
// These sprites can be defined by the polgyons and offsets provided.
// This will blow away any cached sheet with the same fileName.
func (c *Cache) LoadSheetWithOptions(file string, opts ...SheetOption) (*Sheet, error) {
	parseSheet := func(rgba *image.RGBA) (*Sheet, error) {
		return MakeComplexSheet(rgba, opts...)
	}
	return c.loadSheet(file, parseSheet)
}

func (c *Cache) loadSheet(file string, parseSheet func(*image.RGBA) (*Sheet, error)) (*Sheet, error) {
	var rgba *image.RGBA
	var ok bool
	var err error

	if !ok {
		rgba, err = c.loadSprite(file, 0)
		if err != nil {
			return nil, err
		}
	}

	sheet, err := parseSheet(rgba)
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
// This is the legacy variant which fulfills most sprite sheet loads
// If the weightier format becomes more useful this may eventually be deprecated.
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
	wMod := bounds.Max.X % w
	if wMod != 0 {
		return nil, oakerr.InvalidInput{InputName: "cellSize.X"}
	}
	sheetH := bounds.Max.Y / h
	hMod := bounds.Max.Y % h
	if hMod != 0 {
		return nil, oakerr.InvalidInput{InputName: "cellSize.Y"}
	}
	if sheetW < 1 || sheetH < 1 {
		return nil, oakerr.InvalidInput{InputName: "cellSize"}
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

var defaultGenerator = SheetGenerator{
	IgnoreInputValidation: false,
	PerSpriteBuffer:       intgeom.Point2{0, 0},
	PerRowOffsets:         []intgeom.Point2{},
}

var emptyPoint intgeom.Point2

// MakeComplexSheet of sprites from a given image.
// These sheets or may not have a uniform rows and columns.
// This is a laxer enforcement than LoadSheet.
func MakeComplexSheet(rgba *image.RGBA, opts ...SheetOption) (*Sheet, error) {
	g := defaultGenerator
	for _, o := range opts {
		g = o(g)
	}

	cellBoundName := "Polygon"

	cellBounds := g.Bounds
	if g.SheetPolygon.IsEmpty() {
		if g.Bounds == emptyPoint {
			return nil, oakerr.InvalidInput{InputName: "Polygon/Bounds"}
		}
		cellBoundName = "Bounds"
		// g.SheetPolygon = NewPolygon2()
	} else {

		if g.Bounds == emptyPoint {
			floatBounds := g.SheetPolygon.Bounding.Max.Sub(g.SheetPolygon.Bounding.Min)
			cellBounds = intgeom.Point2{int(floatBounds.X()), int(floatBounds.Y())}

		}
	}
	fmt.Println("bounds", cellBounds)
	w := cellBounds.X()
	h := cellBounds.Y()

	// ingest the sheet
	bounds := rgba.Bounds()
	sheetW := bounds.Max.X / w
	sheetH := bounds.Max.Y / h
	if sheetW < 1 || sheetH < 1 {
		return nil, oakerr.InvalidInput{InputName: fmt.Sprintf("%s for image size", cellBoundName)}
	}

	sheet := Sheet{[]*image.RGBA{}}

	if g.CellPosition == nil {
		g.CellPosition = func(i, j int, dims intgeom.Point2) (topLeft intgeom.Point2) {
			return intgeom.Point2{
				i * dims.X(),
				j * dims.Y(),
			}
		}
	}

	i := 0
	j := 0
	for {
		topLeft := g.CellPosition(i, j, cellBounds)
		if topLeft.Y()+h > bounds.Max.Y {
			j = 0
			i++
			sheet = append(sheet, []*image.RGBA{})
			continue
		}
		if topLeft.X()+w > bounds.Max.X {
			break
		}

		candidiateImg := subImage(rgba, topLeft.X(), topLeft.Y(), w, h)
		rect := image.Rect(0, 0, w, h)
		if !g.SheetPolygon.IsEmpty() {
			poly := NewPolygon(g.SheetPolygon)
			poly.FillInverseOnRGBA(candidiateImg, color.RGBA{})

			draw.Draw(candidiateImg, rect, poly.GetRGBA(), image.Point{0, 0}, draw.Src)

			outline := poly.GetColoredOutline(IdentityColorer(colornames.Red), 1)
			if outline.GetRGBA() != nil {
				draw.Draw(candidiateImg, rect, outline.GetRGBA(), image.Point{0, 0}, draw.Src)
			}
		}

		sheet[i] = append(sheet[i], candidiateImg)
		j++
	}
	if len(sheet) != 0 && len(sheet[len(sheet)-1]) == 0 {
		sheet = sheet[:len(sheet)-1]
	}

	return &sheet, nil

}

type SheetGenerator struct {
	IgnoreInputValidation bool
	SheetPolygon          floatgeom.Polygon2
	Bounds                intgeom.Point2
	PerSpriteBuffer       intgeom.Point2
	PerRowOffsets         []intgeom.Point2
	CellPosition          func(i, j int, dims intgeom.Point2) (topLeft intgeom.Point2)
}

// generated via foptgen https://github.com/200sc/foptgen
type SheetOption func(SheetGenerator) SheetGenerator

func WithCellPosition(v func(i, j int, dims intgeom.Point2) intgeom.Point2) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.CellPosition = v
		return s
	}
}

func WithIgnoreInputValidation(v bool) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.IgnoreInputValidation = v
		return s
	}
}

func WithSheetPolygon(v floatgeom.Polygon2) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.SheetPolygon = v
		return s
	}
}

func WithBounds(v intgeom.Point2) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.Bounds = v
		return s
	}
}

func WithPerSpriteBuffer(v intgeom.Point2) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.PerSpriteBuffer = v
		return s
	}
}

func WithPerRowOffsets(v ...intgeom.Point2) SheetOption {
	return func(s SheetGenerator) SheetGenerator {
		s.PerRowOffsets = v
		return s
	}
}
