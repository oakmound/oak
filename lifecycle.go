// Package oak is a game engine...
package oak

import (
	"image"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"

	"golang.org/x/exp/shiny/screen"
)

var (
	worldBuffer   screen.Buffer
	winBuffer     screen.Buffer
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	windowUpdateCH = make(chan bool)

	osCh = make(chan func())
)

//func init() {
//	runtime.LockOSThread()
//}
func lifecycleLoop(s screen.Screen) {
	screenControl = s
	var err error

	// The world buffer represents the total space that is conceptualized by the engine
	// and able to be drawn to. Space outside of this area will appear as smeared
	// white (on windows).
	worldBuffer, err = screenControl.NewBuffer(image.Point{WorldWidth, WorldHeight})
	if err != nil {
		dlog.Error(err)
		return
	}
	defer worldBuffer.Release()

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	winBuffer, err = screenControl.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
		return
	}
	defer winBuffer.Release()

	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	changeWindow(ScreenWidth, ScreenHeight)
	defer windowControl.Release()

	eb = event.GetEventBus()

	go KeyHoldLoop()
	go InputLoop()

	// Initiate the first scene
	//initCh <- true

	if conf.ShowFPS {
		go DrawLoopFPS()
	} else {
		go DrawLoopNoFPS()
	}

	go BindingLoop()
	LogicLoop()

}

// do runs f on the osLocked thread.
func osLockedFunc(f func()) {
	done := make(chan bool, 1)
	osCh <- func() {
		f()
		done <- true
	}
	<-done
}
func BindingLoop() {
	// Handle bind and unbind signals for events
	// (should be made to not use a busy loop eventually)
	for {
		for runEventLoop {
			event.ResolvePending()
		}
	}
}

func LogicLoop() {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then runs any functions bound to when a frame ends.
	// It then allows a scene to perform it's loop operation.
	c := time.Tick(time.Second / time.Duration(int64(FrameRate)))
	for range c {
		if runEventLoop {
			<-eb.TriggerBack("EnterFrame", nil)
			sceneCh <- true
		}
	}
}

func GetScreen() *image.RGBA {
	return winBuffer.RGBA()
}

func GetWorld() *image.RGBA {
	return worldBuffer.RGBA()
}

func SetWorldSize(x, y int) {
	worldBuffer, _ = screenControl.NewBuffer(image.Point{x, y})
}

func changeWindow(width, height int) {
	//windowFlag := windowControl != nil
	//if windowFlag {
	//	windowUpdateCH <- true
	//	windowControl.Publish()
	//	windowControl.Release()
	//}
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := WindowController(screenControl, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	windowControl = wC
	windowRect = image.Rect(0, 0, width, height)

	//if windowFlag {
	//	eFilter = gesture.EventFilter{EventDeque: windowControl}
	//	windowUpdateCH <- true
	//}
}

func ChangeWindow(width, height int) {
	//osLockedFunc(func() { changeWindow(width, height) })
	windowRect = image.Rect(0, 0, width, height)
}
