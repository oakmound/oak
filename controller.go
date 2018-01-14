package oak

import (
	"github.com/oakmound/shiny/screen"
)

func windowController(s screen.Screen, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(conf.Title),
	))
}
