// This can likely build on more than windows but we haven't tested it

// +build windows,linux

package oak

import (
	"golang.org/x/exp/shiny/screen"
)

func WindowController(s screen.Screen, ScreenWidth, ScreenHeight int) (screen.Window, error) {
	return s.NewWindow(&screen.NewWindowOptions{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Title:  conf.Title,
	})
}

// OK so I've been thinking about platform support
// And I realized that shiny actually gave us some nice interfaces to work with
// If we can define a screen.Window satisfying struct which can interact with any
// platform not supported by Shiny (~~JS~~), we can draw our images to that
// platform.
// Given we can do that, all of our code has been split up such that providing a
// WindowController function in an appropriate build file should do the job of
// supporting other platforms.
