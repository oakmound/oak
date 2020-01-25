package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v2/alg"
	"github.com/oakmound/oak/v2/dlog"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/oakmound/shiny/screen"
)

var (
	winBuffer     screen.Image
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	windowUpdateCh = make(chan bool)
)

func lifecycleLoop(s screen.Screen) {
	dlog.Info("Init Lifecycle")

	screenControl = s
	var err error

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	dlog.Info("Creating window buffer")
	winBuffer, err = screenControl.NewImage(image.Point{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
		return
	}

	// Next time:
	// Right here, query the backing scale factor of the physical screen
	// Apply that factor to the scale

	dlog.Info("Creating window controller")
	changeWindow(int32(conf.Screen.X), int32(conf.Screen.Y), ScreenWidth*conf.Screen.Scale, ScreenHeight*conf.Screen.Scale)

	dlog.Info("Starting draw loop")
	go drawLoop()
	dlog.Info("Starting input loop")
	go inputLoop()

	// The quit channel represents a signal
	// for the engine to stop.
	<-quitCh
}

// Quit sends a signal to the window to close itself, ending oak.
func Quit() {
	windowControl.Send(lifecycle.Event{To: lifecycle.StageDead})
}

func changeWindow(x, y int32, width, height int) {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := windowController(screenControl, x, y, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	windowControl = wC
	ChangeWindow(width, height)
}

var (
	// UseAspectRatio determines whether new window changes will distort or
	// maintain the relative width to height ratio of the screen buffer.
	UseAspectRatio = false
	aspectRatio    float64
)

// SetAspectRatio will enforce that the displayed window does not distort the
// input screen away from the given x:y ratio. The screen will not use these
// settings until a new size event is received from the OS.
func SetAspectRatio(xToY float64) {
	UseAspectRatio = true
	aspectRatio = xToY
}

// ChangeWindow sets the width and height of the game window. Although exported,
// calling it without a size event will probably not act as expected.
func ChangeWindow(width, height int) {
	// Draw a black frame to cover up smears
	// Todo: could restrict the black to -just- the area not covered by the
	// scaled screen buffer
	buff, err := screenControl.NewImage(image.Point{width, height})
	if err == nil {
		draw.Draw(buff.RGBA(), buff.Bounds(), Background, zeroPoint, draw.Src)
		windowControl.Upload(zeroPoint, buff, buff.Bounds())
	} else {
		dlog.Error(err)
	}
	var x, y int
	if UseAspectRatio {
		inRatio := float64(width) / float64(height)
		if aspectRatio > inRatio {
			newHeight := alg.RoundF64(float64(height) * (inRatio / aspectRatio))
			y = (newHeight - height) / 2
			height = newHeight - y
		} else {
			newWidth := alg.RoundF64(float64(width) * (aspectRatio / inRatio))
			x = (newWidth - width) / 2
			width = newWidth - x
		}
	}
	windowRect = image.Rect(-x, -y, width, height)
}

// GetScreen returns the current screen as an rgba buffer
func GetScreen() *image.RGBA {
	return winBuffer.RGBA()
}
