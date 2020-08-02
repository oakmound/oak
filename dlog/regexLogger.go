package dlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	_ FullLogger = &RegexLogger{}
)

// RegexLogger is a logger implementation that offers some
// additional features on top of the default logger.
// Todo v3: combine logger implementations.
type RegexLogger struct {
	debugLevel Level

	debugFilter string
	filterRegex *regexp.Regexp
	// FilterOverrideLevel is the log level at which
	// logs will be shown regardless of the filter.
	FilterOverrideLevel Level

	writer io.Writer
	file   io.Writer
}

// NewRegexLogger returns a custom logger that writes to os.Stdout and
// overrides filters on WARN or higher messages.
func NewRegexLogger(level Level) *RegexLogger {
	return &RegexLogger{
		debugLevel:          level,
		writer:              os.Stdout,
		FilterOverrideLevel: WARN,
	}
}

// GetLogLevel returns the current log level, i.e WARN or INFO...
func (l *RegexLogger) GetLogLevel() Level {
	return l.debugLevel
}

// dLog, the primary function of the package,
// prints out and writes to file a string
// containing the logged data separated by spaces,
// prepended with file and line information.
// It only includes logs which pass the current filters.
func (l *RegexLogger) dLog(w io.Writer, override bool, in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if strings.Contains(f, "dlog") {
		_, f, line, ok = runtime.Caller(3)
	}
	if ok {
		var bldr strings.Builder
		f = truncateFileName(f)
		// Note on errors: these functions all return
		// errors, but they are always nil.
		bldr.WriteRune('[')
		bldr.WriteString(f)
		bldr.WriteRune(':')
		bldr.WriteString(strconv.Itoa(line))
		bldr.WriteString("]  ")
		bldr.WriteString(logLevels[l.GetLogLevel()])
		bldr.WriteRune(':')
		for _, elem := range in {
			bldr.WriteString(fmt.Sprintf("%v ", elem))
		}
		bldr.WriteRune('\n')
		fullLog := []byte(bldr.String())

		if !override && !l.checkFilter(fullLog) {
			return
		}

		_, err := w.Write(fullLog)
		if err != nil {
			fmt.Println("Logging error", err)
		}
	}
}

func (l *RegexLogger) checkFilter(fullLog []byte) bool {
	if l.debugFilter == "" {
		return true
	}
	if l.filterRegex != nil {
		return l.filterRegex.Match(fullLog)
	}
	return bytes.Contains(fullLog, []byte(l.debugFilter))
}

// SetDebugFilter sets the string which determines
// what debug messages get printed. Only messages
// which contain the filer as a pseudo-regex
func (l *RegexLogger) SetDebugFilter(filter string) {
	l.debugFilter = filter
	var err error
	l.filterRegex, err = regexp.Compile(filter)
	if err != nil {
		l.Error("could not compile filter regex", err)
	}
}

// SetDebugLevel sets what message levels of debug
// will be printed.
func (l *RegexLogger) SetDebugLevel(dL Level) {
	if dL < NONE || dL > VERBOSE {
		l.Warn("Unknown debug level: ", dL)
		l.debugLevel = NONE
	} else {
		l.debugLevel = dL
	}
}

// CreateLogFile creates a file in the 'logs' directory
// of the starting point of this program to write logs to
func (l *RegexLogger) CreateLogFile() {
	file := "logs/dlog"
	file += time.Now().Format("_Jan_2_15-04-05_2006")
	file += ".txt"
	var err error
	l.file, err = os.Create(file)
	if err != nil {
		fmt.Println("[oak]-------- No logs directory found. No logs will be written to file.")
		return
	}
	l.writer = io.MultiWriter(l.file, l.writer)
}

// FileWrite acts just like a regular write on a RegexLogger. It does
// not respect overrides.
func (l *RegexLogger) FileWrite(in ...interface{}) {
	if l.file == nil {
		return
	}
	l.dLog(l.file, true, in...)
}

// Error will write a dlog if the debug level is not NONE
func (l *RegexLogger) Error(in ...interface{}) {
	if l.debugLevel > NONE {
		l.dLog(l.writer, l.FilterOverrideLevel > NONE, in)
	}
}

// Warn will write a dLog if the debug level is higher than ERROR
func (l *RegexLogger) Warn(in ...interface{}) {
	if l.debugLevel > ERROR {
		l.dLog(l.writer, l.FilterOverrideLevel > ERROR, in)
	}
}

// Info will write a dLog if the debug level is higher than WARN
func (l *RegexLogger) Info(in ...interface{}) {
	if l.debugLevel > WARN {
		l.dLog(l.writer, l.FilterOverrideLevel > WARN, in)
	}
}

// Verb will write a dLog if the debug level is higher than INFO
func (l *RegexLogger) Verb(in ...interface{}) {
	if l.debugLevel > INFO {
		l.dLog(l.writer, l.FilterOverrideLevel > INFO, in)
	}
}

// SetWriter sets the writer that RegexLogger logs to
func (l *RegexLogger) SetWriter(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("cannot write to nil writer")
	}
	l.writer = w
	return nil
}
