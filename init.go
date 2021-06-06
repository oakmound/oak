package oak

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/timing"
)

var (
	zeroPoint = image.Point{0, 0}
)

// Init initializes the oak engine.
// It spawns off an event loop of several goroutines
// and loops through scenes after initialization.
func (c *Controller) Init(firstScene string, configOptions ...ConfigOption) error {

	var err error
	c.config, err = NewConfig(configOptions...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// if c.config.Screen.TargetWidth != 0 && c.config.Screen.TargetHeight != 0 {
	// 	w, h := driver.MonitorSize()
	// 	if w != 0 || h != 0 {
	// 		// Todo: Modify conf.Screen.Scale
	// 	}
	// }

	lvl, err := dlog.ParseDebugLevel(c.config.Debug.Level)
	if err != nil {
		return fmt.Errorf("failed to parse debug config: %w", err)
	}
	dlog.SetFilter(func(msg string) bool {
		return strings.Contains(msg, c.config.Debug.Filter)
	})
	err = dlog.SetDebugLevel(lvl)
	if err != nil {
		return err
	}
	err = oakerr.SetLanguageString(c.config.Language)
	if err != nil {
		return err
	}

	// TODO: languages
	dlog.Info("Oak Init Start")

	c.ScreenWidth = c.config.Screen.Width
	c.ScreenHeight = c.config.Screen.Height
	c.FrameRate = c.config.FrameRate
	c.DrawFrameRate = c.config.DrawFrameRate
	c.IdleDrawFrameRate = c.config.IdleDrawFrameRate
	// assume we are in focus on window creation
	c.inFocus = true

	c.DrawTicker = time.NewTicker(timing.FPSToFrameDelay(c.DrawFrameRate))

	wd, _ := os.Getwd()

	render.SetFontDefaults(wd, c.config.Assets.AssetPath, c.config.Assets.FontPath,
		c.config.Font.Hinting, c.config.Font.Color, c.config.Font.File, c.config.Font.Size,
		c.config.Font.DPI)

	if c.config.TrackInputChanges {
		trackJoystickChanges(c.logicHandler)
	}
	if c.config.EventRefreshRate != 0 {
		c.logicHandler.SetRefreshRate(time.Duration(c.config.EventRefreshRate))
	}

	if !c.config.SkipRNGSeed {
		// seed math/rand with time.Now, useful for minimal examples
		//that would tend to forget to do this.
		rand.Seed(time.Now().UTC().UnixNano())
	}

	imageDir := filepath.Join(wd,
		c.config.Assets.AssetPath,
		c.config.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		c.config.Assets.AssetPath,
		c.config.Assets.AudioPath)

	// TODO: languages
	dlog.Info("Init Scene Loop")
	go c.sceneLoop(firstScene, c.config.TrackInputChanges)
	dlog.Info("Init asset load")
	render.SetAssetPaths(imageDir)
	go c.loadAssets(imageDir, audioDir)
	if c.config.EnableDebugConsole {
		dlog.Info("Init Console")
		go c.debugConsole(os.Stdin)
	}
	dlog.Info("Init Main Driver")
	c.Driver(c.lifecycleLoop)
	return c.exitError
}
