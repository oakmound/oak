package dlog_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/oakmound/oak/v3/dlog"
)

func TestErrorCheck(t *testing.T) {
	buff := &bytes.Buffer{}
	dlog.DefaultLogger.SetOutput(buff)
	dlog.ErrorCheck(nil)
	if buff.Len() != 0 {
		t.Fatal("error should not have been called on nil error")
	}
	dlog.ErrorCheck(fmt.Errorf("err"))
	if buff.Len() == 0 {
		t.Fatal("error should have been called on real error")
	}
	dlog.DefaultLogger.SetOutput(os.Stdout)
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
			in:       "none",
			outLevel: dlog.NONE,
		}, {
			in:          "other",
			outErrors:   true,
			outErrorStr: "invalid input: level",
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
