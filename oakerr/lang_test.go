package oakerr

import (
	"testing"
)

func TestSetLanguageString(t *testing.T) {
	SetLanguageString("Gibberish")
	if currentLanguage != English {
		t.Fatalf("Gibberish did not set language to English")
	}
	SetLanguageString("German")
	if currentLanguage != Deutsch {
		t.Fatalf("German did not set language to Deutsch")
	}
	SetLanguageString("English")
	if currentLanguage != English {
		t.Fatalf("English did not set language to English")
	}
	SetLanguageString("Japanese")
	if currentLanguage != 日本語 {
		t.Fatalf("Japanese did not set language to 日本語")
	}
	SetLanguageString("日本語")
	if currentLanguage != 日本語 {
		t.Fatalf("日本語 did not set language to 日本語")
	}
}
