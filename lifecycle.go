package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v3/alg"
	"github.com/oakmound/oak/v3/debugstream"
	"github.com/oakmound/oak/v3/dlog"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/oakmound/oak/v3/shiny/screen"
)

func (c *Controller) lifecycleLoop(s screen.Screen) {
	dlog.Info("Init Lifecycle")

	c.screenControl = s
	dlog.Info("Creating window buffer")
	err := c.UpdateViewSize(c.ScreenWidth, c.ScreenHeight)
	if err != nil {
		go c.exitWithError(err)
		return
	}

	// Right here, query the backing scale factor of the physical screen
	// Apply that factor to the scale

	dlog.Info("Creating window controller")
	err = c.newWindow(
		int32(c.config.Screen.X),
		int32(c.config.Screen.Y),
		c.ScreenWidth*c.config.Screen.Scale,
		c.ScreenHeight*c.config.Screen.Scale,
	)
	if err != nil {
		go c.exitWithError(err)
		return
	}

	dlog.Info("Starting draw loop")
	go c.drawLoop()
	dlog.Info("Starting input loop")
	go c.inputLoop()

	<-c.quitCh
}

// Quit sends a signal to the window to close itself, closing the window and
// any spun up resources.
func (c *Controller) Quit() {
	c.windowControl.Send(lifecycle.Event{To: lifecycle.StageDead})
	if c.config.EnableDebugConsole {
		debugstream.DefaultCommands.RemoveScope(c.ControllerID)
	}
}

func (c *Controller) newWindow(x, y int32, width, height int) error {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := c.windowController(c.screenControl, x, y, width, height)
	if err != nil {
		return err
	}
	c.windowControl = wC
	c.ChangeWindow(width, height)
	return nil
}

// SetAspectRatio will enforce that the displayed window does not distort the
// input screen away from the given x:y ratio. The screen will not use these
// settings until a new size event is received from the OS.
func (c *Controller) SetAspectRatio(xToY float64) {
	c.UseAspectRatio = true
	c.aspectRatio = xToY
}

// ChangeWindow sets the width and height of the game window. Although exported,
// calling it without a size event will probably not act as expected.
func (c *Controller) ChangeWindow(width, height int) {
	// Draw a black frame to cover up smears
	// Todo: could restrict the black to -just- the area not covered by the
	// scaled screen buffer
	buff, err := c.screenControl.NewImage(image.Point{width, height})
	if err == nil {
		draw.Draw(buff.RGBA(), buff.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
		c.windowControl.Upload(zeroPoint, buff, buff.Bounds())
	} else {
		dlog.Error(err)
	}
	var x, y int
	if c.UseAspectRatio {
		inRatio := float64(width) / float64(height)
		if c.aspectRatio > inRatio {
			newHeight := alg.RoundF64(float64(height) * (inRatio / c.aspectRatio))
			y = (newHeight - height) / 2
			height = newHeight - y
		} else {
			newWidth := alg.RoundF64(float64(width) * (c.aspectRatio / inRatio))
			x = (newWidth - width) / 2
			width = newWidth - x
		}
	}
	c.windowRect = image.Rect(-x, -y, width, height)
}

func (c *Controller) UpdateViewSize(width, height int) error {
	c.ScreenWidth = width
	c.ScreenHeight = height
	newBuffer, err := c.screenControl.NewImage(image.Point{width, height})
	if err != nil {
		return err
	}
	c.winBuffer = newBuffer
	newTexture, err := c.screenControl.NewTexture(c.winBuffer.Bounds().Max)
	if err != nil {
		return err
	}
	c.windowTexture = newTexture
	return nil
}
