package oak

import (
	"strings"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	Lang int
)

const (
	ENGLISH = iota
	GERMAN
)

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
