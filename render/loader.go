package render

import (
	"errors"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/oakmound/oak/oakerr"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/fileutil"
)

var (
	regexpSingleNumber, _ = regexp.Compile(`^\d+$`)
	regexpTwoNumbers, _   = regexp.Compile(`^\d+x\d+$`)
)

var (
	// Form ...main/core.go/assets/images,
	// the image directory.
	wd, _ = os.Getwd()
	dir   = filepath.Join(
		wd,
		"assets",
		"images")
	loadedImages = make(map[string]*image.RGBA)
	loadedSheets = make(map[string]*Sheet)
	// move to some batch load settings
	defaultPad = 0
	loadLock   = sync.Mutex{}
)

func loadImage(directory, fileName string) (*image.RGBA, error) {

	loadLock.Lock()
	if _, ok := loadedImages[fileName]; !ok {
		imgFile, err := fileutil.Open(filepath.Join(directory, fileName))
		if err != nil {
			loadLock.Unlock()
			return nil, err
		}
		defer func() {
			dlog.ErrorCheck(imgFile.Close())
		}()

		ext := strings.ToLower(fileName[len(fileName)-4:])
		decoder, ok := fileDecoders[ext]
		if !ok {
			return nil, errors.New("No decoder available for file type: " + ext)
		}
		img, err := decoder(imgFile)

		if err != nil {
			loadLock.Unlock()
			return nil, err
		}

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
	r := loadedImages[fileName]
	loadLock.Unlock()
	return r, nil
}

// LoadSprite loads the input fileName into a Sprite
func LoadSprite(fileName string) (*Sprite, error) {
	r, err := loadImage(dir, fileName)
	if err != nil {
		dlog.Error(err)
		return nil, err
	}
	return NewSprite(0, 0, r), nil
}

// GetSheet tries to find the given file in the set of loaded sheets.
// If SheetIsLoaded(filename) is not true, this returns an error.
// Otherwise it will return the sheet as a 2d array of sprites
func GetSheet(fileName string) ([][]*Sprite, error) {
	sprites := make([][]*Sprite, 0)
	dlog.Verb(loadedSheets, fileName, loadedSheets[fileName])
	if !SheetIsLoaded(fileName) {
		return sprites, oakerr.NotFound{InputName: fileName}
	}
	sheet, err := LoadSheet(dir, fileName, 0, 0, 0)
	if err != nil {
		return sprites, err
	}
	for x, row := range *sheet {
		sprites = append(sprites, make([]*Sprite, 0))
		for y := range row {
			sprites[x] = append(sprites[x], sheet.SubSprite(x, y))
		}
	}
	return sprites, nil
}

// SheetIsLoaded returns whether when LoadSheet is called, a cached sheet will
// be used, or if false that a new file will attempt to be loaded and stored
func SheetIsLoaded(filename string) bool {
	_, ok := loadedSheets[filename]
	return ok
}

// LoadSheet loads a file in some directory with sheets of (w,h) sized sprites,
// where there is pad pixels of vertical/horizontal pad between each sprite
func LoadSheet(directory, fileName string, w, h, pad int) (*Sheet, error) {
	if _, ok := loadedImages[fileName]; !ok {
		dlog.Verb("Missing file in loaded images: ", fileName)
		r, err := loadImage(directory, fileName)
		if err != nil {
			return nil, err
		}
		loadedImages[fileName] = r
	}
	if sheetP, ok := loadedSheets[fileName]; ok {
		return sheetP, nil
	}
	dlog.Verb("Loading sheet: ", fileName)
	rgba := loadedImages[fileName]
	bounds := rgba.Bounds()

	if w <= 0 {
		return nil, oakerr.InvalidInput{InputName: "w"}
	}
	if h <= 0 {
		return nil, oakerr.InvalidInput{InputName: "h"}
	}

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
		dlog.Error("Bad dimensions given to load sheet")
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

	dlog.Verb("Loaded sheet into map")
	loadedSheets[fileName] = &sheet

	return loadedSheets[fileName], nil
}

// LoadSheetSequence loads a sheet and then calls LoadSequence on that sheet
func LoadSheetSequence(fileName string, w, h, pad int, fps float64, frames ...int) (*Sequence, error) {
	sheet, err := LoadSheet(dir, fileName, w, h, pad)
	if err != nil {
		return nil, err
	}
	return LoadSequence(sheet, w, h, pad, fps, frames...)
}

// LoadSequence takes in a sheet with sheet dimensions, a frame rate and a list of frames where
// frames are in x,y pairs ([0,0,1,0,2,0] for (0,0) (1,0) (2,0)) and returns an animation from that
func LoadSequence(sheet *Sheet, w, h, pad int, fps float64, frames ...int) (*Sequence, error) {
	animation, err := NewSheetSequence(sheet, fps, frames...)
	if err != nil {
		return nil, err
	}
	return animation, nil
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

// SetAssetPaths sets the directories that files are loaded from when using
// the LoadSprite utility (and others). Oak will call this with SetupConfig.Assets
// joined with SetupConfig.Images after Init.
func SetAssetPaths(imagedir string) {
	dir = imagedir
}
