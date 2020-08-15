package render

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/fileutil"
)

var (
	fontdir string

	defaultHinting              = font.HintingNone
	defaultSize                 = 12.0
	defaultDPI                  = 72.0
	defaultColor    image.Image = image.White
	defaultFontFile string

	// DefFontGenerator is a default font generator of no options
	DefFontGenerator = FontGenerator{}

	loadedFonts = make(map[string]*truetype.Font)
)

// A FontGenerator stores information that can be used to create a font
type FontGenerator struct {
	File    string
	RawFile []byte
	Color   image.Image
	Size    float64
	Hinting string
	DPI     float64
}

func (fg FontGenerator) String() string {
	// Don't expose raw file content, it floods outputs
	type cleanFontGenerator struct {
		File    string
		Color   image.Image
		Size    float64
		Hinting string
		DPI     float64
	}
	clg := cleanFontGenerator{
		File:    fg.File,
		Size:    fg.Size,
		Hinting: fg.Hinting,
		DPI:     fg.DPI,
	}
	return fmt.Sprint(clg)
}

// DefFont returns a font built of the parameters set by SetFontDefaults.
func DefFont() *Font {
	return DefFontGenerator.Generate()
}

// Generate creates a font from the FontGenerator. Any parameters not supplied
// will be filled in with defaults set through SetFontDefaults.
func (fg *FontGenerator) Generate() *Font {

	dir := fontdir
	// Replace zero values with defaults
	var fnt *truetype.Font
	if fg.File == "" && len(fg.RawFile) == 0 {
		if defaultFontFile != "" {
			fg.File = defaultFontFile
		} else {
			fg.RawFile = luxisrTTF
		}
	}
	if len(fg.RawFile) != 0 {
		var err error
		fnt, err = truetype.Parse(fg.RawFile)
		if err != nil {
			// Todo: expose error here
			dlog.Error(err)
			return nil
		}
	} else {
		fnt = LoadFont(dir, fg.File)
		if fnt == nil {
			// Todo: expose error here
			return nil
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

	// This logic is copied from truetype for their face scaling
	scl := fixed.Int26_6(0.5 + (fg.Size * fg.DPI * 64 / 72))
	bds := fnt.Bounds(scl)
	intBds := intgeom.NewRect2(
		bds.Min.X.Round(),
		bds.Min.Y.Round(),
		bds.Max.X.Round(),
		bds.Max.Y.Round(),
	)

	return &Font{
		FontGenerator: *fg,
		Drawer: font.Drawer{
			// Color and hinting zero values are replaced
			// by their respective parse functions in the
			// zero case.
			Src: fg.Color,
			Face: truetype.NewFace(fnt, &truetype.Options{
				Size:    fg.Size,
				DPI:     fg.DPI,
				Hinting: parseFontHinting(fg.Hinting),
			}),
		},
		bounds: intBds,
	}

}

// Copy creates a copy of this FontGenerator
func (fg *FontGenerator) Copy() *FontGenerator {
	newFg := new(FontGenerator)
	*newFg = *fg
	return newFg
}

// A Font is obtained as the result of FontGenerator.Generate(). It's used to
// create text type renderables.
type Font struct {
	FontGenerator
	font.Drawer
	bounds intgeom.Rect2
}

// Refresh regenerates this font
func (f *Font) Refresh() {
	*f = *f.Generate()
}

// Copy returns a copy of this font
func (f *Font) Copy() *Font {
	return f.Generate()
}

// Reset sets the font to being a default font
func (f *Font) Reset() {
	// Generate will return all defaults with no args
	f.FontGenerator = FontGenerator{}
	*f = *f.Generate()
}

// SetFontDefaults updates the default font parameters with the passed in arguments
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
	default:
		dlog.Error("Unable to parse font hinting, ", hintType)
		fallthrough
	case "", "none":
		faceHinting = font.HintingNone
	case "vertical":
		faceHinting = font.HintingVertical
	case "full":
		faceHinting = font.HintingFull
	}
	return faceHinting
}

// FontColor accesses x/image/colornames and returns an image.Image for the input
// string. If the string is not defined in x/image/colornames, it will return defaultColor
// as defined by SetFontDefaults. The set of colors as defined by x/image/colornames matches
// the set of colors as defined by the SVG 1.1 spec.
func FontColor(s string) image.Image {
	s = strings.ToLower(s)
	if c, ok := colornames.Map[s]; ok {
		return image.NewUniform(c)
	}
	return defaultColor
}

// LoadFont loads in a font file and stores it with the given fontFile name.
// This is necessary before using that file in a generator, otherwise the default
// directory will be tried at generation time.
func LoadFont(dir string, fontFile string) *truetype.Font {
	if _, ok := loadedFonts[fontFile]; !ok {
		fontBytes, err := fileutil.ReadFile(filepath.Join(dir, fontFile))
		if err != nil {
			dlog.Error(err)
			return nil
		}
		font, err := truetype.Parse(fontBytes)
		if err != nil {
			dlog.Error(err)
			return nil
		}
		loadedFonts[fontFile] = font
	}
	return loadedFonts[fontFile]
}
