package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	plasticPath = "src/bitbucket.org/oakmoundstudio/plasticpiston/plastic"
)

var (
	tmpConf plasticConfig
	err     error
	conf    = plasticConfig{
		Assets{"assets/", "audio/", "images/", "font/"},
		Debug{"", "ERROR"},
		Screen{480, 640},
		Font{"none", 12.0, 72.0, "luxisr.ttf", "white"},
	}
)

type plasticConfig struct {
	Assets Assets `json:"assets"`
	Debug  Debug  `json:"debug"`
	Screen Screen `json:"screen"`
	Font   Font   `json:"font"`
}

type Assets struct {
	AssetPath string `json:"assetPath"`
	AudioPath string `json:"audioPath"`
	ImagePath string `json:"imagePath"`
	FontPath  string `json:"fontPath"`
}
type Debug struct {
	Filter string `json:"filter"`
	Level  string `json:"level"`
}
type Screen struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}
type Font struct {
	Hinting string  `json:"hinting"`
	Size    float64 `json:"size"`
	DPI     float64 `json:"dpi"`
	File    string  `json:"file"`
	Color   string  `json:"color"`
}

func LoadConf(fileName string) error {
	wd, _ := os.Getwd()
	dlog.Verb(conf)

	tmpConf, err = loadPlasticConfig(filepath.Join(wd, fileName))
	if err != nil {
		return err
	}
	return err
}

func loadDefaultConf() error {

	if tmpConf.Assets.AssetPath != "" {
		conf.Assets.AssetPath = tmpConf.Assets.AssetPath
	}
	if tmpConf.Assets.ImagePath != "" {
		conf.Assets.ImagePath = tmpConf.Assets.ImagePath
	}
	if tmpConf.Assets.AudioPath != "" {
		conf.Assets.AudioPath = tmpConf.Assets.AudioPath
	}
	if tmpConf.Assets.FontPath != "" {
		conf.Assets.FontPath = tmpConf.Assets.FontPath
	}

	if tmpConf.Debug.Filter != "" {
		conf.Debug.Filter = tmpConf.Debug.Filter
	}
	if tmpConf.Debug.Level != "" {
		conf.Debug.Level = tmpConf.Debug.Level
	}

	if tmpConf.Screen.Width != 0 {
		conf.Screen.Width = tmpConf.Screen.Width
	}
	if tmpConf.Screen.Height != 0 {
		conf.Screen.Height = tmpConf.Screen.Height
	}

	if tmpConf.Font.Hinting != "" {
		conf.Font.Hinting = tmpConf.Font.Hinting
	}
	if tmpConf.Font.Size != 0 {
		conf.Font.Size = tmpConf.Font.Size
	}
	if tmpConf.Font.DPI != 0 {
		conf.Font.DPI = tmpConf.Font.DPI
	}
	if tmpConf.Font.File != "" {
		conf.Font.File = tmpConf.Font.File
	}
	if tmpConf.Font.Color != "" {
		conf.Font.Color = tmpConf.Font.Color
	}

	dlog.Error(conf)

	return err
}

func loadPlasticConfig(fileName string) (plasticConfig, error) {

	dlog.Error("Loading config:", fileName)

	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		dlog.Error(err)
		return plasticConfig{}, err
	}
	var config plasticConfig
	json.Unmarshal(confFile, &config)
	dlog.Error(config)

	return config, nil
}
