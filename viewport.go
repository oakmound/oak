package oak

import (
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
)

// SetScreen positions the viewport to be at x,y
func (w *Window) SetScreen(x, y int) {
	w.setViewport(intgeom.Point2{x, y})
}

// ShiftScreen shifts the viewport by x,y
func (w *Window) ShiftScreen(x, y int) {
	w.setViewport(w.viewPos.Add(intgeom.Point2{x, y}))
}

func (w *Window) setViewport(pt intgeom.Point2) {
	if w.useViewBounds {
		if w.viewBounds.Min.X() <= pt.X() && w.viewBounds.Max.X() >= pt.X()+w.ScreenWidth {
			w.viewPos[0] = pt.X()
		} else if w.viewBounds.Min.X() > pt.X() {
			w.viewPos[0] = w.viewBounds.Min.X()
		} else if w.viewBounds.Max.X() < pt.X()+w.ScreenWidth {
			w.viewPos[0] = w.viewBounds.Max.X() - w.ScreenWidth
		}
		if w.viewBounds.Min.Y() <= pt.Y() && w.viewBounds.Max.Y() >= pt.Y()+w.ScreenHeight {
			w.viewPos[1] = pt.Y()
		} else if w.viewBounds.Min.Y() > pt.Y() {
			w.viewPos[1] = w.viewBounds.Min.Y()
		} else if w.viewBounds.Max.Y() < pt.Y()+w.ScreenHeight {
			w.viewPos[1] = w.viewBounds.Max.Y() - w.ScreenHeight
		}
	} else {
		w.viewPos = pt
	}
	event.TriggerOn(w.eventHandler, ViewportUpdate, w.viewPos)
}

// GetViewportBounds reports what bounds the viewport has been set to, if any.
func (w *Window) GetViewportBounds() (rect intgeom.Rect2, ok bool) {
	return w.viewBounds, w.useViewBounds
}

// RemoveViewportBounds removes restrictions on the viewport's movement. It will not
// cause the viewport to update immediately.
func (w *Window) RemoveViewportBounds() {
	w.useViewBounds = false
}

// SetViewportBounds sets the minimum and maximum position of the viewport, including
// screen dimensions
func (w *Window) SetViewportBounds(rect intgeom.Rect2) {
	if rect.Max[0] < w.ScreenWidth {
		rect.Max[0] = w.ScreenWidth
	}
	if rect.Max[1] < w.ScreenHeight {
		rect.Max[1] = w.ScreenHeight
	}
	w.useViewBounds = true
	w.viewBounds = rect

	newViewX := w.viewPos.X()
	newViewY := w.viewPos.Y()
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

	if newViewX != w.viewPos.X() || newViewY != w.viewPos.Y() {
		w.setViewport(intgeom.Point2{newViewX, newViewY})
	}
}
