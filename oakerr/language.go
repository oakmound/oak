package oakerr

import (
	"strings"
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
func SetLanguageString(language string) error {
	language = strings.ToUpper(language)
	switch language {
	case "EN", "ENGLISH":
		currentLanguage = EN
	case "DE", "GERMAN", "DEUTSCH":
		currentLanguage = DE
	case "JP", "JAPANESE", "日本語":
		currentLanguage = JP
	default:
		return InvalidInput{InputName: language}
	}
	return nil
}

// Valid languages, approximately matching ISO 639-1
const (
	EN Language = iota
	DE
	JP
)

// Q: Why these languages?
// A: These are the languages I (200sc) know or am actively learning
