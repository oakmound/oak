// Package plastic is a game engine...
package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/audio"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	pmouse "bitbucket.org/oakmoundstudio/plasticpiston/plastic/mouse"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	initCh               = make(chan bool)
	sceneCh              = make(chan bool)
	skipSceneCh          = make(chan bool)
	quitCh               = make(chan bool)
	drawChannel          = make(chan bool)
	debugResetCh         = make(chan bool)
	viewportChannel      = make(chan [2]int)
	drawInit             = false
	runEventLoop         = false
	ScreenWidth          int
	ScreenHeight         int
	press                = key.DirPress
	release              = key.DirRelease
	black                = color.RGBA{0x00, 0x00, 0x00, 0xff}
	b                    screen.Buffer
	winBuffer            screen.Buffer
	eb                   *event.EventBus
	esc                  = false
	l_debug              = false
	wd, _                = os.Getwd()
	imageDir             string
	audioDir             string
	debugResetInProgress = false
)

// Init initializes the plastic engine.
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

	imageDir = filepath.Join(filepath.Dir(wd),
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir = filepath.Join(filepath.Dir(wd),
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	render.SetFontDefaults(wd, conf.Assets.AssetPath, conf.Assets.FontPath,
		conf.Font.Hinting, conf.Font.Color, conf.Font.File, conf.Font.Size,
		conf.Font.DPI)
	// END of loading variables from configuration

	// Init various engine pieces
	collision.Init()
	pmouse.Init()
	render.InitDrawHeap()
	audio.InitWinAudio()

	// Seed the rng
	curSeed := time.Now().UTC().UnixNano()
	curSeed = 1471104995917281000
	rand.Seed(curSeed)
	dlog.Info("The seed is:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~\nTHE SEED IS:", curSeed, "\n~~~~~~~~~~~~~~~\n")

	// Load in assets
	err = render.BatchLoad(imageDir)
	if err != nil {
		dlog.Error(err)
		return
	}
	// err = audio.BatchLoad(audioDir)
	// if err != nil {
	// 	dlog.Error(err)
	// 	return
	// }

	// Spawn off event loop goroutines
	go driver.Main(eventLoop)

	go DebugConsole(debugResetCh, skipSceneCh)

	prevScene := ""
	sceneMap[firstScene].active = true

	<-initCh
	close(initCh)

	// Loop through scenes
	runEventLoop = true
	scene := firstScene
	var data interface{} = nil
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
		pmouse.Clear()
		render.PreDraw(0, nil)

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

func eventLoop(s screen.Screen) {

	// The event loop requires information about
	// the size of the world and screen that is
	// being dealt with, and so initializes it here.
	//
	// Todo: add world size to config
	b, _ = s.NewBuffer(image.Point{4096, 4096})
	winBuffer, _ = s.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	w, err := s.NewWindow(&screen.NewWindowOptions{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
	}
	defer w.Release()

	// This initialization happens here on account of font's initialization
	// requiring a buffer to draw to. Can probably change in the future.
	render.InitFont(b, winBuffer)

	eb = event.GetEventBus()

	// Todo: add frame rate to config
	frameRate := 60
	frameCh := make(chan bool)

	// This goroutine maintains a logical framerate
	go func(frameCh chan bool, frameRate int64) {
		c := time.Tick(time.Second / time.Duration(frameRate))
		for range c {
			frameCh <- true
		}
	}(frameCh, int64(frameRate))

	// Native go event handler
	go func() {
		for {
			e := w.NextEvent()
			// format := "got %#v\n"
			// if _, ok := e.(fmt.Stringer); ok {
			// 	format = "got %v\n"
			// }
			// if l_debug {
			// 	fmt.Printf(format, e)
			// }
			switch e := e.(type) {

			// We only currently respond to death lifecycle events.
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					quitCh <- true
					return
				}

			// Send key events
			//
			// Key events have two varieties:
			// The "KeyDown" and "KeyUp" events, which trigger for all keys
			// and specific "KeyDown$key", etc events which trigger only for $key.
			// The specific key that is pressed is passed as the data interface for
			// the former events, but not for the latter.
			case key.Event:
				k := GetKeyBind(e.Code.String()[4:])
				if e.Direction == press {
					fmt.Println("--------------------", e.Code.String()[4:], k)
					setDown(k)
					eb.Trigger("KeyDown", k)
					eb.Trigger("KeyDown"+k, nil)
				} else if e.Direction == release {
					setUp(k)
					eb.Trigger("KeyUp", k)
					eb.Trigger("KeyUp"+k, nil)
				}

			// Send mouse events
			//
			// Mouse events are parsed based on their button
			// and direction into an event name and then triggered:
			// 'MousePress', 'MouseRelease', 'MouseScrollDown', 'MouseScrollUp', and 'MouseDrag'
			//
			// The basic event name is meant for entities which
			// want to respond to the mouse event happening -anywhere-.
			//
			// For events which have mouse collision enabled, they'll recieve
			// $eventName+"On" when the event occurs within their collision area.
			//
			// Mouse events all recieve an x, y, and button string.
			case mouse.Event:
				button := pmouse.GetMouseButton(int32(e.Button))
				dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
				mevent := pmouse.MouseEvent{e.X, e.Y, button}
				var eventName string
				if e.Direction == mouse.DirPress {
					setDown(button)
					eventName = "MousePress"
				} else if e.Direction == mouse.DirRelease {
					setUp(button)
					eventName = "MouseRelease"
				} else if e.Button == -2 {
					eventName = "MouseScrollDown"
				} else if e.Button == -1 {
					eventName = "MouseScrollUp"
				} else {
					eventName = "MouseDrag"
				}
				eb.Trigger(eventName, mevent)
				pmouse.Propagate(eventName+"On", mevent)

			// I don't really know what a paint event is to be honest.
			case paint.Event:

			// We hypothetically don't allow the user to manually resize
			// their window, so we don't do anything special for such events.
			case size.Event:
				fmt.Println("Window resized")

			case error:
				dlog.Error(e)
			}

			// This is a hardcoded quit function bound to the escape key.
			if IsDown("Escape") {
				if esc {
					dlog.Warn("Quiting plastic from holding ESCAPE")
					w.Send(lifecycle.Event{0, 0, nil})
				}
				esc = true
			} else {
				esc = false
			}
		}
	}()

	// This sends a signal to initiate the first scene
	initCh <- true

	// The draw loop
	// Unless told to stop, the draw channel will repeatedly
	// 1. draw black to a temporary buffer
	// 2. run any functions bound to precede drawing.
	// 3. draw all elements onto the temporary buffer.
	// 4. run any functions bound to follow drawing.
	// 5. draw the buffer's data at the viewport's position to the screen.
	// 6. publish the screen to display in window.
	go func() {
		<-drawChannel
		//cb := render.CompositeFilter(render.NewColorBox(4096, 4096, color.RGBA{0, 0, 0, 125}).Sprite)
		lastTime := time.Now()
		text := render.NewText("", float64(10+ViewX), float64(20+ViewY))
		render.Draw(text, 60000)
		for {
			dlog.Verb("Draw Loop")
		drawSelect:
			select {

			case <-drawChannel:
				dlog.Verb("Got something from draw channel")
				for {
					select {
					case <-drawChannel:
						render.Draw(text, 60000)
						break drawSelect
					case viewPoint := <-viewportChannel:
						dlog.Verb("Got something from viewport channel (waiting on draw)")
						updateScreen(viewPoint[0], viewPoint[1])
					}

				}
			case viewPoint := <-viewportChannel:
				dlog.Verb("Got something from viewport channel")
				updateScreen(viewPoint[0], viewPoint[1])
			default:
				// dlog.Verb("Default")
				eb = event.GetEventBus()
				//cb.Draw(b.RGBA())
				//draw.Draw(b.RGBA(), b.Bounds(), image.Black, image.Point{0, 0}, screen.Src)

				eb.Trigger("PreDraw", nil)
				render.DrawHeap(b, ViewX, ViewY, ScreenWidth, ScreenHeight)
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), b.RGBA(), image.Point{ViewX, ViewY}, screen.Src)
				render.DrawStaticHeap(winBuffer)
				eb.Trigger("PostDraw", b)

				w.Upload(image.Point{0, 0}, winBuffer, winBuffer.Bounds())
				w.Publish()

				timeSince := 1000000000.0 / float64(time.Since(lastTime).Nanoseconds())
				text.SetText(strconv.Itoa(int(timeSince)))
				text.SetPos(float64(10+ViewX), float64(20+ViewY))
				lastTime = time.Now()
			}
		}
	}()

	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	// It then runs any functions bound to when a frame ends.
	for {
		for runEventLoop {
			<-frameCh
			eb.Trigger("EnterFrame", nil)
			eb.Trigger("ExitFrame", nil)
			sceneCh <- true
		}
	}
}

func GetScreen() draw.Image {
	return b.RGBA()
}
