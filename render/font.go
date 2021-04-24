package render

import (
	"image"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/fileutil"
	"github.com/oakmound/oak/v2/oakerr"
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

// DefaultFont returns a font built of the parameters set by SetFontDefaults.
func DefaultFont() *Font {
	fnt, _ := DefFontGenerator.Generate()
	return fnt
}

// Generate creates a font from the FontGenerator. Any parameters not supplied
// will be filled in with defaults set through SetFontDefaults.
func (fg *FontGenerator) Generate() (*Font, error) {

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
			return nil, err
		}
	} else {
		fnt = LoadFont(dir, fg.File)
		if fnt == nil {
			return nil, oakerr.InvalidInput{InputName: "fg.File"}
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
		ttfnt:  fnt,
		bounds: intBds,
	}, nil
}

// A Font is obtained as the result of FontGenerator.Generate(). It's used to
// create text type renderables.
type Font struct {
	FontGenerator
	font.Drawer
	ttfnt  *truetype.Font
	bounds intgeom.Rect2
	Unsafe bool
	mutex  sync.Mutex
}

// Copy returns a copy of this font
func (f *Font) Copy() *Font {
	if f.Unsafe {
		return f
	}
	f2 := &Font{}
	*f2 = *f
	f2.mutex = sync.Mutex{}
	f2.Drawer.Face = truetype.NewFace(f.ttfnt, &truetype.Options{
		Size:    f.FontGenerator.Size,
		DPI:     f.FontGenerator.DPI,
		Hinting: parseFontHinting(f.FontGenerator.Hinting),
	})
	return f2
}

func (f *Font) MeasureString(s string) fixed.Int26_6 {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.Drawer.MeasureString(s)
}

func (f *Font) DrawString(s string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.Drawer.DrawString(s)
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
