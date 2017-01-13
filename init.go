package oak

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/oakmoundstudio/oak/audio"
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

var (
	initCh          = make(chan bool)
	sceneCh         = make(chan bool)
	skipSceneCh     = make(chan bool)
	quitCh          = make(chan bool)
	drawChannel     = make(chan bool)
	debugResetCh    = make(chan bool)
	viewportChannel = make(chan [2]int)

	drawInit             bool
	runEventLoop         bool
	debugResetInProgress bool
	esc                  bool
	startupLoadComplete  bool

	ScreenWidth  int
	ScreenHeight int

	black      = color.RGBA{0x00, 0x00, 0x00, 0xff}
	imageBlack = image.Black

	worldBuffer screen.Buffer
	winBuffer   screen.Buffer
	eb          *event.EventBus
	sscreen     screen.Screen

	wd, _    = os.Getwd()
	imageDir string
	audioDir string

	globalFirstScene string
	scene            string

	zeroPoint = image.Point{0, 0}
)

// Init initializes the oak engine.
// It spawns off an event loop of several goroutines
// and loops through scenes after initalization.
func Init(firstScene string) {
	dlog.CreateLogFile()

	err := loadDefaultConf()

	// Set variables from conf file
	dlog.SetStringDebugLevel(conf.Debug.Level)
	dlog.SetDebugFilter(conf.Debug.Filter)

	if err != nil {
		dlog.Verb(err)
	}

	ScreenWidth = conf.Screen.Width
	ScreenHeight = conf.Screen.Height

	imageDir = filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir = filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	render.SetFontDefaults(wd, conf.Assets.AssetPath, conf.Assets.FontPath,
		conf.Font.Hinting, conf.Font.Color, conf.Font.File, conf.Font.Size,
		conf.Font.DPI)
	// END of loading variables from configuration

	// Init various engine pieces
	collision.Init()
	mouse.Init()
	render.InitDrawHeap()
	audio.InitWinAudio()

	// Seed the rng
	curSeed := time.Now().UTC().UnixNano()
	//curSeed = 1471104995917281000
	rand.Seed(curSeed)
	dlog.Info("The seed is:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~\nTHE SEED IS:", curSeed, "\n~~~~~~~~~~~~~~~\n")

	// Load in assets
	go func() {
		err = render.BatchLoad(imageDir)
		if err != nil {
			dlog.Error(err)
			return
		}

		err = audio.BatchLoad(audioDir)
		if err != nil {
			dlog.Error(err)
		}

		startupLoadComplete = true
	}()

	// Spawn off event loop goroutines
	go driver.Main(lifecycleLoop)

	go DebugConsole(debugResetCh, skipSceneCh)

	prevScene := ""
	sceneMap[firstScene].active = true

	<-initCh
	close(initCh)

	// Loop through scenes
	runEventLoop = true
	globalFirstScene = firstScene
	scene = "loading"
	var data interface{}
	dlog.Info("First Scene Start")
	for {
		ViewX = 0
		ViewY = 0
		useViewBounds = false
		dlog.Info("~~~~~~~~~~~Scene Start~~~~~~~~~")
		sceneMap[scene].start(prevScene, data)
		// Send a signal to resume (or begin) drawing
		drawChannel <- true

		cont := true
		for cont {
			select {
			// The quit channel represents a signal
			// for the engine to stop.
			case <-quitCh:
				return
			case <-sceneCh:
				cont = sceneMap[scene].loop()
			case <-skipSceneCh:
				cont = false
			}
		}
		dlog.Info("~~~~~~~~Scene End~~~~~~~~~~")
		prevScene = scene

		// Send a signal to stop drawing
		drawChannel <- true

		// Reset transient portions of the engine
		event.ResetEntities()
		event.ResetEventBus()
		render.ResetDrawHeap()
		collision.Clear()
		mouse.Clear()
		render.PreDraw()

		scene, data = sceneMap[scene].end()

		eb = event.GetEventBus()
		if !debugResetInProgress {
			debugResetInProgress = true
			go func() {
				debugResetCh <- true
				debugResetInProgress = false
			}()
		}
	}
}
