package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	ViewX         = 0
	ViewY         = 0
	useViewBounds = false
	viewBounds    Rect
)

type Rect struct {
	minX, minY, maxX, maxY int
}

func SetScreen(x, y int) {
	dlog.Verb("Requesting ViewPoint ", x, y)
	viewportChannel <- [2]int{x, y}
	// updateScreen(x, y)
}

func updateScreen(x, y int) {
	if useViewBounds {
		if viewBounds.minX < x && viewBounds.maxX > x+ScreenWidth {
			dlog.Verb("Set ViewX to ", x)
			ViewX = x
		} else if viewBounds.minX > x {
			ViewX = viewBounds.minX
		} else if viewBounds.maxX < x+ScreenWidth {
			ViewX = viewBounds.maxX - ScreenWidth
		}
		if viewBounds.minY < y && viewBounds.maxY > y+ScreenHeight {
			dlog.Verb("Set ViewY to ", y)
			ViewY = y
		} else if viewBounds.minY > y {
			ViewY = viewBounds.minY
		} else if viewBounds.maxY < y+ScreenHeight {
			ViewY = viewBounds.maxY - ScreenHeight
		}
	} else {
		dlog.Verb("Set ViewXY to ", x, " ", y)
		ViewX = x
		ViewY = y
	}
	eb.Trigger("ViewportUpdate", []float64{float64(ViewX), float64(ViewY)})
	dlog.Verb("ViewX, Y: ", ViewX, " ", ViewY)
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
		ViewX += dX * speed
		ViewY += dY * speed
		return 0
	}
}
