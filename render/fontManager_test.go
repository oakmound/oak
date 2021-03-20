package render

import (
	"image"
	"testing"

	"github.com/oakmound/oak/v2/oakerr"
)

func TestFontManager(t *testing.T) {

	initTestFont()

	fm := NewFontManager()

	fm.Get("def")
	// That may or may not be nil depending on if this is being run in a -coverprofile test
	// or not. Todo: fiddle with fonts and fix it
	//assert.NotNil(t, f)
	f := fm.Get("other")
	if f != nil {
		t.Fatalf("other should not be a defined font")
	}

	fg := FontGenerator{
		RawFile: luxisrTTF,
		Color:   image.Black,
	}

	err := fm.NewFont("other", fg)
	if err != nil {
		t.Fatalf("new font should not have failed")
	}

	f = fm.Get("other")
	if f == nil {
		t.Fatalf("other should be a defined font after it was set")
	}

	err = fm.NewFont("def", fg)
	if err == nil {
		t.Fatalf("new font under def name should have errored`")
	}
	if exists, ok := err.(oakerr.ExistingElement); ok {
		if !exists.Overwritten {
			t.Fatalf("def should have been overwritten")
		}
	}
}
