//+build !nolog

// Package dlog provides logging functions with
// caller file and line information,
// logging levels and level and text filters.
package dlog

import (
	"bufio"
	"bytes"
	"fmt"

	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LogLevel represents the levels a debug message can have
type LogLevel int

// Logging levels
const (
	NONE LogLevel = iota
	ERROR
	WARN
	INFO
	VERBOSE
)

var (
	byt         = bytes.NewBuffer(make([]byte, 0))
	debugLevel  = ERROR
	debugFilter = ""
	writer      *bufio.Writer
)

// dLog, the primary function of the package,
// prints out and writes to file a string
// containing the logged data separated by spaces,
// prepended with file and line information.
// It only includes logs which pass the current filters.
// Todo: use io.Multiwriter to simplify the writing to
// both logfiles and stdout
func dLog(console, override bool, in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !checkFilter(f, in) && !override {
			return
		}

		// Note on errors: these functions all return
		// errors, but they are always nil.
		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(strconv.Itoa(line))
		byt.WriteString("]  ")
		for _, elem := range in {
			byt.WriteString(fmt.Sprintf("%v ", elem))
		}
		byt.WriteRune('\n')

		if console {
			fmt.Print(byt.String())
		}

		if writer != nil {
			_, err := writer.WriteString(byt.String())
			if err != nil {
				// We can't log errors while we are in the error
				// logging function.
				panic(err)
			}
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
		}

		byt.Reset()
	}
}

// FileWrite runs dLog, but JUST writes to file instead
// of also to stdout.
func FileWrite(in ...interface{}) {
	dLog(false, true, in...)
}

func truncateFileName(f string) string {
	index := strings.LastIndex(f, "/")
	lIndex := strings.LastIndex(f, ".")
	return f[index+1 : lIndex]
}

func checkFilter(f string, in ...interface{}) bool {
	ret := false
	for _, elem := range in {
		ret = ret || strings.Contains(fmt.Sprintf("%s", elem), debugFilter)
	}
	return ret || strings.Contains(f, debugFilter)
}

// SetDebugFilter sets the string which determines
// what debug messages get printed. Only messages
// which contain the filer as a pseudo-regex
func SetDebugFilter(filter string) {
	debugFilter = filter
}

// SetDebugLevel sets what message levels of debug
// will be printed.
func SetDebugLevel(dL LogLevel) {
	if dL < NONE || dL > VERBOSE {
		Warn("Unknown debug level: ", dL)
		debugLevel = NONE
	} else {
		debugLevel = dL
	}
}

// CreateLogFile creates a file in the 'logs' directory
// of the starting point of this program to write logs to
func CreateLogFile() {
	file := "logs/dlog"
	file += time.Now().Format("_Jan_2_15-04-05_2006")
	file += ".txt"
	fHandle, err := os.Create(file)
	if err != nil {
		// We can't log an error that comes from
		// our error logging functions
		//panic(err)
		// But this is also not an error we want to panic on!
		fmt.Println("[oak]-------- No logs directory found. No logs will be written to file.")
		return
	}
	writer = bufio.NewWriter(fHandle)
}

// Error will write a dlog if the debug level is not NONE
func Error(in ...interface{}) {
	if debugLevel > NONE {
		dLog(true, true, in)
	}
}

// Warn will write a dLog if the debug level is higher than ERROR
func Warn(in ...interface{}) {
	if debugLevel > ERROR {
		dLog(true, true, in)
	}
}

// Info will write a dLog if the debug level is higher than WARN
func Info(in ...interface{}) {
	if debugLevel > WARN {
		dLog(true, false, in)
	}
}

// Verb will write a dLog if the debug level is higher than INFO
func Verb(in ...interface{}) {
	if debugLevel > INFO {
		dLog(true, false, in)
	}
}

// SetStringDebugLevel parses the input string as one of the debug levels
func SetStringDebugLevel(debugL string) {

	var dLevel LogLevel
	switch debugL {
	case "INFO":
		dLevel = INFO
	case "VERBOSE":
		dLevel = VERBOSE
	case "ERROR":
		dLevel = ERROR
	case "WARN":
		dLevel = WARN
	case "NONE":
		dLevel = NONE
	default:
		dLevel = ERROR
		fmt.Println("setting dlog level to \"", debugL, "\" failed, it is now set to ERROR")
	}

	SetDebugLevel(dLevel)
}
