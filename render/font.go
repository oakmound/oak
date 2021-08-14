package render

import (
	"image"
	"image/draw"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

var (
	// DefFontGenerator is a default font generator, using an internally
	// compiled font colored white by default.
	DefFontGenerator = FontGenerator{
		Color:   image.White,
		RawFile: luxisrTTF,
	}
)

// A Font can create text renderables. It should be constructed from
// FontGenerator.Generate().
type Font struct {
	gen FontGenerator
	font.Drawer
	ttfnt  *truetype.Font
	bounds intgeom.Rect2
	Unsafe bool
	mutex  sync.Mutex

	Fallbacks []*Font
}

// A FontGenerator stores information that can be used to create a font
type FontGenerator struct {
	Cache   *Cache
	File    string
	RawFile []byte
	Color   image.Image
	// FontOptions holds all optional font components. Reasonable defaults
	// will be used if these are not provided.
	FontOptions
}

type FontOptions = truetype.Options

// DefaultFont returns a font built from DefFontGenerator.
func DefaultFont() *Font {
	fnt, _ := DefFontGenerator.Generate()
	return fnt
}

func (fg FontGenerator) validate() error {
	if len(fg.File) == 0 && len(fg.RawFile) == 0 {
		return oakerr.InvalidInput{InputName: "File"}
	}
	if fg.Color == nil {
		return oakerr.InvalidInput{InputName: "Color"}
	}
	return nil
}

// Generate generates a font. File or RawFile and Color must be provided.
// If Cache and File are provided, the generated font will be stored in the provided cache.
// If Cache is not provided, it will default to DefaultCache.
func (fg *FontGenerator) Generate() (*Font, error) {
	if err := fg.validate(); err != nil {
		return nil, err
	}
	if fg.Cache == nil {
		fg.Cache = DefaultCache
	}

	var fnt *truetype.Font
	var err error
	if len(fg.RawFile) != 0 {
		fnt, err = truetype.Parse(fg.RawFile)
		if err != nil {
			return nil, err
		}
	} else {
		fnt, err = fg.Cache.LoadFont(fg.File)
		if err != nil {
			return nil, err
		}
	}

	// This logic is copied from truetype for their face scaling
	size := 12.0
	if fg.FontOptions.Size != 0 {
		size = fg.FontOptions.Size
	}
	dpi := 12.0
	if fg.FontOptions.DPI != 0 {
		dpi = fg.FontOptions.DPI
	}
	scl := fixed.Int26_6(0.5 + (size * dpi * 64 / 72))
	bds := fnt.Bounds(scl)
	intBds := intgeom.NewRect2(
		bds.Min.X.Round(),
		bds.Min.Y.Round(),
		bds.Max.X.Round(),
		bds.Max.Y.Round(),
	)

	return &Font{
		gen: *fg,
		Drawer: font.Drawer{
			Src:  fg.Color,
			Face: truetype.NewFace(fnt, &fg.FontOptions),
		},
		ttfnt:  fnt,
		bounds: intBds,
	}, nil
}

// RegenerateWith creates a new font based on this font after changing its generation settings.
func (f *Font) RegenerateWith(fgFunc func(FontGenerator) FontGenerator) (*Font, error) {
	gen := fgFunc(f.gen)
	return gen.Generate()
}

// Copy returns a copy of this font
func (f *Font) Copy() *Font {
	if f.Unsafe {
		return f
	}
	f2 := &Font{
		gen:       f.gen,
		Drawer:    f.Drawer,
		ttfnt:     f.ttfnt,
		bounds:    f.bounds,
		Unsafe:    f.Unsafe,
		Fallbacks: f.Fallbacks,
	}
	f2.Drawer.Face = truetype.NewFace(f.ttfnt, &f.gen.FontOptions)
	return f2
}

// MeasureString calculates the width of a rendered text this font would draw from
// the given input string.
func (f *Font) MeasureString(s string) fixed.Int26_6 {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	prevC := rune(-1)
	var width fixed.Int26_6
	for _, c := range s {
		if prevC >= 0 {
			f.Drawer.Dot.X += f.Drawer.Face.Kern(prevC, c)
		}
		_, _, maskp, advance, ok := f.Drawer.Face.Glyph(f.Drawer.Dot, c)
		if _, empty := emptyboxYValues[maskp.Y]; !ok || empty {
			for _, fallback := range f.Fallbacks {
				_, _, maskp, advance, ok = fallback.Drawer.Face.Glyph(f.Drawer.Dot, c)
				if _, empty := emptyboxYValues[maskp.Y]; !empty && ok {
					break
				}
			}
			if _, empty := emptyboxYValues[maskp.Y]; !ok || empty {
				// TODO: is falling back on the U+FFFD glyph the responsibility of
				// the Drawer or the Face?
				// TODO: set prevC = '\ufffd'?
				continue
			}
		}
		width += advance
		prevC = c
	}
	return width
}

var (
	// In testing, these are the locations where Glyph will return it found a glyph,
	// but return an empty box.
	// TODO: more research--
	// 1. why do the fonts say these characters exist when they don't
	// 2. can we just say < 100 = undefined?
	emptyboxYValues = map[int]struct{}{
		0:  {},
		20: {},
		23: {},
		40: {},
		60: {},
		69: {},
		81: {},
		75: {},
		46: {},
		54: {},
		50: {},
		27: {},
		25: {},
	}
)

func (f *Font) drawString(s string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			f.Drawer.Dot.X += f.Drawer.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := f.Drawer.Face.Glyph(f.Drawer.Dot, c)
		if _, empty := emptyboxYValues[maskp.Y]; !ok || empty {
			for _, fallback := range f.Fallbacks {
				dr, mask, maskp, advance, ok = fallback.Drawer.Face.Glyph(f.Drawer.Dot, c)
				if _, empty := emptyboxYValues[maskp.Y]; !empty && ok {
					break
				}
			}
			if _, empty := emptyboxYValues[maskp.Y]; !ok || empty {
				// TODO: is falling back on the U+FFFD glyph the responsibility of
				// the Drawer or the Face?
				// TODO: set prevC = '\ufffd'?
				continue
			}
		}
		draw.DrawMask(f.Drawer.Dst, dr, f.Drawer.Src, image.Point{}, mask, maskp, draw.Over)
		f.Drawer.Dot.X += advance
		prevC = c
	}
}

// FontColor returns an image.Image color matching the SVG 1.1 spec.
// If the string does not align to a color in the spec, it will error.
func FontColor(s string) (image.Image, error) {
	if c, ok := colornames.Map[strings.ToLower(s)]; ok {
		return image.NewUniform(c), nil
	}
	return nil, oakerr.NotFound{InputName: "s"}
}

// GetFont returns a cached font, or an error if the font is not
// cached.
func (c *Cache) GetFont(file string) (*truetype.Font, error) {
	c.fontLock.RLock()
	f, ok := c.loadedFonts[file]
	c.fontLock.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: "file"}
	}
	return f, nil
}

// LoadFont loads the given font file, parses it, and caches it under
// its full path and its final path element.
func (c *Cache) LoadFont(file string) (*truetype.Font, error) {
	fontBytes, err := fileutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	c.fontLock.Lock()
	c.loadedFonts[file] = font
	c.loadedFonts[filepath.Base(file)] = font
	c.fontLock.Unlock()

	return font, nil
}
