package oak

import (
	"image"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
)

var (
	// ViewPos represents the point in the world which the viewport is anchored at.
	ViewPos       = image.Point{}
	useViewBounds = false
	viewBounds    rect
)

type rect struct {
	minX, minY, maxX, maxY int
}

// SetScreen sends a signal to the draw loop to set the viewport to be at x,y
func SetScreen(x, y int) {
	dlog.Verb("Requesting ViewPoint ", x, y)
	viewportCh <- [2]int{x, y}
}

var (
	defaultUpdateScreen = func(x, y int) {
		if useViewBounds {
			if viewBounds.minX <= x && viewBounds.maxX >= x+ScreenWidth {
				dlog.Verb("Set ViewX to ", x)
				ViewPos.X = x
			} else if viewBounds.minX > x {
				ViewPos.X = viewBounds.minX
			} else if viewBounds.maxX < x+ScreenWidth {
				ViewPos.X = viewBounds.maxX - ScreenWidth
			}
			if viewBounds.minY <= y && viewBounds.maxY >= y+ScreenHeight {
				dlog.Verb("Set ViewY to ", y)
				ViewPos.Y = y
			} else if viewBounds.minY > y {
				ViewPos.Y = viewBounds.minY
			} else if viewBounds.maxY < y+ScreenHeight {
				ViewPos.Y = viewBounds.maxY - ScreenHeight
			}
		} else {
			dlog.Verb("Set ViewXY to ", x, " ", y)
			ViewPos = image.Point{x, y}
		}
		eb.Trigger("ViewportUpdate", []float64{float64(ViewPos.X), float64(ViewPos.Y)})
		dlog.Verb("ViewX, Y: ", ViewPos.X, " ", ViewPos.Y)
	}
	updateScreen = defaultUpdateScreen
)

// SetViewportBounds sets the minimum and maximum position of the viewport, including
// screen dimensions
func SetViewportBounds(x1, y1, x2, y2 int) {
	if x2 < ScreenWidth {
		x2 = ScreenWidth
	}
	if y2 < ScreenHeight {
		y2 = ScreenHeight
	}
	useViewBounds = true
	viewBounds = rect{x1, y1, x2, y2}

	dlog.Info("Viewport bounds set to, ", x1, y1, x2, y2)

	newViewX := ViewPos.X
	newViewY := ViewPos.Y
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

	if newViewX != ViewPos.X || newViewY != ViewPos.Y {
		viewportCh <- [2]int{newViewX, newViewY}
	}
}

func moveViewportBinding(speed int) func(int, interface{}) int {
	return func(cID int, n interface{}) int {
		dX := 0
		dY := 0
		if IsDown("UpArrow") {
			dY--
		}
		if IsDown("DownArrow") {
			dY++
		}
		if IsDown("LeftArrow") {
			dX--
		}
		if IsDown("RightArrow") {
			dX++
		}
		ViewPos.X += dX * speed
		ViewPos.Y += dY * speed
		if viewportLocked {
			return event.UnbindSingle
		}
		return 0
	}
}

// ViewVector returns ViewPos as a Vector
func ViewVector() physics.Vector {
	return physics.NewVector(float64(ViewPos.X), float64(ViewPos.Y))
}
