package render

import (
	"image"
	"sync"
	"testing"
)

func TestTextFns(t *testing.T) {
	initTestFont()

	txt := DefFont().NewStrText("Test", 0, 0)

	fg := FontGenerator{
		RawFile: luxisrTTF,
		Color:   image.Black,
	}

	f, _ := fg.Generate()

	txt.SetFont(f)
	if f != txt.d {
		t.Fatalf("text set font failed")
	}

	txt.SetString("Test2")
	if "Test2" != txt.text.String() {
		t.Fatalf("text SetString failed")
	}
	if "Test2" != txt.StringLiteral() {
		t.Fatalf("text SetString failed")
	}

	n := 100
	txt.SetIntP(&n)

	n = 50
	if txt.text.String() != "50" {
		t.Fatalf("text SetIntP failed")
	}

	txt.SetInt(n + 1)
	if txt.text.String() != "51" {
		t.Fatalf("text SetInt failed")
	}

	txt.SetStringer(dummyStringer{})
	if txt.text.String() != "Dummy" {
		t.Fatalf("text SetText failed")
	}
	if txt.String() != "Text[Dummy]" {
		t.Fatalf("text String() failed")
	}

	txts := txt.Wrap(1, 10)
	if len(txts) != 5 {
		t.Fatalf("wrap did not wrap dummy to multi line")
	}

	for i, wrap := range txts {
		if float64(i*10) != wrap.Y() {
			t.Fatalf("wrapped texts did not have changed y values")
		}
	}

	txts = txt.Wrap(3, 10)
	if len(txts) != 2 {
		t.Fatalf("wrap did not wrap dummy to multi line")
	}

	if txt.ToSprite() == nil {
		t.Fatalf("to sprite failed")
	}

	txt.Center()
	if txt.X() != float64(-2) {
		t.Fatalf("center did not move text's x value")
	}
}

func TestText_StringPtr(t *testing.T) {
	initTestFont()
	s := new(string)
	*s = "hello"
	txt := NewStrPtrText(s, 0, 0)
	if txt.StringLiteral() != "hello" {
		t.Fatalf("str ptr text not set on creation")
	}
	*s = "goodbye"
	if txt.StringLiteral() != "goodbye" {
		t.Fatalf("str ptr text not set by pointer manipulation")
	}
	txt.SetStringPtr(nil)
	if txt.StringLiteral() != "nil" {
		t.Fatalf("nil str ptr text failed")
	}
}

type dummyStringer struct{}

func (d dummyStringer) String() string {
	return "Dummy"
}

var initTestFontOnce sync.Once

// Todo: move this to font_test.go, once we have font_test.go
func initTestFont() {
	initTestFontOnce.Do(func() {
		DefFontGenerator = FontGenerator{RawFile: luxisrTTF}
		SetFontDefaults("", "", "", "", "white", "", 10, 10)
	})
}
