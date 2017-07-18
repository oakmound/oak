package oak

import (
	"image"
	"sync"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"

	"golang.org/x/exp/shiny/screen"
)

var (
	winBuffer     screen.Buffer
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	windowUpdateCh = make(chan bool)

	initControl = sync.Mutex{}

	lifecycleInit bool
)

func lifecycleLoop(s screen.Screen) {
	initControl.Lock()
	if lifecycleInit {
		dlog.Error("Started lifecycle twice, aborting second call")
		return
	}
	lifecycleInit = true
	initControl.Unlock()
	dlog.Info("Init Lifecycle")

	screenControl = s
	var err error

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	dlog.Info("Creating window buffer")
	winBuffer, err = screenControl.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
		return
	}

	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.\
	dlog.Info("Creating window controller")
	changeWindow(ScreenWidth, ScreenHeight)

	dlog.Info("Getting event bus")
	eb = event.GetBus()

	dlog.Info("Starting draw loop")
	go drawLoop()
	dlog.Info("Starting key hold loop")
	go keyHoldLoop()
	dlog.Info("Starting input loop")
	go inputLoop()

	dlog.Info("Starting event handler")
	go event.ResolvePending()
	// The quit channel represents a signal
	// for the engine to stop.
	<-quitCh
}

func changeWindow(width, height int) {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := windowController(screenControl, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	windowControl = wC
	windowRect = image.Rect(0, 0, width, height)
}

// ChangeWindow sets the width and height of the game window. But it doesn't.
func ChangeWindow(width, height int) {
	windowRect = image.Rect(0, 0, width, height)
}

// GetScreen returns the current screen as an rgba buffer
func GetScreen() *image.RGBA {
	return winBuffer.RGBA()
}
