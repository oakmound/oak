package oakerr

import (
	"testing"
)

func TestSetLanguageString(t *testing.T) {
	SetLanguageString("Gibberish")
	if currentLanguage != EN {
		t.Fatalf("Gibberish did not set language to English")
	}
	SetLanguageString("German")
	if currentLanguage != DE {
		t.Fatalf("German did not set language to Deutsch")
	}
	SetLanguageString("English")
	if currentLanguage != EN {
		t.Fatalf("English did not set language to English")
	}
	SetLanguageString("Japanese")
	if currentLanguage != JP {
		t.Fatalf("Japanese did not set language to 日本語")
	}
	SetLanguageString("日本語")
	if currentLanguage != JP {
		t.Fatalf("日本語 did not set language to 日本語")
	}
}
