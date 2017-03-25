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
	"bitbucket.org/oakmoundstudio/oak/timing"
	"fmt"
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

	ScreenWidth  int
	ScreenHeight int
	WorldWidth   int
	WorldHeight  int
	FrameRate    int

	eb *event.EventBus

	wd, _    = os.Getwd()
	imageDir string
	audioDir string

	// GlobalFirstScene is returned by the first
	// loading scene
	globalFirstScene string
	CurrentScene     string

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

	// This check is delayed till after the above lines
	// As otherwise the dlog call would crash/be unseen
	if err != nil {
		dlog.Verb(err)
	}

	ScreenWidth = conf.Screen.Width
	ScreenHeight = conf.Screen.Height
	WorldWidth = conf.World.Width
	WorldHeight = conf.World.Height
	FrameRate = conf.FrameRate

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

	collision.Init()
	mouse.Init()
	render.InitDrawHeap()
	audio.InitWinAudio()

	SeedRNG(DEFAULT_SEED)
	fmt.Println("Got past seeding")
	go LoadAssets()
	go driver.Main(lifecycleLoop)
	go DebugConsole(debugResetCh, skipSceneCh)
	fmt.Println("Got past gofuncs")

	prevScene := ""
	sceneMap[firstScene].active = true

	<-initCh

	fmt.Println("Got past initchan")
	// This is the only time oak closes a channel
	// This should probably change
	close(initCh)

	// Loop through scenes

	runEventLoop = true
	globalFirstScene = firstScene
	CurrentScene = "loading"
	result := new(SceneResult)
	dlog.Info("First Scene Start")
	drawChannel <- true
	drawChannel <- true

	for {
		ViewPos = image.Point{0, 0}
		updateScreen(0, 0)
		useViewBounds = false
		dlog.Info("~~~~~~~~~~~Scene Start~~~~~~~~~")
		go func() {
			sceneMap[CurrentScene].start(prevScene, result.NextSceneInput)
			transitionCh <- true
		}()
		sceneTransition(result)
		// Post transition, begin loading animation
		drawChannel <- true
		<-transitionCh
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
				cont = sceneMap[CurrentScene].loop()
			case <-skipSceneCh:
				cont = false
			}
		}
		dlog.Info("~~~~~~~~Scene End~~~~~~~~~~")
		prevScene = CurrentScene

		// Send a signal to stop drawing
		drawChannel <- true

		// Reset any ongoing delays
	delayLabel:
		for {
			select {
			case timing.ClearDelayCh <- true:
			default:
				break delayLabel
			}
		}
		// Reset transient portions of the engine
		event.ResetEntities()
		event.ResetEventBus()
		render.ResetDrawHeap()
		collision.Clear()
		mouse.Clear()
		render.PreDraw()

		// Todo: Add in customizable loading scene between regular scenes

		CurrentScene, result = sceneMap[CurrentScene].end()
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(SceneResult)
		}

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
