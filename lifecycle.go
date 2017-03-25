// Package oak is a game engine...
package oak

import (
	"image"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"

	"fmt"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

var (
	worldBuffer   screen.Buffer
	winBuffer     screen.Buffer
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	drawInit       bool
	windowUpdateCH = make(chan bool)

	osCh = make(chan func())
)

//func init() {
//	runtime.LockOSThread()
//}
func lifecycleLoop(s screen.Screen) {
	screenControl = s

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
	fmt.Println("changed")
	defer windowControl.Release()

	eb = event.GetEventBus()

	frameCh := make(chan bool)

	go FrameLoop(frameCh, int64(FrameRate))
	go KeyHoldLoop()
	go InputLoop()

	// Initiate the first scene
	initCh <- true

	if conf.ShowFPS {
		go DrawLoopFPS()
	} else {
		go DrawLoopNoFPS()
	}

	go BindingLoop()
	LogicLoop(frameCh)

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
	for runEventLoop {
		event.ResolvePending()
	}
}

// Maintain a frame rate for logical operations
func FrameLoop(frameCh chan bool, frameRate int64) {
	c := time.Tick(time.Second / time.Duration(frameRate))
	for range c {
		frameCh <- true
	}
}

func LogicLoop(frameCh chan bool) {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then runs any functions bound to when a frame ends.
	// It then allows a scene to perform it's loop operation.
	for {
		for runEventLoop {
			<-frameCh
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
	windowControl.Send(size.Event{width, height, geom.Pt(float32(width)), geom.Pt(float32(height)), 1, 0})
	windowRect = image.Rect(0, 0, width, height)
}
