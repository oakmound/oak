package render

import (
	"image"
	"image/color"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/fileutil"
)

var (
	fontdir string

	defaultHinting  font.Hinting
	defaultSize     float64
	defaultDPI      float64
	defaultColor    image.Image
	defaultFontFile string

	// DefFontGenerator is a default font generator of no options
	DefFontGenerator = FontGenerator{}

	loadedFonts = make(map[string]*truetype.Font)
)

//A FontGenerator  stores a set of information that can be used to create a font
type FontGenerator struct {
	File    string
	Color   image.Image
	Size    float64
	Hinting string
	DPI     float64
}

//DefFont returns the default font
func DefFont() *Font {
	return DefFontGenerator.Generate()
}

//Generate creates a font from the FontGenerator
func (fg *FontGenerator) Generate() *Font {

	dir := fontdir
	// Replace zero values with defaults
	if fg.File == "" {
		if defaultFontFile != "" {
			fg.File = defaultFontFile
		} else {
			_, curFile, _, _ := runtime.Caller(1)
			dir = filepath.Join(filepath.Dir(curFile), "default_assets", "font")
			fg.File = "luxisr.ttf"
		}
	}
	if fg.Size == 0 {
		fg.Size = defaultSize
	}
	if fg.DPI == 0 {
		fg.DPI = defaultDPI
	}
	if fg.Color == nil {
		fg.Color = defaultColor
	}

	return &Font{
		FontGenerator: *fg,
		Drawer: font.Drawer{
			// Color and hinting zero values are replaced
			// by their respective parse functions in the
			// zero case.
			Src: fg.Color,
			Face: truetype.NewFace(LoadFont(dir, fg.File), &truetype.Options{
				Size:    fg.Size,
				DPI:     fg.DPI,
				Hinting: parseFontHinting(fg.Hinting),
			}),
		},
	}

}

//Copy cretaes a copy of the FontGenerator
func (fg *FontGenerator) Copy() *FontGenerator {
	newFg := new(FontGenerator)
	*newFg = *fg
	return newFg
}

//A Font can both be generated and drawn
type Font struct {
	FontGenerator
	font.Drawer
}

//Refresh regenerates the font per its generator's parameters
func (f *Font) Refresh() {
	*f = *f.Generate()
}

//Copy returns a new font from the given fonts generate function
func (f *Font) Copy() *Font {
	return f.Generate()
}

//Reset sets the font to being a default font
func (f *Font) Reset() {
	// Generate will return all defaults with no args
	f.FontGenerator = FontGenerator{}
	*f = *f.Generate()
}

//SetFontDefaults updates the default font parameters with the passed in arguments
func SetFontDefaults(wd, assetPath, fontPath, hinting, color, file string, size, dpi float64) {
	fontdir = filepath.Join(
		wd,
		assetPath,
		fontPath)
	defaultHinting = parseFontHinting(hinting)
	defaultSize = size
	defaultDPI = dpi
	defaultColor = FontColor(color)
	defaultFontFile = file
}

func parseFontHinting(hintType string) (faceHinting font.Hinting) {
	hintType = strings.ToLower(hintType)
	switch hintType {
	case "none":
		faceHinting = font.HintingNone
	case "vertical":
		faceHinting = font.HintingVertical
	case "full":
		faceHinting = font.HintingFull
	default:
		dlog.Error("Unable to parse font hinting, ", hintType)
		fallthrough
	case "":
		// Don't warn about undefined hinting
		faceHinting = font.HintingNone
	}
	return faceHinting
}

//FontColor converts a small set of strings to colors
//TODO: Implement a better version or pull in an outside library already doing this as this should be a fairly common utility function
func FontColor(s string) image.Image {
	s = strings.ToLower(s)
	switch s {
	case "white":
		return image.White
	case "black":
		return image.Black
	case "green":
		return image.NewUniform(color.RGBA{0, 255, 0, 255})
	default:
		return defaultColor
	}
}

//LoadFont loads in a font file and stores it with the given name. This is necessary before using the fonttype for a Font
func LoadFont(dir string, fontFile string) *truetype.Font {
	if _, ok := loadedFonts[fontFile]; !ok {
		fontBytes, err := fileutil.ReadFile(filepath.Join(dir, fontFile))
		if err != nil {
			dlog.Error(err.Error())
			return nil
		}
		font, err := truetype.Parse(fontBytes)
		if err != nil {
			dlog.Error(err.Error())
			return nil
		}
		loadedFonts[fontFile] = font
	}
	return loadedFonts[fontFile]
}
