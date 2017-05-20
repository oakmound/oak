// Package oak is a game engine...
package oak

import (
	"image"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/timing"

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

	LogicTicker *timing.DynamicTicker
)

//func init() {
//	runtime.LockOSThread()
//}
func lifecycleLoop(s screen.Screen) {

	screenControl = s
	var err error

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	winBuffer, err = screenControl.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
		return
	}
	//defer winBuffer.Release()

	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	changeWindow(ScreenWidth, ScreenHeight)
	//defer windowControl.Release()

	eb = event.GetEventBus()

	go KeyHoldLoop()
	go InputLoop()
	go DrawLoop()

	event.ResolvePending()
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

func LogicLoop() chan bool {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	ch := make(chan bool)
	go func(doneCh chan bool) {
		LogicTicker = timing.NewDynamicTicker()
		LogicTicker.SetTick(timing.FPSToDuration(FrameRate))
		for {
			select {
			case <-LogicTicker.C:
				<-eb.TriggerBack("EnterFrame", nil)
				sceneCh <- true
			case <-doneCh:
				LogicTicker.Stop()
				return
			}
		}
	}(ch)
	return ch
}

func GetScreen() *image.RGBA {
	return winBuffer.RGBA()
}

func GetWorld() *image.RGBA {
	return worldBuffer.RGBA()
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
