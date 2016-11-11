package render

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"

	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	fontdir string

	defaultHinting  font.Hinting
	defaultSize     float64
	defaultDPI      float64
	defaultColor    image.Image
	defaultFontFile string

	DefFontGenerator = FontGenerator{}

	loadedFonts = make(map[string]*truetype.Font)
)

type FontGenerator struct {
	File    string
	Color   string
	Size    float64
	Hinting string
	DPI     float64
}

func DefFont() *Font {
	return DefFontGenerator.Generate()
}

func (fg *FontGenerator) Generate() *Font {

	// Replace zero values with defaults
	if fg.File == "" {
		fg.File = defaultFontFile
	}
	if fg.Size == 0 {
		fg.Size = defaultSize
	}
	if fg.DPI == 0 {
		fg.DPI = defaultDPI
	}

	return &Font{
		FontGenerator: *fg,
		Drawer: font.Drawer{
			// Color and hinting zero values are replaced
			// by their respective parse functions in the
			// zero case.
			Src: parseFontColor(fg.Color),
			Face: truetype.NewFace(LoadFont(fg.File), &truetype.Options{
				Size:    fg.Size,
				DPI:     fg.DPI,
				Hinting: parseFontHinting(fg.Hinting),
			}),
		},
	}

}

func (fg *FontGenerator) Copy() *FontGenerator {
	newFg := new(FontGenerator)
	*newFg = *fg
	return newFg
}

type Font struct {
	FontGenerator
	font.Drawer
}

func (f *Font) Copy() *Font {
	return f.Generate()
}

func (f *Font) Reset() {
	// Generate will return all defaults with no args
	f.FontGenerator = FontGenerator{}
	f = f.Generate()
}

func SetFontDefaults(wd, assetPath, fontPath, hinting, color, file string, size, dpi float64) {
	fontdir = filepath.Join(filepath.Dir(wd),
		assetPath,
		fontPath)
	defaultHinting = parseFontHinting(hinting)
	defaultSize = size
	defaultDPI = dpi
	defaultColor = parseFontColor(color)
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
		faceHinting = font.HintingNone
	}
	return faceHinting
}

func parseFontColor(s string) image.Image {
	s = strings.ToLower(s)
	switch s {
	case "white":
		return image.White
	case "black":
		return image.Black
	default:
		return defaultColor
	}
}

func LoadFont(fontFile string) *truetype.Font {
	if _, ok := loadedFonts[fontFile]; !ok {
		fontBytes, err := ioutil.ReadFile(filepath.Join(fontdir, fontFile))
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
