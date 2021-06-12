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
func (w *Window) Init(firstScene string, configOptions ...ConfigOption) error {

	var err error
	w.config, err = NewConfig(configOptions...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// if c.config.Screen.TargetWidth != 0 && c.config.Screen.TargetHeight != 0 {
	// 	w, h := driver.MonitorSize()
	// 	if w != 0 || h != 0 {
	// 		// Todo: Modify conf.Screen.Scale
	// 	}
	// }

	lvl, err := dlog.ParseDebugLevel(w.config.Debug.Level)
	if err != nil {
		return fmt.Errorf("failed to parse debug config: %w", err)
	}
	dlog.SetFilter(func(msg string) bool {
		return strings.Contains(msg, w.config.Debug.Filter)
	})
	err = dlog.SetLogLevel(lvl)
	if err != nil {
		return err
	}
	err = oakerr.SetLanguageString(w.config.Language)
	if err != nil {
		return err
	}

	// TODO: languages
	dlog.Info("Oak Init Start")

	w.ScreenWidth = w.config.Screen.Width
	w.ScreenHeight = w.config.Screen.Height
	w.FrameRate = w.config.FrameRate
	w.DrawFrameRate = w.config.DrawFrameRate
	w.IdleDrawFrameRate = w.config.IdleDrawFrameRate
	// assume we are in focus on window creation
	w.inFocus = true

	w.DrawTicker = time.NewTicker(timing.FPSToFrameDelay(w.DrawFrameRate))

	wd, _ := os.Getwd()

	render.SetFontDefaults(wd, w.config.Assets.AssetPath, w.config.Assets.FontPath,
		w.config.Font.Hinting, w.config.Font.Color, w.config.Font.File, w.config.Font.Size,
		w.config.Font.DPI)

	if w.config.TrackInputChanges {
		trackJoystickChanges(w.logicHandler)
	}
	if w.config.EventRefreshRate != 0 {
		w.logicHandler.SetRefreshRate(time.Duration(w.config.EventRefreshRate))
	}

	if !w.config.SkipRNGSeed {
		// seed math/rand with time.Now, useful for minimal examples
		//that would tend to forget to do this.
		rand.Seed(time.Now().UTC().UnixNano())
	}

	imageDir := filepath.Join(wd,
		w.config.Assets.AssetPath,
		w.config.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		w.config.Assets.AssetPath,
		w.config.Assets.AudioPath)

	// TODO: languages
	go w.sceneLoop(firstScene, w.config.TrackInputChanges)
	render.SetAssetPaths(imageDir)
	go w.loadAssets(imageDir, audioDir)
	if w.config.EnableDebugConsole {
		dlog.Info("Init Console")
		go w.debugConsole(os.Stdin, os.Stdout)
	}
	w.Driver(w.lifecycleLoop)
	return w.exitError
}
