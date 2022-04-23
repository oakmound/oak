package dlog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/oakmound/oak/v4/dlog"
)

func TestLogger(t *testing.T) {
	lgr := dlog.NewLogger()

	err := lgr.SetLogLevel(-1)
	if err == nil {
		t.Fatalf("expected -1 log level to error")
	}

	lgr.SetLogLevel(dlog.VERBOSE)

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

	out := buff.String()
	expected := []string{
		"ERROR: error",
		"INFO: info",
		"VERBOSE: verb",
		"VERBOSE: foo",
	}

	lastIndexAt := 0
	for _, s := range expected {
		foundAt := strings.Index(out, s)
		if foundAt < lastIndexAt {
			t.Fatalf("did not find %v in correct order, was at index %v", s, foundAt)
		}
		lastIndexAt = foundAt
	}
}
