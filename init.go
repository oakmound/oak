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
// After the configuration options have been parsed and validated, this will run concurrent
// routines drawing to an OS window or app, forwarding OS inputs to this window's configured
// event handler, and running scenes: first the predefined 'loading' scene, then firstScene
// as provided here, then scenes following commands sent to the window or returned by ending
// scenes.
func (w *Window) Init(firstScene string, configOptions ...ConfigOption) error {

	var err error
	w.config, err = NewConfig(configOptions...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

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

	if !w.config.SkipRNGSeed {
		// seed math/rand with time.Now, useful for minimal examples
		//that would tend to forget to do this.
		rand.Seed(time.Now().UTC().UnixNano())
	}

	overrideInit(w)

	go w.sceneLoop(firstScene, w.config.TrackInputChanges, w.config.BatchLoad)
	if w.config.EnableDebugConsole {
		go w.debugConsole(os.Stdin, os.Stdout)
	}
	w.Driver(w.lifecycleLoop)
	return w.exitError
}
