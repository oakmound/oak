package oakerr

import (
	"strings"

	"github.com/oakmound/oak/v3/dlog"
)

// Language configures the language of returned error strings
type Language int

var (
	currentLanguage Language = EN
)

func SetLanguage(l Language) {
	currentLanguage = l
}

// SetLanguageString parses a string as a language
func SetLanguageString(s string) {
	s = strings.ToUpper(s)
	switch s {
	case "EN", "ENGLISH":
		currentLanguage = EN
	case "DE", "GERMAN", "DEUTSCH":
		currentLanguage = DE
	case "JP", "JAPANESE", "日本語":
		currentLanguage = JP
	default:
		// This should be the only always-english language string logged or returned by the engine
		dlog.Warn("Unknown language:", s, "Language set to English")
		currentLanguage = EN
	}
}

// Valid languages, approximately matching ISO 639-1
const (
	EN Language = iota
	DE
	JP
)

// Q: Why these languages?
// A: These are the languages I (200sc) know or am actively learning
