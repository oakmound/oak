package dlog

import "errors"

// A Logger is a minimal log interface for the content oak wants to log:
// four levels of logging.
type Logger interface {
	Error(...interface{})
	Warn(...interface{})
	Info(...interface{})
	Verb(...interface{})
}

// OakLogger is the Logger which all oak log functions are passed through.
// If this is not manually set through SetLogger, oak will initialize this
// to the an instance of the private logger type
var oakLogger Logger

// ErrorCheck checks that the input is not nil, then calls Error on it if it is
// not. Otherwise it does nothing.
func ErrorCheck(in error) {
	if in != nil {
		Error(in)
	}
}

// Error will write a log if the debug level is not NONE
var Error = func(...interface{}) {}

// Warn will write a log if the debug level is higher than ERROR
var Warn = func(...interface{}) {}

// Info will write a log if the debug level is higher than WARN
var Info = func(...interface{}) {}

// Verb will write a log if the debug level is higher than INFO
var Verb = func(...interface{}) {}

// SetLogger defines what logger should be used for oak's internal logging.
// If this is NOT called before oak.Init is called (assuming this is being
// used with oak), then it will be called with the default logger as a part
// of oak.Init.
func SetLogger(l Logger) {
	_, isDefault := l.(*logger)
	if isDefault && oakLogger != nil {
		// The user set the logger themselves,
		// don't reset to the default logger
		return
	}
	oakLogger = l
	Error = l.Error
	Warn = l.Warn
	Info = l.Info
	Verb = l.Verb
	// If this logger supports the additional functionality described
	// by the FullLogger interface, enable those functions. Otherwise
	// they are NOPs. (the default logger supports these functions.)
	if fl, ok := l.(FullLogger); ok {
		fullOakLogger = fl
		FileWrite = fl.FileWrite
		GetLogLevel = fl.GetLogLevel
		SetDebugFilter = fl.SetDebugFilter
		SetDebugLevel = fl.SetDebugLevel
		CreateLogFile = fl.CreateLogFile
	}
}

// ParseDebugLevel parses the input string as a known debug levels
func ParseDebugLevel(s string) (Level, error) {
	switch s {
	case "INFO":
		return INFO, nil
	case "VERBOSE":
		return VERBOSE, nil
	case "ERROR":
		return ERROR, nil
	case "WARN":
		return WARN, nil
	case "NONE":
		return NONE, nil
	default:
		return ERROR, errors.New("parsing dlog level of \"" + s + "\" failed")
	}
}
