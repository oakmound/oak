package oak

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	err := LoadConf("default.config")
	assert.Nil(t, err)
	assert.Equal(t, conf, SetupConfig)
	f, err := os.Open("default.config")
	assert.Nil(t, err)
	err = LoadConfData(f)
	assert.Nil(t, err)
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
	assert.Equal(t, SetupConfig, conf)
	// Failure to load
	err = LoadConf("nota.config")
	assert.NotNil(t, err)
}
