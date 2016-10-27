package render

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"

	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	fontdir string

	d        *font.Drawer
	static_d *font.Drawer
	f        *truetype.Font

	defaultHinting  font.Hinting
	defaultSize     float64
	defaultDPI      float64
	defaultColor    image.Image
	defaultFontFile string

	faceHinting font.Hinting
	faceSize    float64
	faceDPI     float64
	faceColor   image.Image
	loadedFonts = make(map[string]*truetype.Font)
)

func SetFontDefaults(wd, assetPath, fontPath, hinting, color, file string, size, dpi float64) {
	fontdir = filepath.Join(filepath.Dir(wd),
		assetPath,
		fontPath)
	parseFontHinting(hinting)
	defaultHinting = faceHinting
	faceSize = size
	defaultSize = faceSize
	faceDPI = dpi
	defaultDPI = faceDPI
	faceColor = parseFontColor(color)
	defaultColor = faceColor
	defaultFontFile = file
}

func InitFont(b screen.Buffer, static_b screen.Buffer) {
	LoadFont(defaultFontFile)
	f = loadedFonts[defaultFontFile]
	d = &font.Drawer{
		Dst: b.RGBA(),
		Src: faceColor,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    faceSize,
			DPI:     faceDPI,
			Hinting: faceHinting,
		}),
	}
	static_d = &font.Drawer{
		Dst: static_b.RGBA(),
		Src: faceColor,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    faceSize,
			DPI:     faceDPI,
			Hinting: faceHinting,
		}),
	}
}

func setFace() {
	d.Face = truetype.NewFace(f, &truetype.Options{
		Size:    faceSize,
		DPI:     faceDPI,
		Hinting: faceHinting,
	})
	static_d.Face = truetype.NewFace(f, &truetype.Options{
		Size:    faceSize,
		DPI:     faceDPI,
		Hinting: faceHinting,
	})
}

func SetFontColor(im image.Image) {
	d.Src = im
}
func SetFontSize(fontSize float64) {
	faceSize = fontSize
	setFace()
}
func SetFontDPI(dpi float64) {
	faceDPI = dpi
	setFace()
}
func SetFontHinting(hintType string) {
	parseFontHinting(hintType)
	setFace()
}

func parseFontHinting(hintType string) {
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
}

func parseFontColor(s string) image.Image {
	s = strings.ToLower(s)
	switch s {
	case "white":
		return image.White
	case "black":
		return image.Black
	default:
		return image.Black
	}
}

func ResetFontFormat() {
	faceHinting = defaultHinting
	faceSize = defaultSize
	faceDPI = defaultDPI
	setFace()
	d.Src = defaultColor
	static_d.Src = defaultColor
}

func LoadFont(fontFile string) {
	fontBytes, err := ioutil.ReadFile(filepath.Join(fontdir, fontFile))
	if err != nil {
		dlog.Error(err.Error())
		return
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		dlog.Error(err.Error())
		return
	}
	loadedFonts[fontFile] = font
}
