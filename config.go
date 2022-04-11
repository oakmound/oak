package oak

import (
	"encoding/json"
	"io"

	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/shiny/driver"
)

// Config stores initialization settings for oak.
type Config struct {
	Driver                 Driver           `json:"-"`
	Assets                 Assets           `json:"assets"`
	Debug                  Debug            `json:"debug"`
	Screen                 Screen           `json:"screen"`
	BatchLoadOptions       BatchLoadOptions `json:"batchLoadOptions"`
	FrameRate              int              `json:"frameRate"`
	DrawFrameRate          int              `json:"drawFrameRate"`
	IdleDrawFrameRate      int              `json:"idleDrawFrameRate"`
	Language               string           `json:"language"`
	Title                  string           `json:"title"`
	BatchLoad              bool             `json:"batchLoad"`
	GestureSupport         bool             `json:"gestureSupport"`
	LoadBuiltinCommands    bool             `json:"loadBuiltinCommands"`
	TrackInputChanges      bool             `json:"trackInputChanges"`
	EnableDebugConsole     bool             `json:"enableDebugConsole"`
	TopMost                bool             `json:"topmost"`
	Borderless             bool             `json:"borderless"`
	Fullscreen             bool             `json:"fullscreen"`
	SkipRNGSeed            bool             `json:"skip_rng_seed"`
	UnlimitedDrawFrameRate bool             `json:"unlimitedDrawFrameRate"`
}

// NewConfig creates a config from a set of transformation options.
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
	c.Driver = driver.Main
	c.Assets = Assets{
		AudioPath: "assets/audio/",
		ImagePath: "assets/images/",
	}
	c.Debug = Debug{
		Level: "ERROR",
	}
	c.Screen = Screen{
		Height: 480,
		Width:  640,
		Scale:  1,
	}
	c.FrameRate = 60
	c.DrawFrameRate = 60
	c.IdleDrawFrameRate = 60
	c.Language = "English"
	c.Title = "Oak Window"
	return c
}

// Assets is a json type storing paths to different asset folders
type Assets struct {
	AudioPath string `json:"audioPath"`
	ImagePath string `json:"imagePath"`
}

// Debug is a json type storing the starting debug filter and level
type Debug struct {
	Filter string `json:"filter"`
	Level  string `json:"level"`
}

// Screen is a json type storing the starting screen width and height
type Screen struct {
	X      int     `json:"X"`
	Y      int     `json:"Y"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Scale  float64 `json:"scale"`
}

// BatchLoadOptions is a json type storing customizations for batch loading.
// These settings do not take effect unless batch load is true.
type BatchLoadOptions struct {
	BlankOutAudio    bool  `json:"blankOutAudio"`
	MaxImageFileSize int64 `json:"maxImageFileSize"`
}

// FileConfig loads a config file, that could exist inside
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

// A ConfigOption transforms a Config object.
type ConfigOption func(Config) (Config, error)

// ReaderConfig reads a Config as json from the given reader.
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
	if c2.Driver != nil {
		c.Driver = c2.Driver
	}
	if c2.Assets.AudioPath != "" {
		c.Assets.AudioPath = c2.Assets.AudioPath
	}
	if c2.Assets.ImagePath != "" {
		c.Assets.ImagePath = c2.Assets.ImagePath
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
	c.SkipRNGSeed = c2.SkipRNGSeed
	c.UnlimitedDrawFrameRate = c2.UnlimitedDrawFrameRate
	return c
}
