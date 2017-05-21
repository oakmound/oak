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
	conf    = oakConfig{
		Assets{"assets/", "audio/", "images/", "font/"},
		Debug{"", "ERROR"},
		Screen{480, 640},
		Font{"none", 12.0, 72.0, "", "white"},
		60,
		60,
		false,
		"English",
		"Oak Window",
	}
)

type oakConfig struct {
	Assets        Assets `json:"assets"`
	Debug         Debug  `json:"debug"`
	Screen        Screen `json:"screen"`
	Font          Font   `json:"font"`
	FrameRate     int    `json:"frameRate"`
	DrawFrameRate int    `json:"drawFrameRate"`
	ShowFPS       bool   `json:"showFPS"`
	Language      string `json:"language"`
	Title         string `json:"title"`
}

// Assets is a json type storing paths to different asset folders
type Assets struct {
	AssetPath string `json:"assetPath"`
	AudioPath string `json:"audioPath"`
	ImagePath string `json:"imagePath"`
	FontPath  string `json:"fontPath"`
}

// Debug is a json type storing the starting debug filter and level
type Debug struct {
	Filter string `json:"filter"`
	Level  string `json:"level"`
}

// Screen is a json type storing the starting screen width and height
type Screen struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// Font is a json type storing the default font settings
type Font struct {
	Hinting string  `json:"hinting"`
	Size    float64 `json:"size"`
	DPI     float64 `json:"dpi"`
	File    string  `json:"file"`
	Color   string  `json:"color"`
}

// LoadConf loads a config file
func LoadConf(fileName string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dlog.Verb(conf)

	tmpConf, err = loadOakConfig(filepath.Join(wd, fileName))
	return
}

func loadDefaultConf() {

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

	if tmpConf.FrameRate != 0 {
		conf.FrameRate = tmpConf.FrameRate
	}

	if tmpConf.DrawFrameRate != 0 {
		conf.DrawFrameRate = tmpConf.DrawFrameRate
	}

	conf.ShowFPS = tmpConf.ShowFPS

	if tmpConf.Language != "" {
		conf.Language = tmpConf.Language
	}

	if tmpConf.Title != "" {
		conf.Title = tmpConf.Title
	}

	dlog.Error(conf)
}

func loadOakConfig(fileName string) (oakConfig, error) {

	dlog.Error("Loading config:", fileName)

	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		dlog.Error(err)
		return oakConfig{}, err
	}
	var config oakConfig
	err = json.Unmarshal(confFile, &config)
	dlog.Error(config)

	return config, err
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
