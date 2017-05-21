package oak

import (
	"strings"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

// Language is hypothetically something games might care about in their text
type Language int

var (
	// Lang is the current langugae
	Lang Language
)

// Lang enumerator
const (
	ENGLISH Language = iota
	GERMAN
)

// SetLang parses a string as a language
func SetLang(s string) {
	s = strings.ToUpper(s)
	switch s {
	case "ENGLISH":
		Lang = ENGLISH
	case "GERMAN":
		Lang = GERMAN
	default:
		dlog.Warn("Unknown language string:", s, "Language set to English")
		Lang = ENGLISH
	}
}
