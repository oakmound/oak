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

	esc      bool
	drawInit bool
)

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
	windowControl, err := WindowController(screenControl, ScreenWidth, ScreenHeight)
	if err != nil {
		dlog.Error(err)
		return
	}
	defer windowControl.Release()

	eb = event.GetEventBus()

	frameCh := make(chan bool)

	go FrameLoop(frameCh, int64(FrameRate))
	go InputLoop(windowControl)

	// Initiate the first scene
	initCh <- true

	go DrawLoop(windowControl)
	go BindingLoop()
	LogicLoop(frameCh)
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
			<-eb.Trigger("EnterFrame", nil)
			// ExitFrame shouldn't need to exist given event priorities
			<-eb.Trigger("ExitFrame", nil)
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
