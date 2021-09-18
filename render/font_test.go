package render

import (
	"errors"
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/oakmound/oak/v3/oakerr"
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
