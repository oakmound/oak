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

	callLogger := func() {
		lgr.FileWrite("fileWrite")
		lgr.Error("error")
		lgr.Warn("warn")
		lgr.Info("info")
		lgr.Verb("verb")

		lgr.SetDebugFilter("foo")
		lgr.Verb("bar")
		lgr.Verb("foo")
	}
	callLogger()

	expectedOut := `[default_test:39]  VERBOSE:fileWrite 
[default_test:39]  VERBOSE:error 
[default_test:39]  VERBOSE:warn 
[default_test:39]  VERBOSE:info 
[default_test:39]  VERBOSE:verb 
[default_test:39]  VERBOSE:foo 
`

	out := string(buff.Bytes())

	if out != expectedOut {
		t.Fatalf("logged output did not match: got %q expected %q", out, expectedOut)
	}

	lgr.CreateLogFile()
}
