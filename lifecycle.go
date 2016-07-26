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
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	initCh        = make(chan bool)
	sceneCh       = make(chan bool)
	quitCh        = make(chan bool)
	drawChannel   = make(chan bool)
	drawInit      = false
	runEventLoop  = false
	ScreenWidth   int
	ScreenHeight  int
	press         = key.DirPress
	release       = key.DirRelease
	black         = color.RGBA{0x00, 0x00, 0x00, 0xff}
	b             screen.Buffer
	winBuffer     screen.Buffer
	eb            *event.EventBus
	ViewX         = 0
	ViewY         = 0
	useViewBounds = false
	viewBounds    []int
	esc           = false
	l_debug       = false
	wd, _         = os.Getwd()
	imageDir      string
	audioDir      string
)

// Scene loop initialization
func Init(firstScene string) {
	dlog.CreateLogFile()

	loadDefaultConf()
	// Set variables from conf file
	dlog.SetStringDebugLevel(conf.Debug.Level)
	dlog.SetDebugFilter(conf.Debug.Filter)

	ScreenWidth = conf.Screen.Width
	ScreenHeight = conf.Screen.Height

	imageDir = filepath.Join(filepath.Dir(wd),
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir = filepath.Join(filepath.Dir(wd),
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)
	// END of loading variables from configuration

	collision.Init()
	pmouse.Init()
	render.InitDrawHeap()
	audio.InitWinAudio()

	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)
	dlog.Info("The seed is:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~\nTHE SEED IS:", curSeed, "\n~~~~~~~~~~~~~~~\n")

	go driver.Main(eventLoop)

	prevScene := ""
	sceneMap[firstScene].active = true
	<-initCh
	close(initCh)
	runEventLoop = true
	scene := firstScene
	dlog.Info("First Scene Start")
	for {
		dlog.Info("~~~~~~~~~~~Scene Start~~~~~~~~~")
		sceneMap[scene].start(prevScene)
		// Send a signal to resume (or begin) drawing
		drawChannel <- true

		cont := true
		for cont {
			select {
			case <-quitCh:
				return

			case <-sceneCh:
				cont = sceneMap[scene].loop()
			}
		}
		dlog.Info("~~~~~~~~Scene End~~~~~~~~~~")
		prevScene = scene

		// Send a signal to stop drawing
		drawChannel <- true

		event.ResetEventBus()
		render.ResetDrawHeap()
		collision.Clear()
		render.PreDraw(0, nil)

		scene = sceneMap[scene].end()

		eb = event.GetEventBus()
	}
}

func eventLoop(s screen.Screen) {
	// Todo: ratify these world size points
	b, _ = s.NewBuffer(image.Point{4000, 4000})
	winBuffer, _ = s.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	w, err := s.NewWindow(&screen.NewWindowOptions{ScreenWidth, ScreenHeight})
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()
	render.InitFont(&b)

	err = render.BatchLoad(imageDir)
	if err != nil {
		dlog.Error(err)
		return
	}
	err = audio.BatchLoad(audioDir)
	if err != nil {
		dlog.Error(err)
		return
	}

	eb = event.GetEventBus()

	frameRate := 60
	frameCh := make(chan bool, 100)

	go func(frameCh chan bool, frameRate int64) {
		c := time.Tick(time.Second / time.Duration(frameRate))
		for range c {
			frameCh <- true
		}
	}(frameCh, int64(frameRate))

	go func() {
		for {
			// Handle window events
			e := w.NextEvent()
			// This print message is to help programmers learn what events this
			// example program generates. A real program shouldn't print such
			// messages; they're not important to end users.
			format := "got %#v\n"
			if _, ok := e.(fmt.Stringer); ok {
				format = "got %v\n"
			}
			if l_debug {
				fmt.Printf(format, e)
			}
			switch e := e.(type) {

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					quitCh <- true
					return
				}

			case key.Event:
				if e.Direction == press {
					fmt.Println("--------------------", e.Code.String()[4:])
					SetDown(e.Code.String()[4:])
					eb.Trigger("KeyDown", e.Code.String()[4:])
					eb.Trigger("KeyDown"+e.Code.String()[4:], nil)
				} else if e.Direction == release {
					SetUp(e.Code.String()[4:])
					eb.Trigger("KeyUp", e.Code.String()[4:])
					eb.Trigger("KeyUp"+e.Code.String()[4:], nil)
				}

			case mouse.Event:
				button := pmouse.GetMouseButton(int32(e.Button))
				dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
				mevent := pmouse.MouseEvent{e.X, e.Y, button}
				var eventName string
				if e.Direction == mouse.DirPress {
					SetDown(button)
					eventName = "MousePress"
				} else if e.Direction == mouse.DirRelease {
					SetUp(button)
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

			case paint.Event:

			case size.Event:
				fmt.Println("Window resized")

			case error:
				log.Print(e)
			}
			if IsDown("Escape") {
				if esc {
					dlog.Warn("\n\n~~~~~~~~~~~~Now Escaping~~~~~~~~~~~~~~\n\n\n")
					ev := lifecycle.Event{0, 0, nil}
					w.Send(ev)
				}
				esc = true
			} else {
				esc = false
			}
		}
	}()

	initCh <- true

	// Draw loop
	// Pulled away from the framerate loop below
	go func() {
		<-drawChannel
		for {
			select {
			case <-drawChannel:
				<-drawChannel
			default:
				eb = event.GetEventBus()
				draw.Draw(b.RGBA(), b.Bounds(), image.Black, image.Point{0, 0}, screen.Src)

				eb.Trigger("PreDraw", nil)
				render.DrawHeap(b)
				eb.Trigger("PostDraw", b)
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), b.RGBA(), image.Point{ViewX, ViewY}, screen.Src)

				w.Upload(image.Point{0, 0}, winBuffer, winBuffer.Bounds())
				w.Publish()
			}
		}
	}()

	for {
		for runEventLoop {
			<-frameCh

			eb.Trigger("EnterFrame", nil)
			sceneCh <- true

			eb.Trigger("ExitFrame", nil)
		}
	}
}

func fillScreen(w screen.Window, c color.RGBA) {
	w.Fill(b.Bounds(), black, screen.Src)
}

func SetScreen(x, y int) {
	if useViewBounds {
		if viewBounds[0] < x && viewBounds[2] > x+ScreenWidth {
			dlog.Verb("Set ViewX to ", x)
			ViewX = x
		} else if viewBounds[0] > x {
			ViewX = viewBounds[0]
		} else if viewBounds[2] < x+ScreenWidth {
			ViewX = viewBounds[2] - ScreenWidth
		}
		if viewBounds[1] < y && viewBounds[3] > y+ScreenHeight {
			dlog.Verb("Set ViewY to ", y)
			ViewY = y
		} else if viewBounds[1] > y {
			ViewY = viewBounds[1]
		} else if viewBounds[3] < y+ScreenHeight {
			ViewY = viewBounds[3] - ScreenHeight
		}

	} else {
		dlog.Verb("Set ViewXY to ", x, " ", y)
		ViewX = x
		ViewY = y
	}
	dlog.Verb("ViewX, Y: ", ViewX, " ", ViewY)
	eb.Trigger("ViewportUpdate", []float64{float64(ViewX), float64(ViewY)})
}
func SetViewportBounds(x1, y1, x2, y2 int) {
	dlog.Info("Viewport bounds set to, ", x1, y1, x2, y2)
	useViewBounds = true
	viewBounds = []int{x1, y1, x2, y2}
}
