package oak

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/oakerr"
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

	w.ScreenWidth = w.config.Screen.Width
	w.ScreenHeight = w.config.Screen.Height
	w.FrameRate = w.config.FrameRate
	w.DrawFrameRate = w.config.DrawFrameRate
	w.IdleDrawFrameRate = w.config.IdleDrawFrameRate
	// assume we are in focus on window creation
	w.inFocus = true
	w.Driver = w.config.Driver

	w.DrawTicker = time.NewTicker(timing.FPSToFrameDelay(w.DrawFrameRate))

	if w.config.TrackInputChanges {
		trackJoystickChanges(w.eventHandler)
	}
	if w.config.EventRefreshRate != 0 {
		w.eventHandler.SetRefreshRate(time.Duration(w.config.EventRefreshRate))
	}

	if !w.config.SkipRNGSeed {
		// seed math/rand with time.Now, useful for minimal examples
		//that would tend to forget to do this.
		rand.Seed(time.Now().UTC().UnixNano())
	}

	overrideInit(w)

	go w.sceneLoop(firstScene, w.config.TrackInputChanges)
	if w.config.BatchLoad {
		w.startupLoading = true
		go func() {
			w.loadAssets(w.config.Assets.ImagePath, w.config.Assets.AudioPath)
			w.endLoad()
		}()
	}
	if w.config.EnableDebugConsole {
		go w.debugConsole(os.Stdin, os.Stdout)
	}
	w.Driver(w.lifecycleLoop)
	return w.exitError
}
