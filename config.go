package oak

import (
	"encoding/json"
	"io"

	"github.com/oakmound/oak/fileutil"

	"github.com/oakmound/oak/dlog"
)

var (
	// SetupConfig is the config struct read from at initialization time
	// when oak starts. When oak.Init() is called, the variables behind
	// SetupConfig are passed to their appropriate places in the engine, and
	// afterword the variable is unused.
	SetupConfig Config

	// These are the default settings of a project. Anything within SetupConfig
	// that is set to its zero value will not overwrite these settings.
	conf = Config{
		Assets{"assets/", "audio/", "images/", "font/"},
		Debug{"", "ERROR"},
		Screen{480, 640, 1},
		Font{"none", 12.0, 72.0, "", "white"},
		60,
		60,
		"English",
		"Oak Window",
		false,
		false,
		false,
	}
)

// Config stores initialization settings for oak.
type Config struct {
	Assets         Assets `json:"assets"`
	Debug          Debug  `json:"debug"`
	Screen         Screen `json:"screen"`
	Font           Font   `json:"font"`
	FrameRate      int    `json:"frameRate"`
	DrawFrameRate  int    `json:"drawFrameRate"`
	Language       string `json:"language"`
	Title          string `json:"title"`
	BatchLoad      bool   `json:"batchLoad"`
	GestureSupport bool   `json:"gestureSupport"`
	// DisableKeyhold is deprecated. Keyhold functionality
	// no longer has a significant performance impact and so can't
	// be disabled.
	DisableKeyhold bool `json:"disableKeyHold"`
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
	Scale  int `json:"scale"`
}

// Font is a json type storing the default font settings
type Font struct {
	Hinting string  `json:"hinting"`
	Size    float64 `json:"size"`
	DPI     float64 `json:"dpi"`
	File    string  `json:"file"`
	Color   string  `json:"color"`
}

// LoadConf loads a config file, that could exist inside
// oak's binary data storage (see fileutil), to SetupConfig
func LoadConf(filePath string) error {
	r, err := fileutil.Open(filePath)
	defer r.Close()
	if err != nil {
		dlog.Warn(err)
		return err
	}
	err = LoadConfData(r)
	dlog.Info(SetupConfig)
	return err
}

// LoadConfData takes in an io.Reader and decodes it to SetupConfig
func LoadConfData(r io.Reader) error {
	return json.NewDecoder(r).Decode(&SetupConfig)
}

func initConfAssets() {
	if SetupConfig.Assets.AssetPath != "" {
		conf.Assets.AssetPath = SetupConfig.Assets.AssetPath
	}
	if SetupConfig.Assets.ImagePath != "" {
		conf.Assets.ImagePath = SetupConfig.Assets.ImagePath
	}
	if SetupConfig.Assets.AudioPath != "" {
		conf.Assets.AudioPath = SetupConfig.Assets.AudioPath
	}
	if SetupConfig.Assets.FontPath != "" {
		conf.Assets.FontPath = SetupConfig.Assets.FontPath
	}
}

func initConfDebug() {
	if SetupConfig.Debug.Filter != "" {
		conf.Debug.Filter = SetupConfig.Debug.Filter
	}
	if SetupConfig.Debug.Level != "" {
		conf.Debug.Level = SetupConfig.Debug.Level
	}
}

func initConfScreen() {

	if SetupConfig.Screen.Width != 0 {
		conf.Screen.Width = SetupConfig.Screen.Width
	}
	if SetupConfig.Screen.Height != 0 {
		conf.Screen.Height = SetupConfig.Screen.Height
	}
	if SetupConfig.Screen.Scale != 0 {
		conf.Screen.Scale = SetupConfig.Screen.Scale
	}
}

func initConfFont() {
	if SetupConfig.Font.Hinting != "" {
		conf.Font.Hinting = SetupConfig.Font.Hinting
	}
	if SetupConfig.Font.Size != 0 {
		conf.Font.Size = SetupConfig.Font.Size
	}
	if SetupConfig.Font.DPI != 0 {
		conf.Font.DPI = SetupConfig.Font.DPI
	}
	if SetupConfig.Font.File != "" {
		conf.Font.File = SetupConfig.Font.File
	}
	if SetupConfig.Font.Color != "" {
		conf.Font.Color = SetupConfig.Font.Color
	}
}

func initConf() {

	initConfAssets()

	initConfDebug()

	initConfScreen()

	initConfFont()

	if SetupConfig.FrameRate != 0 {
		conf.FrameRate = SetupConfig.FrameRate
	}

	if SetupConfig.DrawFrameRate != 0 {
		conf.DrawFrameRate = SetupConfig.DrawFrameRate
	}

	if SetupConfig.Language != "" {
		conf.Language = SetupConfig.Language
	}

	if SetupConfig.Title != "" {
		conf.Title = SetupConfig.Title
	}

	conf.BatchLoad = SetupConfig.BatchLoad

	conf.GestureSupport = SetupConfig.GestureSupport

	conf.DisableKeyhold = SetupConfig.DisableKeyhold

	dlog.Error(conf)
}
