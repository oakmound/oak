package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	ViewX         = 0
	ViewY         = 0
	useViewBounds = false
	viewBounds    []int
)

func SetScreen(x, y int) {
	if useViewBounds {
		if viewBounds[0] < x && viewBounds[2] > x+ScreenWidth {
			dlog.Verb("Set ViewX to ", x)
			ViewX = x
		} else if viewBounds[0] > x {
			ViewX = viewBounds[0]
		} else if viewBounds[2] < x+ScreenWidth {
			ViewX = viewBounds[2] - ScreenWidth
		}
		if viewBounds[1] < y && viewBounds[3] > y+ScreenHeight {
			dlog.Verb("Set ViewY to ", y)
			ViewY = y
		} else if viewBounds[1] > y {
			ViewY = viewBounds[1]
		} else if viewBounds[3] < y+ScreenHeight {
			ViewY = viewBounds[3] - ScreenHeight
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
	viewBounds = []int{x1, y1, x2, y2}
}
