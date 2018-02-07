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
