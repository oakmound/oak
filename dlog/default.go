package dlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/oakmound/oak/v3/oakerr"
)

var (
	_ Logger = &logger{}
)

type logger struct {
	bytPool     sync.Pool
	debugLevel  Level
	debugFilter func(string) bool
	writer      io.Writer
}

// NewLogger returns an instance of the default logger with no filter,
// no file, and level set to ERROR
func NewLogger() Logger {
	return &logger{
		bytPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		debugLevel: ERROR,
		writer:     os.Stdout,
	}
}

// dLog, the primary function of the package,
// prints out and writes to file a string
// containing the logged data separated by spaces,
// prepended with file and line information.
// It only includes logs which pass the current filters.
func (l *logger) dLog(level Level, in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	// TODO: restructure so dlog functions work like t.Helper,
	// incrementing this traceBack value for us
	traceBack := 2
	_, f, line, ok := runtime.Caller(traceBack)
	for strings.Contains(f, "dlog") {
		traceBack++
		_, f, line, ok = runtime.Caller(traceBack)
	}
	if ok {
		f = truncateFileName(f)
		if !l.checkFilter(f, in) {
			return
		}

		buffer := l.bytPool.Get().(*bytes.Buffer)
		// Note on errors: these functions all return
		// errors, but they are always nil.
		buffer.WriteRune('[')
		buffer.WriteString(f)
		buffer.WriteRune(':')
		buffer.WriteString(strconv.Itoa(line))
		buffer.WriteString("]  ")
		buffer.WriteString(logLevels[level])
		buffer.WriteRune(':')
		for _, elem := range in {
			buffer.WriteString(fmt.Sprintf(" %v", elem))
		}
		buffer.WriteRune('\n')

		// This can error, but we can't do anything about it if it does.
		l.writer.Write(buffer.Bytes())

		buffer.Reset()
		l.bytPool.Put(buffer)
	}
}

func truncateFileName(f string) string {
	directoryIndex := strings.LastIndex(f, "/")
	extensionIndex := strings.LastIndex(f, ".")
	return f[directoryIndex+1 : extensionIndex]
}

func (l *logger) checkFilter(f string, in ...interface{}) bool {
	if l.debugFilter == nil {
		return true
	}
	check := f
	for _, elem := range in {
		check += fmt.Sprintf(" %v", elem)
	}
	return l.debugFilter(check)
}

// SetFilter defines a custom filter function. Log lines that
// return false when passed to this function will not be output.
func (l *logger) SetFilter(filter func(string) bool) {
	l.debugFilter = filter
}

// SetLogLevel sets what message levels of debug
// will be printed.
func (l *logger) SetLogLevel(level Level) error {
	if level < NONE || level > VERBOSE {
		return oakerr.InvalidInput{
			InputName: "level",
		}
	}
	l.debugLevel = level
	return nil
}

// Error will write a dlog if the debug level is not NONE
func (l *logger) Error(in ...interface{}) {
	if l.debugLevel > NONE {
		l.dLog(ERROR, in...)
	}
}

// Info will write a dLog if the debug level is higher than WARN
func (l *logger) Info(in ...interface{}) {
	if l.debugLevel > ERROR {
		l.dLog(INFO, in...)
	}
}

// Verb will write a dLog if the debug level is higher than INFO
func (l *logger) Verb(in ...interface{}) {
	if l.debugLevel > INFO {
		l.dLog(VERBOSE, in...)
	}
}

func (l *logger) SetOutput(w io.Writer) {
	l.writer = w
}
