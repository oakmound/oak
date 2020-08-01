package dlog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCustomLogger(t *testing.T) {
	cl := NewCustomLogger(NONE)
	require.NotNil(t, cl)
	require.Equal(t, cl.debugLevel, NONE)
	require.Equal(t, cl.FilterOverrideLevel, WARN)
}

func TestCustomLogger_GetLogLevel(t *testing.T) {
	cl := NewCustomLogger(NONE)
	level := cl.GetLogLevel()
	require.Equal(t, level, NONE)
}

func testCustomLoggerContains(t *testing.T, f func(cl *CustomLogger, buff *bytes.Buffer), s string) {
	cl := NewCustomLogger(VERBOSE)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	f(cl, buff)
	if len(s) == 0 {
		require.Equal(t, "", buff.String())
	} else {
		require.Contains(t, buff.String(), s)
	}
}

func TestCustomLogger_SetDebugFilter(t *testing.T) {
	type testCase struct {
		name     string
		fn       func(cl *CustomLogger, buff *bytes.Buffer)
		contains string
	}
	tcs := []testCase{
		{
			name: "valid regex: no match",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("20145")
			},
		}, {
			name: "valid regex: no match 2",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("myfilter")
			},
		}, {
			name: "valid regex: match",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("myfilter11111")
			},
			contains: "myfilter11111",
		}, {
			name: "invalid regex",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
			},
			contains: "could not compile filter regex",
		}, {
			name: "invalid regex: no match",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
				buff.Reset()
				cl.Verb("1423")
			},
		}, {
			name: "invalid regex: match",
			fn: func(cl *CustomLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
				buff.Reset()
				cl.Verb("([1-9)+")
			},
			contains: "([1-9)+",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testCustomLoggerContains(t, tc.fn, tc.contains)
		})
	}
}

func TestCustomLogger_SetDebugLevel(t *testing.T) {
	cl := NewCustomLogger(NONE)
	cl.SetDebugLevel(INFO)
	require.Equal(t, cl.debugLevel, INFO)

	level := cl.GetLogLevel()
	require.Equal(t, level, INFO)
}

func TestCustomLogger_CreateLogFile(t *testing.T) {

}

func TestCustomLogger_FileWrite(t *testing.T) {

}

func TestCustomLogger_Error(t *testing.T) {

}

func TestCustomLogger_Warn(t *testing.T) {

}

func TestCustomLogger_Info(t *testing.T) {

}

func TestCustomLogger_Verb(t *testing.T) {

}

func TestCustomLogger_SetWriter(t *testing.T) {

}
