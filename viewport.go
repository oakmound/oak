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
	dlog.Verb("ViewX, Y: ", ViewX, " ", ViewY)
	eb.Trigger("ViewportUpdate", []float64{float64(ViewX), float64(ViewY)})
}

func SetViewportBounds(x1, y1, x2, y2 int) {
	dlog.Info("Viewport bounds set to, ", x1, y1, x2, y2)
	useViewBounds = true
	viewBounds = Rect{x1, y1, x2, y2}
}
