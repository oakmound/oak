package oakerr

import (
	"testing"
)

func TestSetLanguageString(t *testing.T) {
	err := SetLanguageString("Gibberish")
	if err == nil {
		t.Fatal("Setting to language Gibberish did not error")
	}
	err = SetLanguageString("German")
	if err != nil {
		t.Fatalf("SetLanguageString failed: %v", err)
	}
	if CurrentLanguage != DEU {
		t.Fatalf("German did not set language to Deutsch")
	}
	err = SetLanguageString("English")
	if err != nil {
		t.Fatalf("SetLanguageString failed: %v", err)
	}
	if CurrentLanguage != ENG {
		t.Fatalf("English did not set language to English")
	}
	err = SetLanguageString("Japanese")
	if err != nil {
		t.Fatalf("SetLanguageString failed: %v", err)
	}
	if CurrentLanguage != JPN {
		t.Fatalf("Japanese did not set language to 日本語")
	}
	err = SetLanguageString("日本語")
	if err != nil {
		t.Fatalf("SetLanguageString failed: %v", err)
	}
	if CurrentLanguage != JPN {
		t.Fatalf("日本語 did not set language to 日本語")
	}
}
