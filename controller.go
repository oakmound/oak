package oak

import (
	"github.com/oakmound/shiny/screen"
)

func windowController(s screen.Screen, ScreenWidth, ScreenHeight int) (screen.Window, error) {
	return s.NewWindow(&screen.NewWindowOptions{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Title:  conf.Title,
	})
}
