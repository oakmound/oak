package oak

import (
	"image"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
)

var (
	ViewPos       = image.Point{}
	useViewBounds = false
	viewBounds    Rect
)

type Rect struct {
	minX, minY, maxX, maxY int
}

func SetScreen(x, y int) {
	dlog.Verb("Requesting ViewPoint ", x, y)
	viewportChannel <- [2]int{x, y}
}

func updateScreen(x, y int) {
	if useViewBounds {
		if viewBounds.minX < x && viewBounds.maxX > x+ScreenWidth {
			dlog.Verb("Set ViewX to ", x)
			ViewPos.X = x
		} else if viewBounds.minX > x {
			ViewPos.X = viewBounds.minX
		} else if viewBounds.maxX < x+ScreenWidth {
			ViewPos.X = viewBounds.maxX - ScreenWidth
		}
		if viewBounds.minY < y && viewBounds.maxY > y+ScreenHeight {
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

func SetViewportBounds(x1, y1, x2, y2 int) {
	if x2 < ScreenWidth {
		x2 = ScreenWidth
	}
	if y2 < ScreenHeight {
		y2 = ScreenHeight
	}
	dlog.Info("Viewport bounds set to, ", x1, y1, x2, y2)
	useViewBounds = true
	viewBounds = Rect{x1, y1, x2, y2}
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
			return event.UNBIND_SINGLE
		}
		return 0
	}
}
