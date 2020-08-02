package dlog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRegexLogger(t *testing.T) {
	cl := NewRegexLogger(NONE)
	require.NotNil(t, cl)
	require.Equal(t, cl.debugLevel, NONE)
	require.Equal(t, cl.FilterOverrideLevel, WARN)
}

func TestRegexLogger_GetLogLevel(t *testing.T) {
	cl := NewRegexLogger(NONE)
	level := cl.GetLogLevel()
	require.Equal(t, level, NONE)
}

func testRegexLoggerContains(t *testing.T, f func(cl *RegexLogger, buff *bytes.Buffer), s string) {
	cl := NewRegexLogger(VERBOSE)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	f(cl, buff)
	if len(s) == 0 {
		require.Equal(t, "", buff.String())
	} else {
		require.Contains(t, buff.String(), s)
	}
}

func TestRegexLogger_SetDebugFilter(t *testing.T) {
	type testCase struct {
		name     string
		fn       func(cl *RegexLogger, buff *bytes.Buffer)
		contains string
	}
	tcs := []testCase{
		{
			name: "valid regex: no match",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("20145")
			},
		}, {
			name: "valid regex: no match 2",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("myfilter")
			},
		}, {
			name: "valid regex: match",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("myfilter(1)+")
				cl.Verb("myfilter11111")
			},
			contains: "myfilter11111",
		}, {
			name: "invalid regex",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
			},
			contains: "could not compile filter regex",
		}, {
			name: "invalid regex: no match",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
				buff.Reset()
				cl.Verb("1423")
			},
		}, {
			name: "invalid regex: match",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("([1-9)+")
				buff.Reset()
				cl.Verb("([1-9)+")
			},
			contains: "([1-9)+",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testRegexLoggerContains(t, tc.fn, tc.contains)
		})
	}
}

func TestRegexLogger_FilterOverrideLevel(t *testing.T) {
	type testCase struct {
		name     string
		fn       func(cl *RegexLogger, buff *bytes.Buffer)
		contains string
	}
	tcs := []testCase{
		{
			name: "default override: too low does not emit",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("filter1")
				cl.Verb("does not contain")
			},
		}, {
			name: "default override: emits",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("filter1")
				cl.Warn("does not contain")
			},
			contains: "does not contain",
		}, {
			name: "custom override: too low does not emit",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("filter1")
				cl.FilterOverrideLevel = ERROR
				cl.Warn("does not contain")
			},
		}, {
			name: "custom override: emits",
			fn: func(cl *RegexLogger, buff *bytes.Buffer) {
				cl.SetDebugFilter("filter1")
				cl.FilterOverrideLevel = ERROR
				cl.Error("does not contain")
			},
			contains: "does not contain",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testRegexLoggerContains(t, tc.fn, tc.contains)
		})
	}
}

func TestRegexLogger_SetDebugLevel(t *testing.T) {
	cl := NewRegexLogger(NONE)
	cl.SetDebugLevel(INFO)
	require.Equal(t, cl.debugLevel, INFO)

	level := cl.GetLogLevel()
	require.Equal(t, level, INFO)
}

func TestRegexLogger_CreateLogFile(t *testing.T) {
	t.Skip("not worth mocking os")
}

func TestRegexLogger_FileWrite(t *testing.T) {
	cl := NewRegexLogger(NONE)
	buff := bytes.NewBuffer([]byte{})
	cl.file = buff
	cl.FileWrite("test")
	require.Contains(t, buff.String(), "test")
}

func TestRegexLogger_Error(t *testing.T) {
	cl := NewRegexLogger(ERROR)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	logged := buff.String()
	require.NotContains(t, logged, "verbose")
	require.NotContains(t, logged, "info")
	require.NotContains(t, logged, "warn")
	require.Contains(t, logged, "error")
}

func TestRegexLogger_Warn(t *testing.T) {
	cl := NewRegexLogger(WARN)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	logged := buff.String()
	require.NotContains(t, logged, "verbose")
	require.NotContains(t, logged, "info")
	require.Contains(t, logged, "warn")
	require.Contains(t, logged, "error")
}

func TestRegexLogger_Info(t *testing.T) {
	cl := NewRegexLogger(INFO)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	logged := buff.String()
	require.NotContains(t, logged, "verbose")
	require.Contains(t, logged, "info")
	require.Contains(t, logged, "warn")
	require.Contains(t, logged, "error")
}

func TestRegexLogger_Verb(t *testing.T) {
	cl := NewRegexLogger(VERBOSE)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	logged := buff.String()
	require.Contains(t, logged, "verbose")
	require.Contains(t, logged, "info")
	require.Contains(t, logged, "warn")
	require.Contains(t, logged, "error")
}

func TestRegexLogger_SetWriter(t *testing.T) {
	cl := NewRegexLogger(VERBOSE)
	err := cl.SetWriter(nil)
	require.NotNil(t, err)
	err = cl.SetWriter(bytes.NewBuffer([]byte{}))
	require.Nil(t, err)
}
