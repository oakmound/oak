package render

import (
	"errors"
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/oakmound/oak/v4/oakerr"
	"golang.org/x/image/colornames"
)

func TestFont_UnsafeCopy(t *testing.T) {
	f := DefaultFont()
	f.Unsafe = true
	f2 := f.Copy()
	if f2 != f {
		t.Fatalf("unsafe should have prevented the copy from actually copying")
	}
}

func TestFontColorSuccess(t *testing.T) {
	_, err := FontColor(colornames.Names[rand.Intn(len(colornames.Names)-1)])
	if err != nil {
		t.Fatalf("failed to load color: %v", err)
	}
}

func TestFontColorNotFound(t *testing.T) {
	_, err := FontColor("notacolor")
	expected := &oakerr.NotFound{}
	if !errors.As(err, expected) {
		t.Fatalf("expected not found error, got: %v", err)
	}
}

func TestFontGenerator_validate(t *testing.T) {
	fg := FontGenerator{}
	_, err := fg.Generate()
	expected := &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "File" {
		t.Fatalf("expected invalid File, got %v", expected.InputName)
	}

	fg = FontGenerator{File: "filename"}
	_, err = fg.Generate()
	expected = &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "Color" {
		t.Fatalf("expected invalid Color, got %v", expected.InputName)
	}
}

func TestFontGenerator_Generate_Failure(t *testing.T) {
	t.Run("BadRawFile", func(t *testing.T) {
		fg := FontGenerator{
			RawFile: []byte("notafontfile"),
			Color:   image.NewUniform(color.RGBA{255, 0, 0, 255}),
			FontOptions: FontOptions{
				Size: 13.0,
				DPI:  44.0,
			},
		}
		_, err := fg.Generate()
		if err == nil {
			t.Fatalf("generate should have failed")
		}
	})
	t.Run("BadLoadFont", func(t *testing.T) {
		fg := FontGenerator{
			File:  "file that does not exist",
			Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
			FontOptions: FontOptions{
				Size: 13.0,
				DPI:  44.0,
			},
		}
		_, err := fg.Generate()
		if err == nil {
			t.Fatalf("generate should have failed")
		}
	})
}

func TestFontGenerator_Generate_Success(t *testing.T) {
	fg := FontGenerator{
		File:  "testdata/assets/fonts/luxisr.ttf",
		Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
		FontOptions: FontOptions{
			Size: 13.0,
			DPI:  44.0,
		},
	}
	_, err := fg.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
}

func TestFont_Height(t *testing.T) {
	ht := rand.Float64() * 10
	fg := FontGenerator{
		File:  "testdata/assets/fonts/luxisr.ttf",
		Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
		FontOptions: FontOptions{
			Size: ht,
			DPI:  44.0,
		},
	}
	f, err := fg.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if f.Height() != ht {
		t.Fatalf("size did not match height: got %v expected %v", f.Height(), ht)
	}
}

func TestFont_RegenerateWith(t *testing.T) {
	fg := FontGenerator{
		File:  "testdata/assets/fonts/luxisr.ttf",
		Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
		FontOptions: FontOptions{
			Size: 13.0,
			DPI:  44.0,
		},
	}
	f, err := fg.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	f2, err := f.RegenerateWith(func(fg FontGenerator) FontGenerator {
		fg.Size = 100
		return fg
	})
	if err != nil {
		t.Fatalf("regenerate failed: %v", err)
	}
	if f2.Height() != 100 {
		t.Fatalf("size did not match height: got %v expected %v", f.Height(), 100)
	}
}

func TestCache_LoadFont(t *testing.T) {
	t.Run("NotExists", func(t *testing.T) {
		c := NewCache()
		_, err := c.LoadFont("bogusfilepath")
		if err == nil {
			t.Fatal("expected error loading bad file")
		}
	})
	t.Run("NotFontFile", func(t *testing.T) {
		c := NewCache()
		_, err := c.LoadFont("testdata/assets/images/devfile.pdn")
		if err == nil {
			t.Fatal("expected error loading non-font")
		}
	})
	t.Run("GetCached", func(t *testing.T) {
		c := NewCache()
		_, err := c.LoadFont("testdata/assets/fonts/luxisr.ttf")
		if err != nil {
			t.Fatal("failed to load font into cache")
		}
		_, err = c.GetFont("luxisr.ttf")
		if err != nil {
			t.Fatalf("failed to get cached font: %v", err)
		}
	})
	t.Run("GetUncached", func(t *testing.T) {
		c := NewCache()
		_, err := c.GetFont("luxisr.ttf")
		if err == nil {
			t.Fatalf("expected error getting uncached font")
		}
	})
}

func TestFont_Fallback(t *testing.T) {
	fg := FontGenerator{
		File:  "testdata/assets/fonts/luxisr.ttf",
		Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
		FontOptions: FontOptions{
			Size: 13.0,
			DPI:  44.0,
		},
	}
	f, err := fg.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	fg.File = "testdata/assets/fonts/seguiemj.ttf"
	emjfont, err := fg.Generate()
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	f.Fallbacks = append(f.Fallbacks, emjfont)

	f.MeasureString("a😀b😃c😄d😁e本")
	txt := f.NewText("a😀b😃c😄d😁e本", 0, 0)
	txt.Draw(image.NewRGBA(image.Rect(0, 0, 200, 200)), 0, 0)
}
