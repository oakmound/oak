package dlog

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRegexLogger(t *testing.T) {
	cl := NewRegexLogger(NONE)
	if cl.debugLevel != NONE {
		t.Fatalf("expected %v debug level, got %v", NONE, cl.debugLevel)
	}
	if cl.FilterOverrideLevel != WARN {
		t.Fatalf("expected %v override level, got %v", WARN, cl.FilterOverrideLevel)
	}
}

func TestRegexLogger_GetLogLevel(t *testing.T) {
	cl := NewRegexLogger(NONE)
	level := cl.GetLogLevel()
	if level != NONE {
		t.Fatalf("expected %v debug level, got %v", NONE, level)
	}
}

func testRegexLoggerContains(t *testing.T, f func(cl *RegexLogger, buff *bytes.Buffer), s string) {
	cl := NewRegexLogger(VERBOSE)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	f(cl, buff)
	if len(s) == 0 {
		bs := buff.String()
		if bs != "" {
			t.Fatalf("expected empty string, got %v", bs)
		}
	} else {
		bs := buff.String()
		if !strings.Contains(bs, s) {
			t.Fatalf("expected buffer to contain %v, got %v", s, bs)
		}
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
	if cl.debugLevel != INFO {
		t.Fatalf("expected %v debug level, got %v", INFO, cl.debugLevel)
	}

	level := cl.GetLogLevel()
	if level != INFO {
		t.Fatalf("expected %v debug level, got %v", INFO, level)
	}
}

func TestRegexLogger_CreateLogFile(t *testing.T) {
	cl := NewRegexLogger(NONE)
	cl.CreateLogFile()
}

func TestRegexLogger_FileWrite(t *testing.T) {
	cl := NewRegexLogger(NONE)
	buff := bytes.NewBuffer([]byte{})
	cl.FileWrite("whoops")
	cl.file = buff
	cl.FileWrite("test")
	bs := buff.String()
	if strings.Contains(bs, "whoops") {
		t.Fatalf("expected buffer not to contain %v, was %v", "whoops", bs)
	}
	if !strings.Contains(bs, "test") {
		t.Fatalf("expected buffer to contain %v, was %v", "test", bs)
	}
}

func TestRegexLogger_Error(t *testing.T) {
	cl := NewRegexLogger(ERROR)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	bs := buff.String()
	if strings.Contains(bs, "verbose") {
		t.Fatalf("expected buffer not to contain %v, was %v", "verbose", bs)
	}
	if strings.Contains(bs, "info") {
		t.Fatalf("expected buffer not to contain %v, was %v", "info", bs)
	}
	if strings.Contains(bs, "warn") {
		t.Fatalf("expected buffer not to contain %v, was %v", "warn", bs)
	}
	if !strings.Contains(bs, "error") {
		t.Fatalf("expected buffer to contain %v, was %v", "error", bs)
	}
}

func TestRegexLogger_Warn(t *testing.T) {
	cl := NewRegexLogger(WARN)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	bs := buff.String()
	if strings.Contains(bs, "verbose") {
		t.Fatalf("expected buffer not to contain %v, was %v", "verbose", bs)
	}
	if strings.Contains(bs, "info") {
		t.Fatalf("expected buffer not to contain %v, was %v", "info", bs)
	}
	if !strings.Contains(bs, "warn") {
		t.Fatalf("expected buffer to contain %v, was %v", "warn", bs)
	}
	if !strings.Contains(bs, "error") {
		t.Fatalf("expected buffer to contain %v, was %v", "error", bs)
	}
}

func TestRegexLogger_Info(t *testing.T) {
	cl := NewRegexLogger(INFO)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	bs := buff.String()
	if strings.Contains(bs, "verbose") {
		t.Fatalf("expected buffer not to contain %v, was %v", "verbose", bs)
	}
	if !strings.Contains(bs, "info") {
		t.Fatalf("expected buffer to contain %v, was %v", "info", bs)
	}
	if !strings.Contains(bs, "warn") {
		t.Fatalf("expected buffer to contain %v, was %v", "warn", bs)
	}
	if !strings.Contains(bs, "error") {
		t.Fatalf("expected buffer to contain %v, was %v", "error", bs)
	}
}

func TestRegexLogger_Verb(t *testing.T) {
	cl := NewRegexLogger(VERBOSE)
	buff := bytes.NewBuffer([]byte{})
	cl.SetWriter(buff)
	cl.Verb("verbose")
	cl.Info("info")
	cl.Warn("warn")
	cl.Error("error")
	bs := buff.String()
	if !strings.Contains(bs, "verbose") {
		t.Fatalf("expected buffer to contain %v, was %v", "verbose", bs)
	}
	if !strings.Contains(bs, "info") {
		t.Fatalf("expected buffer to contain %v, was %v", "info", bs)
	}
	if !strings.Contains(bs, "warn") {
		t.Fatalf("expected buffer to contain %v, was %v", "warn", bs)
	}
	if !strings.Contains(bs, "error") {
		t.Fatalf("expected buffer to contain %v, was %v", "error", bs)
	}
}

func TestRegexLogger_SetWriter(t *testing.T) {
	cl := NewRegexLogger(VERBOSE)
	err := cl.SetWriter(nil)
	if err == nil {
		t.Fatalf("expected setWriter(nil) to error")
	}
	err = cl.SetWriter(bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("expected setWriter([]byte) not to error: %v", err)
	}
}
