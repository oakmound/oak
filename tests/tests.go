package tests

import (
	"testing"
)

func ExpectError(err error, t *testing.T) {
	FailIf(err == nil, "Expected error, error was nil", t)
}

func FailError(err error, t *testing.T) {
	FailIf(err != nil, err.Error(), t)
}

func FailIf(b bool, log string, t *testing.T) {
	if b {
		t.Log(log)
		t.Fail()
	}
}

func FatalIf(b bool, log string, t *testing.T) {
	if b {
		t.Fatal(log)
	}
}
