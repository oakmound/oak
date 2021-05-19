package dlog_test

import (
	"testing"
	"github.com/oakmound/oak/v3/dlog"
)


func TestLevelsString(t *testing.T) {
	type testCase struct {
		in  dlog.Level
		out string
	}
	tcs := []testCase{
		{
			in:  dlog.NONE,
			out: "NONE",
		}, {
			in:  dlog.ERROR,
			out: "ERROR",
		}, {
			in:  dlog.WARN,
			out: "WARN",
		}, {
			in:  dlog.INFO,
			out: "INFO",
		}, {
			in:  dlog.VERBOSE,
			out: "VERBOSE",
		}, {
			in:  dlog.Level(100),
			out: "",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.out, func(t *testing.T) {
			out := tc.in.String()
			if out != tc.out {
				t.Fatalf("mismatched output, got %v expected %v", out, tc.out)
			}
		})
	}
}
