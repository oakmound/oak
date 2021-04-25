package oak

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
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
