package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"

	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
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
	ScreenWidth   = 640
	ScreenHeight  = 480
	press         = key.DirPress
	release       = key.DirRelease
	runEventLoop  = false
	drawChannel   = make(chan bool)
	black         = color.RGBA{0x00, 0x00, 0x00, 0xff}
	b             screen.Buffer
	winBuffer     screen.Buffer
	eb            event.EventBus
	viewX         = 0
	viewY         = 0
	useViewBounds = false
	viewBounds    []int
	esc           = false
	l_debug       = false
)

// Scene loop initialization
func Init(firstScene string) {
	dlog.CreateLogFile()
	collision.Init()
	render.InitDrawHeap()
	winaudio.InitWinAudio()

	curSeed := time.Now().UTC().UnixNano()
	// Basic seed
	//curSeed = 1463358974925095300
	// Seed that required modifying connection algorithm 7/2
	//curSeed = 1467565587127684400
	rand.Seed(curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~\nTHE SEED IS:", curSeed, "\n~~~~~~~~~~~~~~~\n")

	go driver.Main(eventLoop)

	prevScene := ""
	sceneMap[firstScene].active = true
	<-initCh
	close(initCh)
	runEventLoop = true
	scene := firstScene
	for {
		dlog.Info("~~~~~~~~~~~Scene Start~~~~~~~~~")
		sceneMap[scene].start(prevScene)
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
		prevScene = scene
		scene = sceneMap[scene].end()
	}
}

func eventLoop(s screen.Screen) {
	b, _ = s.NewBuffer(image.Point{4000, 4000})
	winBuffer, _ = s.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	w, err := s.NewWindow(&screen.NewWindowOptions{ScreenWidth, ScreenHeight})
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()
	render.InitFont(&b)
	render.SetScreen((&s))

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
				} else if e.Direction == release {
					SetUp(e.Code.String()[4:])
					eb.Trigger("KeyUp", e.Code.String()[4:])
				}

			case mouse.Event:
				button := getMouseButton(int32(e.Button))
				dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
				mevent := MouseEvent{e.X, e.Y, button}
				if e.Direction == mouse.DirPress {
					SetDown(button)
					eb.Trigger("MousePress", mevent)
				} else if e.Direction == mouse.DirRelease {
					SetUp(button)
					eb.Trigger("MouseRelease", mevent)
				} else if e.Button == -2 {
					eb.Trigger("MouseScrollDown", mevent)
				} else if e.Button == -1 {
					eb.Trigger("MouseScrollUp", mevent)
				} else {
					eb.Trigger("MouseDrag", mevent)
				}

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
			// Comment out this for smearing, but visible text
			draw.Draw(b.RGBA(), b.Bounds(), image.Black, image.Point{0, 0}, screen.Src)

			eb.Trigger("PreDraw", nil)
			render.DrawHeap(b)
			eb.Trigger("PostDraw", b)
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), b.RGBA(), image.Point{-viewX, -viewY}, screen.Src)

			w.Upload(image.Point{0, 0}, winBuffer, winBuffer.Bounds())
			w.Publish()
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
		if viewBounds[0] > x && viewBounds[2] < x-ScreenWidth {
			viewX = x
		}
		if viewBounds[1] > y && viewBounds[3] < y-ScreenHeight {
			viewY = y
		}

	} else {
		viewX = x
		viewY = y
	}
}
func SetViewportBounds(x1, y1, x2, y2 int) {
	useViewBounds = true
	viewBounds = []int{x1, y1, x2, y2}
}
