package dlog

import "testing"

func TestSetCustomLogger(t *testing.T) {
	type customLogger struct {
		FullLogger
	}
	cl := customLogger{}
	SetLogger(cl)

	SetLogger(&logger{})

	_, isCustom := oakLogger.(customLogger)
	if !isCustom {
		t.Fatal("custom logger should not have been overwritten by default logger")
	}

	oakLogger = nil
	fullOakLogger = nil
}
