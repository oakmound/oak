package oak

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfigFileMatchesEmptyConfig(t *testing.T) {
	c1, err := NewConfig(func(c Config) (Config, error) {
		// clear out defaults
		return Config{}, nil
	}, FileConfig(filepath.Join("testdata", "default.config")))
	if err != nil {
		t.Fatalf("failed to load default.config (1): %v", err)
	}
	c2, _ := NewConfig()
	if c1 != c2 {
		t.Fatalf("config from file did not match default config (1)")
	}
	f, err := os.Open(filepath.Join("testdata", "default.config"))
	if err != nil {
		t.Fatalf("failed to load deafult.config (1)")
	}
	c3, err := NewConfig(func(c Config) (Config, error) {
		// clear out defaults
		return Config{}, nil
	}, ReaderConfig(f))
	if err != nil {
		t.Fatalf("failed to load config data from reader")
	}
	if c2 != c3 {
		t.Fatalf("config from reader did not match default config (1)")
	}
}

func TestConfig_overwriteFrom(t *testing.T) {
	// coverage test
	c2 := Config{
		Debug: Debug{
			Filter: "filter",
		},
		Screen: Screen{
			X:            1,
			Y:            1,
			TargetWidth:  1,
			TargetHeight: 1,
		},
		BatchLoadOptions: BatchLoadOptions{
			MaxImageFileSize: 10000,
		},
	}
	c1 := Config{}
	c1.overwriteFrom(c2)
}

func TestFileConfigBadFile(t *testing.T) {
	_, err := NewConfig(FileConfig("badpath"))
	if err == nil {
		t.Fatalf("expected error loading bad file")
	}
	// This error is an stdlib error, not ours, so we don't care
	// about its type
}

func TestReaderConfigBadJSON(t *testing.T) {
	b := bytes.NewBuffer([]byte("this isn't json"))
	_, err := NewConfig(ReaderConfig(b))
	if err == nil {
		t.Fatalf("expected error loading bad file")
	}
	// This error is an stdlib error, not ours, so we don't care
	// about its type
}

func TestDuration_HappyPath(t *testing.T) {
	d := Duration(time.Second)
	marshalled, err := d.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal duration failed: %v", err)
	}
	d2 := new(Duration)
	err = d2.UnmarshalJSON(marshalled)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	marshalled2, err := d2.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal duration 2 failed: %v", err)
	}
	if !bytes.Equal(marshalled, marshalled2) {
		t.Fatalf("marshals not equal: %v vs %v", string(marshalled), string(marshalled2))
	}
}

func TestDuration_UnmarshalJSON_Float(t *testing.T) {
	f := []byte("10.0")
	d2 := new(Duration)
	err := d2.UnmarshalJSON(f)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
}

func TestDuration_UnmarshalJSON_Boolean(t *testing.T) {
	f := []byte("false")
	d2 := new(Duration)
	err := d2.UnmarshalJSON(f)
	if err == nil {
		t.Fatalf("expected failure in unmarshal")
	}
}

func TestDuration_UnmarshalJSON_BadString(t *testing.T) {
	f := []byte("\"10mmmm\"")
	d2 := new(Duration)
	err := d2.UnmarshalJSON(f)
	if err == nil {
		t.Fatalf("expected failure in unmarshal")
	}
}

func TestDuration_UnmarshalJSON_BadJSON(t *testing.T) {
	f := []byte("\"1mm")
	d2 := new(Duration)
	err := d2.UnmarshalJSON(f)
	if err == nil {
		t.Fatalf("expected failure in unmarshal")
	}
}
