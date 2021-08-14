package render

import (
	"image"
	"testing"
)

func TestTextFns(t *testing.T) {
	txt := DefaultFont().NewText("Test", 0, 0)

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
	if txt.text.String() != "Test2" {
		t.Fatalf("text SetString failed")
	}
	if txt.StringLiteral() != "Test2" {
		t.Fatalf("text SetString failed")
	}

	n := 100
	txt.SetIntPtr(&n)

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
	if txt.X() != float64(-20) {
		t.Fatalf("center did not move text's x value: expected %v got %v", -20, txt.X())
	}
}

func TestText_StringPtr(t *testing.T) {
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
