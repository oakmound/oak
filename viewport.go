package oak

import (
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
)

type Viewport struct {
	Position       intgeom.Point2
	Bounds         intgeom.Rect2
	BoundsEnforced bool
}

// ShiftViewport shifts the viewport by x,y
func (w *Window) ShiftViewport(delta intgeom.Point2) {
	w.SetViewport(w.viewPos.Add(delta))
}

// SetViewport positions the viewport to be at x,y
func (w *Window) SetViewport(pt intgeom.Point2) {
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

// ViewportBounds returns the boundary of this window's viewport, or the rectangle
// that the viewport is not allowed to exit as it moves around. It often represents
// the total size of the world within a given scene. If bounds are not enforced, ok will
// be false.
func (w *Window) ViewportBounds() (rect intgeom.Rect2, ok bool) {
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

	newView := rect.Clamp(w.viewPos)
	if newView != w.viewPos {
		w.SetViewport(newView)
	}
}

// Viewport returns the viewport's position. Its width and height are the window's
// width and height. This position plus width/height cannot exceed ViewportBounds.
func (w *Window) Viewport() intgeom.Point2 {
	return w.viewPos
}
