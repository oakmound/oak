package oak

import (
	"github.com/oakmound/shiny/screen"
)

func windowController(s screen.Screen, x, y int32, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(conf.Title),
		screen.Position(x, y),
		screen.Fullscreen(SetupFullscreen),
		screen.Borderless(SetupBorderless),
		screen.TopMost(SetupTopMost),
	))
}
