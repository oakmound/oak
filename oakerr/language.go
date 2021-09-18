package oakerr

import (
	"strings"
)

// Language configures the language of returned error strings
type Language int

var (
	// CurrentLanguage is the current language for error and log strings
	CurrentLanguage Language
)

// SetLanguageString parses a string as a language
func SetLanguageString(language string) error {
	language = strings.ToUpper(language)
	switch language {
	case "EN", "ENGLISH":
		CurrentLanguage = ENG
	case "DE", "GERMAN", "DEUTSCH":
		CurrentLanguage = DEU
	case "JP", "JAPANESE", "日本語":
		CurrentLanguage = JPN
	default:
		return InvalidInput{InputName: language}
	}
	return nil
}

// Valid languages, uppercase ISO 639-2
const (
	// English
	ENG Language = iota
	// German
	DEU
	// Japanese
	JPN
)

// Q: Why these languages?
// A: These are the languages I (200sc) know or am actively learning
