package oak

import (
	"strings"

	"github.com/oakmound/oak/dlog"
)

// Language is hypothetically something games might care about in their text
// work: Consider moving this to oakerr, and also moving all strings to oakerr so
// messages output by the engine are localized.
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
