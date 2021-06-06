package dlog

import (
	"strings"

	"github.com/oakmound/oak/v3/oakerr"
)

// Level represents the levels a debug message can have
type Level int

// Level values const
const (
	NONE Level = iota
	ERROR
	INFO
	VERBOSE
)

var logLevels = map[Level]string{
	NONE:    "NONE",
	ERROR:   "ERROR",
	INFO:    "INFO",
	VERBOSE: "VERBOSE",
}

func (l Level) String() string {
	return logLevels[l]
}

// ParseDebugLevel parses the input string as a known debug levels
func ParseDebugLevel(level string) (Level, error) {
	level = strings.ToUpper(level)
	switch level {
	case "INFO":
		return INFO, nil
	case "VERBOSE":
		return VERBOSE, nil
	case "ERROR":
		return ERROR, nil
	case "NONE":
		return NONE, nil
	default:
		return ERROR, oakerr.InvalidInput{InputName: "level"}
	}
}
