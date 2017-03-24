package oak

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	tmpConf oakConfig
	err     error
	conf    = oakConfig{
		Assets{"assets/", "audio/", "images/", "font/"},
		Debug{"", "ERROR"},
		Screen{480, 640},
		Font{"none", 12.0, 72.0, "luxisr.ttf", "white"},
		World{4000, 4000},
		60,
		false,
		"English",
		"Oak Window",
	}
)

type oakConfig struct {
	Assets    Assets `json:"assets"`
	Debug     Debug  `json:"debug"`
	Screen    Screen `json:"screen"`
	Font      Font   `json:"font"`
	World     World  `json:"world"`
	FrameRate int    `json:"frameRate"`
	ShowFPS   bool   `json:showFPS`
	Language  string `json:"language"`
	Title     string `json:"title"`
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
type World struct {
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

	tmpConf, err = loadOakConfig(filepath.Join(wd, fileName))
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

	if tmpConf.World.Width != 0 {
		conf.World.Width = tmpConf.World.Width
	}
	if tmpConf.World.Height != 0 {
		conf.World.Height = tmpConf.World.Height
	}

	if tmpConf.FrameRate != 0 {
		conf.FrameRate = tmpConf.FrameRate
	}

	if tmpConf.ShowFPS != false {
		conf.ShowFPS = tmpConf.ShowFPS
	}

	if tmpConf.Language != "" {
		conf.Language = tmpConf.Language
	}

	if tmpConf.Title != "" {
		conf.Title = tmpConf.Title
	}

	dlog.Error(conf)

	return err
}

func loadOakConfig(fileName string) (oakConfig, error) {

	dlog.Error("Loading config:", fileName)

	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		dlog.Error(err)
		return oakConfig{}, err
	}
	var config oakConfig
	json.Unmarshal(confFile, &config)
	dlog.Error(config)

	return config, nil
}

func (oc *oakConfig) String() string {
	st := "Config:\n{"
	st += oc.Debug.String()
	st += "\n}"
	return st
}

func (d *Debug) String() string {
	st := "Debug:\n{"
	st += "Level: " + d.Level
	st += "\nFilter:" + d.Filter
	st += "\n}"
	return st
}
