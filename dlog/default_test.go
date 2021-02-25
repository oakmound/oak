package dlog

import (
	"bufio"
	"bytes"
	"testing"
)

func TestLogger(t *testing.T) {
	lgr := NewLogger().(*logger)

	defaultLevel := lgr.GetLogLevel()
	if defaultLevel != ERROR {
		t.Fatalf("expected default log level to be ERROR, was: %v", defaultLevel)
	}

	lgr.SetDebugLevel(-1)
	if lgr.GetLogLevel() != NONE {
		t.Fatalf("expected -1 log level to be NONE, was: %v", lgr.GetLogLevel())
	}

	lgr.SetDebugLevel(VERBOSE)

	var buff = new(bytes.Buffer)

	lgr.writer = bufio.NewWriter(buff)

	lgr.FileWrite("fileWrite")
	lgr.Error("error")
	lgr.Warn("warn")
	lgr.Info("info")
	lgr.Verb("verb")

	lgr.SetDebugFilter("foo")
	lgr.Verb("bar")
	lgr.Verb("foo")

	expectedOut := `[testing:1194]  VERBOSE:fileWrite 
[testing:1194]  VERBOSE:error 
[testing:1194]  VERBOSE:warn 
[testing:1194]  VERBOSE:info 
[testing:1194]  VERBOSE:verb 
[testing:1194]  VERBOSE:foo 
`

	out := string(buff.Bytes())

	if out != expectedOut {
		t.Fatalf("logged output did not match: got %q expected %q", out, expectedOut)
	}

	lgr.CreateLogFile()
}
