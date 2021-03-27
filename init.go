package oak

import (
	"image"
	"os"
	"path/filepath"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/oakerr"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/shiny/driver"
)

var (
	zeroPoint = image.Point{0, 0}
)

// Init initializes the oak engine.
// It spawns off an event loop of several goroutines
// and loops through scenes after initialization.
func (c *Controller) Init(firstScene string) {
	dlog.SetLogger(dlog.NewLogger())
	dlog.CreateLogFile()

	initConf()

	if conf.Screen.TargetWidth != 0 && conf.Screen.TargetHeight != 0 {
		w, h := driver.MonitorSize()
		if w != 0 || h != 0 {
			// Todo: Modify conf.Screen.Scale
		}
	}

	// Set variables from conf file
	lvl, err := dlog.ParseDebugLevel(conf.Debug.Level)
	dlog.SetDebugLevel(lvl)
	// We are intentionally using the lvl value before checking error,
	// because we can only log errors through dlog itself anyway

	// We do this knowing that the default debug level when SetDebugLevel fails
	// is ERROR, so this will be recorded.
	dlog.ErrorCheck(err)
	dlog.SetDebugFilter(conf.Debug.Filter)
	oakerr.SetLanguageString(conf.Language)

	// TODO: languages
	dlog.Info("Oak Init Start")

	c.ScreenWidth = conf.Screen.Width
	c.ScreenHeight = conf.Screen.Height
	c.FrameRate = conf.FrameRate
	c.DrawFrameRate = conf.DrawFrameRate

	wd, _ := os.Getwd()

	render.SetFontDefaults(wd, conf.Assets.AssetPath, conf.Assets.FontPath,
		conf.Font.Hinting, conf.Font.Color, conf.Font.File, conf.Font.Size,
		conf.Font.DPI)

	if conf.TrackInputChanges {
		trackJoystickChanges()
	}
	if conf.EventRefreshRate != 0 {
		if cfgHandler, ok := c.logicHandler.(event.ConfigHandler); ok {
			cfgHandler.SetRefreshRate(conf.EventRefreshRate)
		}
	}
	// END of loading variables from configuration

	seedRNG()

	imageDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	// TODO: languages
	dlog.Info("Init Scene Loop")
	go c.sceneLoop(firstScene, conf.TrackInputChanges, conf.DisableDebugConsole)
	dlog.Info("Init asset load")
	render.SetAssetPaths(imageDir)
	go c.loadAssets(imageDir, audioDir)
	if !conf.DisableDebugConsole {
		dlog.Info("Init Console")
		go c.debugConsole(c.debugResetCh, c.skipSceneCh, os.Stdin)
	}
	dlog.Info("Init Main Driver")
	c.Driver(c.lifecycleLoop)
}
