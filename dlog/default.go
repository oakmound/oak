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

type logger struct {
	byt         *bytes.Buffer
	debugLevel  Level
	debugFilter string
	writer      *bufio.Writer
}

// NewLogger returns an instance of the default logger with no filter,
// no file, and level set to ERROR
func NewLogger() Logger {
	return &logger{
		byt:        bytes.NewBuffer(make([]byte, 0)),
		debugLevel: ERROR,
	}
}

// GetLogLevel returns the current log level, i.e WARN or INFO...
func (l *logger) GetLogLevel() Level {
	return l.debugLevel
}

// dLog, the primary function of the package,
// prints out and writes to file a string
// containing the logged data separated by spaces,
// prepended with file and line information.
// It only includes logs which pass the current filters.
// Todo: use io.Multiwriter to simplify the writing to
// both logfiles and stdout
func (l *logger) dLog(console, override bool, in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !l.checkFilter(f, in) && !override {
			return
		}

		// Note on errors: these functions all return
		// errors, but they are always nil.
		l.byt.WriteRune('[')
		l.byt.WriteString(f)
		l.byt.WriteRune(':')
		l.byt.WriteString(strconv.Itoa(line))
		l.byt.WriteString("]  ")
		for _, elem := range in {
			l.byt.WriteString(fmt.Sprintf("%v ", elem))
		}
		l.byt.WriteRune('\n')

		if console {
			fmt.Print(l.byt.String())
		}

		if l.writer != nil {
			_, err := l.writer.WriteString(l.byt.String())
			if err != nil {
				// We can't log errors while we are in the error
				// logging function.
				fmt.Println("Logging error", err)
			}
			err = l.writer.Flush()
			if err != nil {
				fmt.Println("Logging error", err)
			}
		}

		l.byt.Reset()
	}
}

// FileWrite runs dLog, but JUST writes to file instead
// of also to stdout.
func (l *logger) FileWrite(in ...interface{}) {
	l.dLog(false, true, in...)
}

func truncateFileName(f string) string {
	index := strings.LastIndex(f, "/")
	lIndex := strings.LastIndex(f, ".")
	return f[index+1 : lIndex]
}

func (l *logger) checkFilter(f string, in ...interface{}) bool {
	ret := false
	for _, elem := range in {
		ret = ret || strings.Contains(fmt.Sprintf("%s", elem), l.debugFilter)
	}
	return ret || strings.Contains(f, l.debugFilter)
}

// SetDebugFilter sets the string which determines
// what debug messages get printed. Only messages
// which contain the filer as a pseudo-regex
func (l *logger) SetDebugFilter(filter string) {
	l.debugFilter = filter
}

// SetDebugLevel sets what message levels of debug
// will be printed.
func (l *logger) SetDebugLevel(dL Level) {
	if dL < NONE || dL > VERBOSE {
		Warn("Unknown debug level: ", dL)
		l.debugLevel = NONE
	} else {
		l.debugLevel = dL
	}
}

// CreateLogFile creates a file in the 'logs' directory
// of the starting point of this program to write logs to
func (l *logger) CreateLogFile() {
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
	l.writer = bufio.NewWriter(fHandle)
}

// Error will write a dlog if the debug level is not NONE
func (l *logger) Error(in ...interface{}) {
	if l.debugLevel > NONE {
		l.dLog(true, true, in)
	}
}

// Warn will write a dLog if the debug level is higher than ERROR
func (l *logger) Warn(in ...interface{}) {
	if l.debugLevel > ERROR {
		l.dLog(true, true, in)
	}
}

// Info will write a dLog if the debug level is higher than WARN
func (l *logger) Info(in ...interface{}) {
	if l.debugLevel > WARN {
		l.dLog(true, false, in)
	}
}

// Verb will write a dLog if the debug level is higher than INFO
func (l *logger) Verb(in ...interface{}) {
	if l.debugLevel > INFO {
		l.dLog(true, false, in)
	}
}
