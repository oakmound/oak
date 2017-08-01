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
		Assets{"a/", "a/", "i/", "f/"},
		Debug{"FILTER", "INFO"},
		Screen{240, 320, 2},
		Font{"hint", 20.0, 36.0, "luxisr.ttf", "green"},
		30,
		30,
		"German",
		"Some Window",
		true,
		true,
		true,
	}
	initConf()
	assert.Equal(t, SetupConfig, conf)
	// Failure to load
	err = LoadConf("nota.config")
	assert.NotNil(t, err)
}
