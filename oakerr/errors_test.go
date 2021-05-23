package oakerr

import (
	"testing"
)

func TestErrorsAreErrors(t *testing.T) {
	languages := []Language{English, Deutsch}
	for _, lang := range languages {
		SetLanguage(lang)
		var err error = NotFound{}
		if err.Error() == "" {
			t.Fatalf("NotFound error was empty")
		}
		err = ExistingElement{}
		if err.Error() == "" {
			t.Fatalf("ExistingElement error was empty")
		}
		err = ExistingElement{Overwritten: true}
		if err.Error() == "" {
			t.Fatalf("ExistingElement error was empty")
		}
		err = InsufficientInputs{}
		if err.Error() == "" {
			t.Fatalf("InsufficientInputs error was empty")
		}
		err = InvalidInput{}
		if err.Error() == "" {
			t.Fatalf("InvalidInput error was empty")
		}
		err = NilInput{}
		if err.Error() == "" {
			t.Fatalf("NilInput error was empty")
		}
		err = IndivisibleInput{}
		if err.Error() == "" {
			t.Fatalf("IndivisibleInput error was empty")
		}
		err = UnsupportedFormat{}
		if err.Error() == "" {
			t.Fatalf("UnsupportedFormat error was empty")
		}
		err = UnsupportedPlatform{}
		if err.Error() == "" {
			t.Fatalf("UnsupportedPlatform error was empty")
		}
	}
	// Assert nothing crashed
}

func TestErrorFallback(t *testing.T) {
	SetLanguage(日本語)
	s := errorString(codeIndivisibleInput, "a", "b")
	if s != "a was not divisible by b" {
		t.Fatalf("language fallback to english failed")
	}
}
