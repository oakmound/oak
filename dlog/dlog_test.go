package dlog_test

import (
	"fmt"
	"testing"

	"github.com/oakmound/oak/v3/dlog"
)

func TestErrorCheck(t *testing.T) {
	called := false
	dlog.Error = func(...interface{}) {
		called = true
	}
	dlog.ErrorCheck(nil)
	if called {
		t.Fatal("error should not have been called on nil error")
	}
	dlog.ErrorCheck(fmt.Errorf("err"))
	if !called {
		t.Fatal("error should have been called on real error")
	}
	dlog.Error = func(...interface{}) {}
}

func TestParseDebugLevel(t *testing.T) {
	type testCase struct {
		in          string
		outLevel    dlog.Level
		outErrors   bool
		outErrorStr string
	}
	tcs := []testCase{
		{
			in:       "info",
			outLevel: dlog.INFO,
		}, {
			in:       "InFo",
			outLevel: dlog.INFO,
		}, {
			in:       "verbose",
			outLevel: dlog.VERBOSE,
		}, {
			in:       "ERROR",
			outLevel: dlog.ERROR,
		}, {
			in:       "warN",
			outLevel: dlog.WARN,
		}, {
			in:       "none",
			outLevel: dlog.NONE,
		}, {
			in:          "other",
			outErrors:   true,
			outErrorStr: "parsing dlog level of \"OTHER\" failed",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			lvl, err := dlog.ParseDebugLevel(tc.in)
			if tc.outErrors {
				if err == nil {
					t.Fatal("expected error")
				}
				if tc.outErrorStr != "" {
					if tc.outErrorStr != err.Error() {
						t.Fatalf("error did not match: got %v expected %v", err.Error(), tc.outErrorStr)
					}
				}
				return
			}
			if lvl != tc.outLevel {
				t.Fatalf("level did not match: got %v expected %v", lvl, tc.outLevel)
			}
		})
	}
}
