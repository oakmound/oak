// Package window provides a common interface for oak-created windows.
package window

import (
	"image"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
)

// Window is an interface of methods on an oak.Window available on platforms which have distinct app windows
// (osx, linux, windows). It is not available on other platforms (js, android)
type Window interface {
	App

	SetFullScreen(bool) error
	SetBorderless(bool) error
	SetTopMost(bool) error
	SetTitle(string) error
	SetIcon(image.Image) error
	MoveWindow(x, y, w, h int) error
	HideCursor() error
}

// App is an interface of methods available to all oak programs.
type App interface {
	Bounds() intgeom.Point2

	Viewport() intgeom.Point2
	SetViewportBounds(intgeom.Rect2)
	ShiftViewport(intgeom.Point2)
	SetViewport(intgeom.Point2)

	NextScene()
	GoToScene(string)

	InFocus() bool
	Quit()

	EventHandler() event.Handler
}
