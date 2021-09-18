package dlog

import "github.com/oakmound/oak/v3/oakerr"

type logCode int

// Constant log string identifiers. All log strings output by oak
// should be enumerated here.
const (
	WindowClosed logCode = iota
	SceneStarting
	SceneLooping
	SceneEnding
	UnknownScene
)

func (lc logCode) String() string {
	s := logstrings[oakerr.CurrentLanguage][lc]
	if s == "" {
		return logstrings[oakerr.ENG][lc]
	}
	return s
}

var logstrings = map[oakerr.Language]map[logCode]string{
	oakerr.ENG: {
		WindowClosed:  "Window closed",
		SceneStarting: "Scene start:",
		SceneLooping:  "Looping scene",
		SceneEnding:   "Scene end:",
		UnknownScene:  "Unknown scene:",
	},
	oakerr.DEU: {},
	oakerr.JPN: {},
}
