package oak

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/oakmound/oak/v2/fileutil"
)

// Config stores initialization settings for oak.
type Config struct {
	Assets              Assets           `json:"assets"`
	Debug               Debug            `json:"debug"`
	Screen              Screen           `json:"screen"`
	Font                Font             `json:"font"`
	BatchLoadOptions    BatchLoadOptions `json:"batchLoadOptions"`
	FrameRate           int              `json:"frameRate"`
	DrawFrameRate       int              `json:"drawFrameRate"`
	IdleDrawFrameRate   int              `json:"idleDrawFrameRate"`
	Language            string           `json:"language"`
	Title               string           `json:"title"`
	EventRefreshRate    Duration         `json:"refreshRate"`
	BatchLoad           bool             `json:"batchLoad"`
	GestureSupport      bool             `json:"gestureSupport"`
	LoadBuiltinCommands bool             `json:"loadBuiltinCommands"`
	TrackInputChanges   bool             `json:"trackInputChanges"`
	EnableDebugConsole  bool             `json:"enableDebugConsole"`
	TopMost             bool             `json:"topmost"`
	Borderless          bool             `json:"borderless"`
	Fullscreen          bool             `json:"fullscreen"`
}

// A Duration is a wrapper arouind time.Duration that allows for easier json formatting.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration type")
	}
}

func NewConfig(opts ...ConfigOption) (Config, error) {
	c := Config{}
	c = c.setDefaults()
	var err error
	for _, o := range opts {
		c, err = o(c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func (c Config) setDefaults() Config {
	c.Assets = Assets{
		AssetPath: "assets/",
		AudioPath: "audio/",
		ImagePath: "images/",
		FontPath:  "font/",
	}
	c.Debug = Debug{
		Level: "ERROR",
	}
	c.Screen = Screen{
		Height: 480,
		Width:  640,
		Scale:  1,
	}
	c.Font = Font{
		Hinting: "none",
		Size:    12.0,
		DPI:     72.0,
		File:    "",
		Color:   "white",
	}
	c.FrameRate = 60
	c.DrawFrameRate = 60
	c.IdleDrawFrameRate = 60
	c.Language = "English"
	c.Title = "Oak Window"
	c.EventRefreshRate = Duration(50 * time.Millisecond)
	return c
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
	X      int `json:"X"`
	Y      int `json:"Y"`
	Height int `json:"height"`
	Width  int `json:"width"`
	Scale  int `json:"scale"`
	// Target sets the expected dimensions of the monitor the game will be opened on, in pixels.
	// If Fullscreen is false, then a scaling will be applied to correct the game screen size to be
	// appropriate for the Target size. If no TargetWidth or Height is provided, scaling will not
	// be adjusted.
	TargetWidth  int `json:"targetHeight"`
	TargetHeight int `json:"targetWidth"`
}

// Font is a json type storing the default font settings
type Font struct {
	Hinting string  `json:"hinting"`
	Size    float64 `json:"size"`
	DPI     float64 `json:"dpi"`
	File    string  `json:"file"`
	Color   string  `json:"color"`
}

// BatchLoadOptions is a json type storing customizations for batch loading.
// These settings do not take effect unless batch load is true.
type BatchLoadOptions struct {
	BlankOutAudio    bool  `json:"blankOutAudio"`
	MaxImageFileSize int64 `json:"maxImageFileSize"`
}

// LoadConf loads a config file, that could exist inside
// oak's binary data storage (see fileutil), to SetupConfig
func FileConfig(filePath string) ConfigOption {
	return func(c Config) (Config, error) {
		r, err := fileutil.Open(filePath)
		if err != nil {
			return c, err
		}
		defer r.Close()
		return ReaderConfig(r)(c)
	}
}

type ConfigOption func(Config) (Config, error)

func ReaderConfig(r io.Reader) ConfigOption {
	return func(c Config) (Config, error) {
		c2 := Config{}
		decoder := json.NewDecoder(r)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&c2)
		if err != nil {
			return c, err
		}
		c2 = c.overwriteFrom(c2)
		return c2, nil
	}
}

func (c Config) overwriteFrom(c2 Config) Config {
	// TODO: is this the right place for these configuration pieces?
	// TODO: is there other configuration that should go here?
	if c2.Assets.AssetPath != "" {
		c.Assets.AssetPath = c2.Assets.AssetPath
	}
	if c2.Assets.AudioPath != "" {
		c.Assets.AudioPath = c2.Assets.AudioPath
	}
	if c2.Assets.ImagePath != "" {
		c.Assets.ImagePath = c2.Assets.ImagePath
	}
	if c2.Assets.FontPath != "" {
		c.Assets.FontPath = c2.Assets.FontPath
	}
	if c2.Debug.Filter != "" {
		c.Debug.Filter = c2.Debug.Filter
	}
	if c2.Debug.Level != "" {
		c.Debug.Level = c2.Debug.Level
	}
	if c2.Screen.X != 0 {
		c.Screen.X = c2.Screen.X
	}
	if c2.Screen.Y != 0 {
		c.Screen.Y = c2.Screen.Y
	}
	if c2.Screen.Height != 0 {
		c.Screen.Height = c2.Screen.Height
	}
	if c2.Screen.Width != 0 {
		c.Screen.Width = c2.Screen.Width
	}
	if c2.Screen.Scale != 0 {
		c.Screen.Scale = c2.Screen.Scale
	}
	if c2.Screen.TargetWidth != 0 {
		c.Screen.TargetWidth = c2.Screen.TargetWidth
	}
	if c2.Screen.TargetHeight != 0 {
		c.Screen.TargetHeight = c2.Screen.TargetHeight
	}
	if c2.Font.Hinting != "" {
		c.Font.Hinting = c2.Font.Hinting
	}
	if c2.Font.Size != 0 {
		c.Font.Size = c2.Font.Size
	}
	if c2.Font.DPI != 0 {
		c.Font.DPI = c2.Font.DPI
	}
	if c2.Font.File != "" {
		c.Font.File = c2.Font.File
	}
	if c2.Font.Color != "" {
		c.Font.Color = c2.Font.Color
	}
	c.BatchLoadOptions.BlankOutAudio = c2.BatchLoadOptions.BlankOutAudio
	if c2.BatchLoadOptions.MaxImageFileSize != 0 {
		c.BatchLoadOptions.MaxImageFileSize = c2.BatchLoadOptions.MaxImageFileSize
	}
	if c2.FrameRate != 0 {
		c.FrameRate = c2.FrameRate
	}
	if c2.DrawFrameRate != 0 {
		c.DrawFrameRate = c2.DrawFrameRate
	}
	if c2.IdleDrawFrameRate != 0 {
		c.IdleDrawFrameRate = c2.IdleDrawFrameRate
	}
	if c2.Language != "" {
		c.Language = c2.Language
	}
	if c2.Title != "" {
		c.Title = c2.Title
	}
	if c2.EventRefreshRate != 0 {
		c.EventRefreshRate = c2.EventRefreshRate
	}
	// Booleans can be directly overwritten-- all booleans in a Config
	// default to false, if they were unset they will stay false.
	c.BatchLoad = c2.BatchLoad
	c.GestureSupport = c2.GestureSupport
	c.LoadBuiltinCommands = c2.LoadBuiltinCommands
	c.TrackInputChanges = c2.TrackInputChanges
	c.EnableDebugConsole = c2.EnableDebugConsole
	c.TopMost = c2.TopMost
	c.Borderless = c2.Borderless
	c.Fullscreen = c2.Fullscreen
	return c
}
