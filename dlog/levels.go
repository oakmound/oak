package dlog

// Level represents the levels a debug message can have
type Level int

// Level values const
const (
	NONE Level = iota
	ERROR
	WARN
	INFO
	VERBOSE
)

var logLevels = map[Level]string{
	NONE:    "NONE",
	ERROR:   "ERROR",
	WARN:    "WARN",
	INFO:    "INFO",
	VERBOSE: "VERBOSE",
}
