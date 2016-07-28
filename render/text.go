package render

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"

	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	fontdir string

	d *font.Drawer
	f *truetype.Font

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

type Text struct {
	Point
	Layered
	text string
	d_p  *font.Drawer
}

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

func InitFont(b_p *screen.Buffer) {
	b := *b_p
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
}

func setFace() {
	d.Face = truetype.NewFace(f, &truetype.Options{
		Size:    faceSize,
		DPI:     faceDPI,
		Hinting: faceHinting,
	})
}

func NewText(str string, x, y float64) *Text {
	return &Text{
		Point: Point{
			x,
			y,
		},
		text: str,
		d_p:  d,
	}
}

func (t_p *Text) GetRGBA() *image.RGBA {
	return nil
}

func (t_p *Text) Draw(buff screen.Buffer) {
	t_p.d_p.Dot = fixed.P(int(t_p.X), int(t_p.Y))
	t_p.d_p.DrawString(t_p.text)
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t_p *Text) Center() {
	textWidth := t_p.d_p.MeasureString(t_p.text).Round()
	t_p.ShiftX(float64(-textWidth / 2))
}

func DrawText(str string, x, y int) {
	d.Dot = fixed.P(x, y)
	d.DrawString(str)
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
