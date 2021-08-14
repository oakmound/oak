package oakerr

import (
	"testing"
)

func TestSetLanguageString(t *testing.T) {
	err := SetLanguageString("Gibberish")
	if err == nil {
		t.Fatal("Setting to language Gibberish did not error")
	}
	SetLanguageString("German")
	if currentLanguage != DEU {
		t.Fatalf("German did not set language to Deutsch")
	}
	SetLanguageString("English")
	if currentLanguage != ENG {
		t.Fatalf("English did not set language to English")
	}
	SetLanguageString("Japanese")
	if currentLanguage != JPN {
		t.Fatalf("Japanese did not set language to 日本語")
	}
	SetLanguageString("日本語")
	if currentLanguage != JPN {
		t.Fatalf("日本語 did not set language to 日本語")
	}
}
