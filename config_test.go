package oak

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	err := LoadConf(filepath.Join("testdata", "default.config"))
	if err != nil {
		t.Fatalf("failed to load deafult.config (1): %v", err)
	}
	if conf != SetupConfig {
		t.Fatalf("load conf did not match default config (1)")
	}
	f, err := os.Open(filepath.Join("testdata", "default.config"))
	if err != nil {
		t.Fatalf("failed to load deafult.config (1)")
	}
	err = LoadConfData(f)
	if err != nil {
		t.Fatalf("failed to load config data from file")
	}
	SetupConfig = Config{
		Assets:              Assets{"a/", "a/", "i/", "f/"},
		Debug:               Debug{"FILTER", "INFO"},
		Screen:              Screen{0, 0, 240, 320, 2, 0, 0},
		Font:                Font{"hint", 20.0, 36.0, "luxisr.ttf", "green"},
		FrameRate:           30,
		DrawFrameRate:       30,
		Language:            "German",
		Title:               "Some Window",
		BatchLoad:           true,
		GestureSupport:      true,
		LoadBuiltinCommands: true,
	}
	initConf()
	if conf != SetupConfig {
		t.Fatalf("load conf did not match default config (2)")
	}
	// Failure to load
	err = LoadConf("nota.config")
	if err == nil {
		t.Fatalf("loading bad config file did not fail")
	}
}
