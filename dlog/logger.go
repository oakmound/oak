package dlog

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	byt        = bytes.NewBuffer(make([]byte, 0))
	debugLevel = 0
	theFilter  = ""
)

func dLog(s string) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !checkFilter(f, s) {
			return
		}

		lineStr := strconv.Itoa(line)

		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(lineStr)
		byt.WriteString("]  ")
		byt.WriteString(s)
		byt.WriteRune('\n')

		fmt.Print(byt.String())

		byt.Reset()

		// [filename:lineNum]  output
	}
}

func dILog(s string, nums ...int) {

	_, f, line, ok := runtime.Caller(2)
	if ok {

		f = truncateFileName(f)
		lineStr := strconv.Itoa(line)
		if !checkFilter(f, s) {
			return
		}

		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(lineStr)
		byt.WriteString("]  ")
		byt.WriteString(s)
		for _, num := range nums {
			byt.WriteRune(' ')
			byt.WriteString(strconv.Itoa(num))
		}
		byt.WriteRune('\n')

		fmt.Print(byt.String())

		byt.Reset()
		// [filename:lineNum]  output
	}
}

func truncateFileName(f string) string {
	index := strings.LastIndex(f, "/")
	lIndex := strings.LastIndex(f, ".")
	return f[index+1 : lIndex]
}

func checkFilter(f, s string) bool {
	return strings.Contains(s, theFilter) || strings.Contains(f, theFilter)
}

func SetDebugFilter(filter string) {
	theFilter = filter
}

func SetDebugLevel(dL string) {
	switch dL {
	case "VERBOSE":
		debugLevel = 4
	case "INFO":
		debugLevel = 3
	case "WARN":
		debugLevel = 2
	case "ERROR":
		debugLevel = 1
	default:
		debugLevel = 0
	}
}

func Error(s string) {
	if debugLevel > 0 {
		dLog(s)
	}
}
func ErrorI(s string, nums ...int) {
	if debugLevel > 0 {
		dILog(s, nums...)
	}
}

func Warn(s string) {
	if debugLevel > 1 {
		dLog(s)
	}
}
func WarnI(s string, nums ...int) {
	if debugLevel > 1 {
		dILog(s, nums...)
	}
}

func Info(s string) {
	if debugLevel > 2 {
		dLog(s)
	}
}
func InfoI(s string, nums ...int) {
	if debugLevel > 2 {
		dILog(s, nums...)
	}
}

func Verb(s string) {
	if debugLevel > 3 {
		dLog(s)
	}
}
func VerbI(s string, nums ...int) {
	if debugLevel > 3 {
		dILog(s, nums...)
	}
}

// dlog.WarnI(somestring)
// dlog.InfoI()
// dlog.VerboseI()

// Verbose
// Info
// Warn
