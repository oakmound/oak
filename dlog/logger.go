package dlog

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	byt       = bytes.NewBuffer(make([]byte, 0))
	theFilter = ""
)

func DLog(s string) {
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

func DILog(s string, nums ...int) {

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

// dlog.WarnI(somestring)
// dlog.InfoI()
// dlog.VerboseI()

// Verbose
// Info
// Warn
