package oak

import (
	"image"
	"os"
	"path/filepath"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"golang.org/x/exp/shiny/driver"
)

var (

	// The init channel communicates between
	// initializing goroutines for when significant
	// steps in initialization have been reached
	// initCh = make(chan bool)
	// currently unused

	//
	transitionCh = make(chan bool)

	// The Scene channel receives a signal
	// when a scene's .loop() function should
	// be called.
	sceneCh = make(chan bool)

	// The skip scene channel receives a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh = make(chan bool)

	// The quit channel receives a signal when
	// the program should stop.
	quitCh = make(chan bool)

	// The draw channel receives a signal when
	// drawing should cease (or resume)
	drawCh = make(chan bool)

	// The debug reset channel represents
	// when the debug console should forget the
	// commands that have been sent to it.
	debugResetCh = make(chan bool)

	// The viewport channel controls when new
	// viewport positions should be drawn
	viewportCh = make(chan [2]int)

	debugResetInProgress bool

	// ScreenWidth is the width of the screen
	ScreenWidth int
	// ScreenHeight is the height of the screen
	ScreenHeight int

	// FrameRate is the current logical frame rate.
	// Changing this won't directly effect frame rate, that
	// requires changing the LogicTicker, but it will take
	// effect next scene
	FrameRate int

	// DrawFrameRate is the unused equivalent to FrameRate
	DrawFrameRate int

	eb *event.Bus

	// GlobalFirstScene is returned by the first
	// loading scene
	globalFirstScene string

	// CurrentScene is the scene currently running in oak
	CurrentScene string

	zeroPoint = image.Point{0, 0}
)

// Init initializes the oak engine.
// It spawns off an event loop of several goroutines
// and loops through scenes after initialization.
func Init(firstScene string) {
	dlog.CreateLogFile()

	initConf()

	// Set variables from conf file
	dlog.SetStringDebugLevel(conf.Debug.Level)
	dlog.SetDebugFilter(conf.Debug.Filter)

	dlog.Info("Oak Init Start")

	ScreenWidth = conf.Screen.Width
	ScreenHeight = conf.Screen.Height
	FrameRate = conf.FrameRate
	DrawFrameRate = conf.DrawFrameRate

	wd, _ := os.Getwd()

	render.SetFontDefaults(wd, conf.Assets.AssetPath, conf.Assets.FontPath,
		conf.Font.Hinting, conf.Font.Color, conf.Font.File, conf.Font.Size,
		conf.Font.DPI)
	// END of loading variables from configuration

	SeedRNG(DefaultSeed)

	imageDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	dlog.Info("Init Scene Loop")
	go sceneLoop(firstScene)
	dlog.Info("Init asset load")
	go loadAssets(imageDir, audioDir)
	dlog.Info("Init Console")
	go debugConsole(debugResetCh, skipSceneCh)
	dlog.Info("Init Main Driver")
	driver.Main(lifecycleLoop)
}
