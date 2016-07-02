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

const (
	NONE = iota
	ERROR
	WARN
	INFO
	VERBOSE
)

var (
	byt         = bytes.NewBuffer(make([]byte, 0))
	debugLevel  = NONE
	debugFilter = ""
	writer_p    *bufio.Writer
)

func dLog(in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !checkFilter(f, in) {
			return
		}

		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(strconv.Itoa(line))
		byt.WriteString("]  ")
		for _, elem := range in {
			byt.WriteString(fmt.Sprintf("%v ", elem))
		}
		byt.WriteRune('\n')

		fmt.Print(byt.String())
		//w := writer(io.Writer)
		//fmt.Fprintln(w, byt.String())
		writer := *writer_p
		writer.WriteString(byt.String())
		writer.Flush()

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
		ret = ret || strings.Contains(fmt.Sprintf("%s", elem), debugFilter)
	}
	return ret || strings.Contains(f, debugFilter)
}

func SetDebugFilter(filter string) {
	debugFilter = filter
}

func SetDebugLevel(dL int) {
	if dL < NONE || dL > VERBOSE {
		Warn("Unknown debug level: ", dL)
		debugLevel = NONE
	} else {
		debugLevel = dL
	}
}

func CreateLogFile() {
	file := "../logs/dlog"
	file += strconv.Itoa(int(time.Now().Unix()))
	file += ".txt"
	fHandle, _ := os.Create(file)
	//fHandle, _ := os.OpenFile("./dlog.txt", os.O_APPEND, 0666)

	writer_p = bufio.NewWriter(fHandle)
	//writer = *writer_p
	//defer fHandle.Close()
}

func Error(in ...interface{}) {
	if debugLevel > NONE {
		dLog(in)
	}
}

func Warn(in ...interface{}) {
	if debugLevel > ERROR {
		dLog(in)
	}
}

func Info(in ...interface{}) {
	if debugLevel > WARN {
		dLog(in)
	}
}

func Verb(in ...interface{}) {
	if debugLevel > INFO {
		dLog(in)
	}
}

// dlog.Warn()
// dlog.Info()
// dlog.Verb()

// Verbose
// Info
// Warn
