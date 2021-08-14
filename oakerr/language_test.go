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
	if CurrentLanguage != DEU {
		t.Fatalf("German did not set language to Deutsch")
	}
	SetLanguageString("English")
	if CurrentLanguage != ENG {
		t.Fatalf("English did not set language to English")
	}
	SetLanguageString("Japanese")
	if CurrentLanguage != JPN {
		t.Fatalf("Japanese did not set language to 日本語")
	}
	SetLanguageString("日本語")
	if CurrentLanguage != JPN {
		t.Fatalf("日本語 did not set language to 日本語")
	}
}
