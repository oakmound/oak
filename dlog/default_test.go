package dlog

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	lgr := NewLogger().(*logger)

	defaultLevel := lgr.GetLogLevel()
	if defaultLevel != ERROR {
		t.Fatalf("expected default log level to be ERROR, was: %v", defaultLevel)
	}

	err := lgr.SetDebugLevel(-1)
	if err == nil {
		t.Fatalf("expected -1 log level to error")
	}

	lgr.SetDebugLevel(VERBOSE)

	var buff = new(bytes.Buffer)

	lgr.SetOutput(buff)

	// This function wrapper corrects the logged file generated
	calllogger := func() {
		lgr.Error("error")
		lgr.Info("info")
		lgr.Verb("verb")

		lgr.SetFilter(func(s string) bool { return strings.Contains(s, "foo") })
		lgr.Verb("bar")
		lgr.Verb("foo")
	}

	calllogger()

	expectedOut := `[default_test:39]  ERROR: error
[default_test:39]  INFO: info
[default_test:39]  VERBOSE: verb
[default_test:39]  VERBOSE: foo
`
	out := buff.String()

	if out != expectedOut {
		t.Fatalf("logged output did not match: got %q expected %q", out, expectedOut)
	}
}
