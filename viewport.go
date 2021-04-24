package oak

import (
	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
)

// SetScreen sends a signal to the draw loop to set the viewport to be at x,y
func (c *Controller) SetScreen(x, y int) {
	c.viewportCh <- intgeom.Point2{x, y}
}

// ShiftScreen sends a signal to the draw loop to shift the viewport by x,y
func (c *Controller) ShiftScreen(x, y int) {
	c.viewportShiftCh <- intgeom.Point2{x, y}
}

func (c *Controller) updateScreen(pt intgeom.Point2) {
	c.viewPosMutex.Lock()
	c.setViewport(pt)
	c.viewPosMutex.Unlock()
}

func (c *Controller) shiftViewPort(pt intgeom.Point2) {
	c.viewPosMutex.Lock()
	pt = pt.Add(c.ViewPos)
	c.setViewport(pt)
	c.viewPosMutex.Unlock()
}
func (c *Controller) setViewport(pt intgeom.Point2) {
	if c.useViewBounds {
		if c.viewBounds.Min.X() <= pt.X() && c.viewBounds.Max.X() >= pt.X()+c.ScreenWidth {
			c.ViewPos[0] = pt.X()
		} else if c.viewBounds.Min.X() > pt.X() {
			c.ViewPos[0] = c.viewBounds.Min.X()
		} else if c.viewBounds.Max.X() < pt.X()+c.ScreenWidth {
			c.ViewPos[0] = c.viewBounds.Max.X() - c.ScreenWidth
		}
		if c.viewBounds.Min.Y() <= pt.Y() && c.viewBounds.Max.Y() >= pt.Y()+c.ScreenHeight {
			c.ViewPos[1] = pt.Y()
		} else if c.viewBounds.Min.Y() > pt.Y() {
			c.ViewPos[1] = c.viewBounds.Min.Y()
		} else if c.viewBounds.Max.Y() < pt.Y()+c.ScreenHeight {
			c.ViewPos[1] = c.viewBounds.Max.Y() - c.ScreenHeight
		}
	} else {
		c.ViewPos = pt
	}
	c.logicHandler.Trigger(event.ViewportUpdate, c.ViewPos)
}

// GetViewportBounds reports what bounds the viewport has been set to, if any.
func (c *Controller) GetViewportBounds() (rect intgeom.Rect2, ok bool) {
	return c.viewBounds, c.useViewBounds
}

// RemoveViewportBounds removes restrictions on the viewport's movement. It will not
// cause ViewPos to update immediately.
func (c *Controller) RemoveViewportBounds() {
	c.useViewBounds = false
}

// SetViewportBounds sets the minimum and maximum position of the viewport, including
// screen dimensions
func (c *Controller) SetViewportBounds(rect intgeom.Rect2) {
	if rect.Max[0] < c.ScreenWidth {
		rect.Max[0] = c.ScreenWidth
	}
	if rect.Max[1] < c.ScreenHeight {
		rect.Max[1] = c.ScreenHeight
	}
	c.useViewBounds = true
	c.viewBounds = rect

	newViewX := c.ViewPos.X()
	newViewY := c.ViewPos.Y()
	if newViewX < rect.Min[0] {
		newViewX = rect.Min[0]
	} else if newViewX > rect.Max[0] {
		newViewX = rect.Max[0]
	}
	if newViewY < rect.Min[1] {
		newViewY = rect.Min[1]
	} else if newViewY > rect.Max[1] {
		newViewY = rect.Max[1]
	}

	if newViewX != c.ViewPos.X() || newViewY != c.ViewPos.Y() {
		c.viewportCh <- intgeom.Point2{newViewX, newViewY}
	}
}

func (c *Controller) moveViewportBinding(speed int) func(event.CID, interface{}) int {
	return func(cID event.CID, n interface{}) int {
		dX := 0
		dY := 0
		if c.IsDown(key.UpArrow) {
			dY--
		}
		if c.IsDown(key.DownArrow) {
			dY++
		}
		if c.IsDown(key.LeftArrow) {
			dX--
		}
		if c.IsDown(key.RightArrow) {
			dX++
		}
		c.ViewPos[0] += dX * speed
		c.ViewPos[1] += dY * speed
		if c.viewportLocked {
			return event.UnbindSingle
		}
		return 0
	}
}
