package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v2/alg"
	"github.com/oakmound/oak/v2/dlog"
	"golang.org/x/mobile/event/lifecycle"

	"github.com/oakmound/oak/v2/shiny/screen"
)

func (c *Controller) lifecycleLoop(s screen.Screen) {
	dlog.Info("Init Lifecycle")

	c.screenControl = s
	dlog.Info("Creating window buffer")
	err := c.UpdateViewSize(c.ScreenWidth, c.ScreenHeight)
	if err != nil {
		dlog.Error(err)
		return
	}

	// Right here, query the backing scale factor of the physical screen
	// Apply that factor to the scale

	dlog.Info("Creating window controller")
	c.newWindow(
		int32(c.config.Screen.X),
		int32(c.config.Screen.Y),
		c.ScreenWidth*c.config.Screen.Scale,
		c.ScreenHeight*c.config.Screen.Scale,
	)

	dlog.Info("Starting draw loop")
	go c.drawLoop()
	dlog.Info("Starting input loop")
	go c.inputLoop()

	<-c.quitCh
}

// Quit sends a signal to the window to close itself, ending oak.
func (c *Controller) Quit() {
	c.windowControl.Send(lifecycle.Event{To: lifecycle.StageDead})
}

func (c *Controller) newWindow(x, y int32, width, height int) {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := c.windowController(c.screenControl, x, y, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	c.windowControl = wC
	c.ChangeWindow(width, height)
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

// GetScreen returns the current screen as an rgba buffer
func (c *Controller) GetScreen() *image.RGBA {
	return c.winBuffer.RGBA()
}
