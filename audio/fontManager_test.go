package audio

import (
	"testing"

	"github.com/200sc/klangsynthese/font"
)

func TestFontManager(t *testing.T) {
	fm := NewFontManager()
	if fm.NewFont("unused", font.New()) != nil {
		t.Fatalf("expected new font to succeed")
	}
	if fm.NewFont("unused", font.New()) == nil {
		t.Fatalf("expected duplicate font to fail")
	}
	if fm.Get("notafont") != nil {
		t.Fatalf("expected non existant get font to fail")
	}
	if fm.GetDefault() == nil {
		t.Fatalf("expected def get font to succeed")
	}
}
