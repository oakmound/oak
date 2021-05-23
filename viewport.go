package oak

import (
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
)

// SetScreen positions the viewport to be at x,y
func (c *Controller) SetScreen(x, y int) {
	c.setViewport(intgeom.Point2{x, y})
}

// ShiftScreen shifts the viewport by x,y
func (c *Controller) ShiftScreen(x, y int) {
	c.setViewport(c.viewPos.Add(intgeom.Point2{x, y}))
}

func (c *Controller) setViewport(pt intgeom.Point2) {
	if c.useViewBounds {
		if c.viewBounds.Min.X() <= pt.X() && c.viewBounds.Max.X() >= pt.X()+c.ScreenWidth {
			c.viewPos[0] = pt.X()
		} else if c.viewBounds.Min.X() > pt.X() {
			c.viewPos[0] = c.viewBounds.Min.X()
		} else if c.viewBounds.Max.X() < pt.X()+c.ScreenWidth {
			c.viewPos[0] = c.viewBounds.Max.X() - c.ScreenWidth
		}
		if c.viewBounds.Min.Y() <= pt.Y() && c.viewBounds.Max.Y() >= pt.Y()+c.ScreenHeight {
			c.viewPos[1] = pt.Y()
		} else if c.viewBounds.Min.Y() > pt.Y() {
			c.viewPos[1] = c.viewBounds.Min.Y()
		} else if c.viewBounds.Max.Y() < pt.Y()+c.ScreenHeight {
			c.viewPos[1] = c.viewBounds.Max.Y() - c.ScreenHeight
		}
	} else {
		c.viewPos = pt
	}
	c.logicHandler.Trigger(event.ViewportUpdate, c.viewPos)
}

// GetViewportBounds reports what bounds the viewport has been set to, if any.
func (c *Controller) GetViewportBounds() (rect intgeom.Rect2, ok bool) {
	return c.viewBounds, c.useViewBounds
}

// RemoveViewportBounds removes restrictions on the viewport's movement. It will not
// cause the viewport to update immediately.
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

	newViewX := c.viewPos.X()
	newViewY := c.viewPos.Y()
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

	if newViewX != c.viewPos.X() || newViewY != c.viewPos.Y() {
		c.setViewport(intgeom.Point2{newViewX, newViewY})
	}
}
