package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"fmt"
	"image"
	"image/color"
	"log"
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
	initCh       = make(chan bool)
	sceneCh      = make(chan bool)
	quitCh       = make(chan bool)
	screenWidth  = 640
	screenHeight = 480
	press        = key.DirPress
	release      = key.DirRelease
	runEventLoop = false
	black        = color.RGBA{0x00, 0x00, 0x00, 0xff}
	b            screen.Buffer
	esc          = false
	l_debug      = false
)

func Init(scene string) {
	collision.Init()

	go driver.Main(eventLoop)

	prevScene := ""
	sceneMap[scene].active = true
	<-initCh
	close(initCh)
	for {
		sceneMap[scene].start(prevScene)
		cont := true
		runEventLoop = true
		for cont {
			select {
			case <-quitCh:
				return

			case <-sceneCh:
				cont = sceneMap[scene].loop()
			}
		}
		runEventLoop = false
		prevScene = scene
		scene = sceneMap[scene].end()
	}
}

func eventLoop(s screen.Screen) {
	b, _ = s.NewBuffer(image.Point{screenWidth, screenHeight})
	w, err := s.NewWindow(&screen.NewWindowOptions{screenWidth, screenHeight})
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()

	render.SetScreen((&s))

	frameRate := 60
	frameCh := make(chan bool, 100)

	go func(frameCh chan bool, frameRate int64) {
		c := time.Tick(time.Second / time.Duration(frameRate))
		for range c {
			frameCh <- true
		}
	}(frameCh, int64(frameRate))

	eb := event.GetEventBus()

	go func(eb event.EventBus, w screen.Window) {
		for {
			for runEventLoop {

				<-frameCh

				eb.Trigger("EnterFrame", nil)

				sceneCh <- true

				// To satisfy pc master race,
				// could pull this out into another
				// channel which happens as fast as possible
				// fillScreen(w, black)
				eb.Trigger("Draw", b)

				w.Upload(image.Point{0, 0}, b, b.Bounds())
				// x := w.Publish()
				// fmt.Println(x)

				eb.Trigger("ExitFrame", nil)
			}
		}
	}(eb, w)

	initCh <- true

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
			} else if e.Direction == release {
				SetUp(e.Code.String()[4:])
			}

		case mouse.Event:

		case paint.Event:

		case size.Event:

		case error:
			log.Print(e)
		}
		if IsDown("Escape") {
			if esc {
				fmt.Println("\n\n~~~~~~~~~~~~Now Escaping~~~~~~~~~~~~~~\n\n\n")
				ev := lifecycle.Event{0, 0, nil}
				w.Send(ev)
			}
			esc = true
		} else {
			esc = false
		}
	}
}

func fillScreen(w screen.Window, c color.RGBA) {
	w.Fill(b.Bounds(), black, screen.Src)
}
