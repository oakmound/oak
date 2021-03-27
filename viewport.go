package oak

import (
	"image"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
)

type rect struct {
	minX, minY, maxX, maxY int
}

// SetScreen sends a signal to the draw loop to set the viewport to be at x,y
func (c *Controller) SetScreen(x, y int) {
	dlog.Verb("Requesting ViewPoint ", x, y)
	c.viewportCh <- [2]int{x, y}
}

// ShiftScreen sends a signal to the draw loop to shift the viewport by x,y
func (c *Controller) ShiftScreen(x, y int) {
	dlog.Verb("Requesting shift of ViewPoint by ", x, y)
	c.viewportShiftCh <- [2]int{x, y}
}

func (c *Controller) updateScreen(x, y int) {
	c.viewPosMutex.Lock()
	c.setViewport(x, y)
	c.viewPosMutex.Unlock()
}

func (c *Controller) shiftViewPort(x, y int) {
	c.viewPosMutex.Lock()
	c.setViewport(c.ViewPos.X+x, c.ViewPos.Y+y)
	c.viewPosMutex.Unlock()
}
func (c *Controller) setViewport(x, y int) {
	if c.useViewBounds {
		if c.viewBounds.minX <= x && c.viewBounds.maxX >= x+c.ScreenWidth {
			dlog.Verb("Set ViewX to ", x)
			c.ViewPos.X = x
		} else if c.viewBounds.minX > x {
			c.ViewPos.X = c.viewBounds.minX
		} else if c.viewBounds.maxX < x+c.ScreenWidth {
			c.ViewPos.X = c.viewBounds.maxX - c.ScreenWidth
		}
		if c.viewBounds.minY <= y && c.viewBounds.maxY >= y+c.ScreenHeight {
			dlog.Verb("Set ViewY to ", y)
			c.ViewPos.Y = y
		} else if c.viewBounds.minY > y {
			c.ViewPos.Y = c.viewBounds.minY
		} else if c.viewBounds.maxY < y+c.ScreenHeight {
			c.ViewPos.Y = c.viewBounds.maxY - c.ScreenHeight
		}
	} else {
		dlog.Verb("Set ViewXY to ", x, " ", y)
		c.ViewPos = image.Point{x, y}
	}
	c.logicHandler.Trigger(event.ViewportUpdate, []float64{float64(c.ViewPos.X), float64(c.ViewPos.Y)})
	dlog.Verb("ViewX, Y: ", c.ViewPos.X, " ", c.ViewPos.Y)
}

// GetViewportBounds reports what bounds the viewport has been set to, if any.
func (c *Controller) GetViewportBounds() (x1, y1, x2, y2 int, ok bool) {
	if c.useViewBounds {
		return c.viewBounds.minX, c.viewBounds.minY, c.viewBounds.maxX, c.viewBounds.maxY, true
	}
	return 0, 0, 0, 0, false
}

// RemoveViewportBounds removes restrictions on the viewport's movement. It will not
// cause ViewPos to update immediately.
func (c *Controller) RemoveViewportBounds() {
	c.useViewBounds = false
}

// SetViewportBounds sets the minimum and maximum position of the viewport, including
// screen dimensions
func (c *Controller) SetViewportBounds(x1, y1, x2, y2 int) {
	if x2 < c.ScreenWidth {
		x2 = c.ScreenWidth
	}
	if y2 < c.ScreenHeight {
		y2 = c.ScreenHeight
	}
	c.useViewBounds = true
	c.viewBounds = rect{x1, y1, x2, y2}

	dlog.Info("Viewport bounds set to, ", x1, y1, x2, y2)

	newViewX := c.ViewPos.X
	newViewY := c.ViewPos.Y
	if newViewX < x1 {
		newViewX = x1
	} else if newViewX > x2 {
		newViewX = x2
	}
	if newViewY < y1 {
		newViewY = y1
	} else if newViewY > y2 {
		newViewY = y2
	}

	if newViewX != c.ViewPos.X || newViewY != c.ViewPos.Y {
		c.viewportCh <- [2]int{newViewX, newViewY}
	}
}

func (c *Controller) moveViewportBinding(speed int) func(event.CID, interface{}) int {
	return func(cID event.CID, n interface{}) int {
		dX := 0
		dY := 0
		if c.IsDown("UpArrow") {
			dY--
		}
		if c.IsDown("DownArrow") {
			dY++
		}
		if c.IsDown("LeftArrow") {
			dX--
		}
		if c.IsDown("RightArrow") {
			dX++
		}
		c.ViewPos.X += dX * speed
		c.ViewPos.Y += dY * speed
		if c.viewportLocked {
			return event.UnbindSingle
		}
		return 0
	}
}
