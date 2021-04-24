package oakerr

import (
	"strings"

	"github.com/oakmound/oak/v2/dlog"
)

// Language configures the language of returned error strings
type Language int

var (
	currentLanguage Language = English
)

func SetLanguage(l Language) {
	currentLanguage = l
}

// SetLanguageString parses a string as a language
func SetLanguageString(s string) {
	s = strings.ToUpper(s)
	switch s {
	case "ENGLISH":
		currentLanguage = English
	case "GERMAN", "DEUTSCH":
		currentLanguage = Deutsch
	case "JAPANESE", "日本語":
		currentLanguage = 日本語
	default:
		// This should be the only always-english language string logged or returned by the engine
		dlog.Warn("Unknown language:", s, "Language set to English")
		currentLanguage = English
	}
}

// Valid languages
// TODO: should we be using ISO 639 codes? If so, ISO 639-1? ISO 639-3?
const (
	English Language = iota
	Deutsch
	日本語
)

// Q: Why these languages?
// A: These are the languages I (200sc) know or am actively learning
