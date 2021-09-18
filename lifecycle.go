package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v3/alg"
	"github.com/oakmound/oak/v3/debugstream"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/oakmound/oak/v3/shiny/screen"
)

func (w *Window) lifecycleLoop(s screen.Screen) {
	w.screenControl = s
	err := w.UpdateViewSize(w.ScreenWidth, w.ScreenHeight)
	if err != nil {
		go w.exitWithError(err)
		return
	}

	// Right here, query the backing scale factor of the physical screen
	// Apply that factor to the scale

	err = w.newWindow(
		int32(w.config.Screen.X),
		int32(w.config.Screen.Y),
		int(float64(w.ScreenWidth)*w.config.Screen.Scale),
		int(float64(w.ScreenHeight)*w.config.Screen.Scale),
	)
	if err != nil {
		go w.exitWithError(err)
		return
	}

	go w.drawLoop()
	go w.inputLoop()

	<-w.quitCh
	w.windowControl.Release()
}

// Quit sends a signal to the window to close itself, closing the window and
// any spun up resources. It should not be called before Init. After it is called,
// it must not be called again.
func (w *Window) Quit() {
	// We could have hit this before the window was created
	if w.windowControl == nil {
		close(w.quitCh)
	} else {
		w.windowControl.Send(lifecycle.Event{To: lifecycle.StageDead})
	}
	if w.config.EnableDebugConsole {
		debugstream.DefaultCommands.RemoveScope(w.ControllerID)
	}
}

func (w *Window) newWindow(x, y int32, width, height int) error {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := w.windowController(w.screenControl, x, y, width, height)
	if err != nil {
		return err
	}
	w.windowControl = wC
	return w.ChangeWindow(width, height)
}

// SetAspectRatio will enforce that the displayed window does not distort the
// input screen away from the given x:y ratio. The screen will not use these
// settings until a new size event is received from the OS.
func (w *Window) SetAspectRatio(xToY float64) {
	w.UseAspectRatio = true
	w.aspectRatio = xToY
}

// ChangeWindow sets the width and height of the game window. Although exported,
// calling it without a size event will probably not act as expected.
func (w *Window) ChangeWindow(width, height int) error {
	// Draw a black frame to cover up smears
	// Todo: could restrict the black to -just- the area not covered by the
	// scaled screen buffer
	buff, err := w.screenControl.NewImage(image.Point{width, height})
	if err == nil {
		draw.Draw(buff.RGBA(), buff.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
		w.windowControl.Upload(zeroPoint, buff, buff.Bounds())
	} else {
		return err
	}
	var x, y int
	if w.UseAspectRatio {
		inRatio := float64(width) / float64(height)
		if w.aspectRatio > inRatio {
			newHeight := alg.RoundF64(float64(height) * (inRatio / w.aspectRatio))
			y = (newHeight - height) / 2
			height = newHeight - y
		} else {
			newWidth := alg.RoundF64(float64(width) * (w.aspectRatio / inRatio))
			x = (newWidth - width) / 2
			width = newWidth - x
		}
	}
	w.windowRect = image.Rect(-x, -y, width, height)
	return nil
}

// UpdateViewSize updates the size of this window's viewport.
func (w *Window) UpdateViewSize(width, height int) error {
	w.ScreenWidth = width
	w.ScreenHeight = height
	newBuffer, err := w.screenControl.NewImage(image.Point{width, height})
	if err != nil {
		return err
	}
	w.winBuffer = newBuffer
	newTexture, err := w.screenControl.NewTexture(w.winBuffer.Bounds().Max)
	if err != nil {
		return err
	}
	w.windowTexture = newTexture
	return nil
}
