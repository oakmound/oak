package oak

import (
	"image"
	"os"
	"path/filepath"

	"bitbucket.org/oakmoundstudio/oak/audio"
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"golang.org/x/exp/shiny/driver"
)

var (

	// The init channel communicates between
	// initializing goroutines for when significant
	// steps in initialization have been reached
	initCh = make(chan bool)

	//
	transitionCh = make(chan bool)

	// The Scene channel recieves a signal
	// when a scene's .loop() function should
	// be called.
	sceneCh = make(chan bool)

	// The skip scene channel recieves a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh = make(chan bool)

	// The quit channel recieves a signal when
	// the program should stop.
	quitCh = make(chan bool)

	// The draw channel recieves a signal when
	// drawing should cease (or resume)
	drawChannel = make(chan bool)

	// The debug reset channel represents
	// when the debug console should forget the
	// commands that have been sent to it.
	debugResetCh = make(chan bool)

	// The viewport channel controls when new
	// viewport positions should be drawn
	viewportChannel = make(chan [2]int)

	runEventLoop         bool
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
// and loops through scenes after initalization.
func Init(firstScene string) {
	dlog.CreateLogFile()

	loadDefaultConf()

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

	collision.Init()
	mouse.Init()
	audio.InitWinAudio()

	SeedRNG(DefaultSeed)

	imageDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	go loadAssets(imageDir, audioDir)
	go driver.Main(lifecycleLoop)
	go debugConsole(debugResetCh, skipSceneCh)

	// Loop through scenes
	sceneLoop(firstScene)
}
