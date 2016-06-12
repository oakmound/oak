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

func dLog(in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !checkFilter(f, in) {
			return
		}

		lineStr := strconv.Itoa(line)

		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(lineStr)
		byt.WriteString("]  ")
		for _, elem := range in {
			byt.WriteString(fmt.Sprintf("%s ", elem))
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

func checkFilter(f string, in ...interface{}) bool {
	ret := false
	for _, elem := range in {
		ret = ret || strings.Contains(fmt.Sprintf("%s", elem), theFilter)
	}
	return ret || strings.Contains(f, theFilter)
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

func Error(in ...interface{}) {
	if debugLevel > 0 {
		dLog(in)
	}
}

func Warn(in ...interface{}) {
	if debugLevel > 1 {
		dLog(in)
	}
}

func Info(in ...interface{}) {
	if debugLevel > 2 {
		dLog(in)
	}
}

func Verb(in ...interface{}) {
	if debugLevel > 3 {
		dLog(in)
	}
}

// dlog.Warn()
// dlog.Info()
// dlog.Verb()

// Verbose
// Info
// Warn
