package oakerr

import (
	"strings"
)

// Language configures the language of returned error strings
type Language int

var (
	currentLanguage Language
)

func SetLanguage(l Language) {
	currentLanguage = l
}

// SetLanguageString parses a string as a language
func SetLanguageString(language string) error {
	language = strings.ToUpper(language)
	switch language {
	case "EN", "ENGLISH":
		currentLanguage = ENG
	case "DE", "GERMAN", "DEUTSCH":
		currentLanguage = DEU
	case "JP", "JAPANESE", "日本語":
		currentLanguage = JPN
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
