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

	// SetFullscreen causes a window to expand and fill a display.
	SetFullScreen(bool) error
	// SetBorderless causes a window to lose its OS-provided border definitions, e.g. window title, close button.
	SetBorderless(bool) error
	// SetTopMost causes a window to remain above other windows even when it is clicked out of.
	SetTopMost(bool) error
	// SetTitle changes the title of this window, usually displayed in the top left of the window next to the icon.
	SetTitle(string) error
	// SetIcon changes the icon of this window, usually displayed both in the top left of the window and in a taskbar
	// component.
	SetIcon(image.Image) error
	// MoveWindow moves a window to the given x,y coordinates with the given dimensions.
	// TODO v4: intgeom.Rect2?
	MoveWindow(x, y, w, h int) error
	// HideCursor will cause the mouse cursor to not display when it lies within this window.
	HideCursor() error
}

// App is an interface of methods available to all oak programs.
type App interface {
	// Bounds returns the boundaries of the application client area measured in pixels. This is not the size
	// of the window or app on the operating system necessarily; it is the area able to be rendered to within oak.
	// On some platforms these two concepts will usually be equal (js); on some they will have a built in scaling factor
	// (osx, for retina displays), and if a window is manually scaled by a user and oak is not instructed to resize to
	// match the scale, this area will be unchanged and the view will be stretched to fit the window.
	Bounds() intgeom.Point2

	// Viewport relates Bounds() to the entire content available for display. Viewport returns where the top left corner
	// of the application client area is.
	Viewport() intgeom.Point2
	// SetViewportBounds defines the limits of where the viewport may be positioned. In other words, the total viewable
	// content of a scene. Unless impossible, the rectangle (viewport, viewport+bounds) will never leave the area defined
	// by SetViewportBounds.
	SetViewportBounds(intgeom.Rect2)
	// ShiftViewport is a helper method calling a.SetViewport(a.Viewport()+delta)
	ShiftViewport(delta intgeom.Point2)
	// SetViewport changes where the viewport position. If the resulting rectangle (viewport, viewport+bounds) would
	// exceed the boundary set by SetViewportBounds, viewport will be clamped to the edges of that boundary.
	SetViewport(intgeom.Point2)

	// NextScene causes the End function to be triggered for the current scene.
	NextScene()
	// GoToScene causes the End function to be triggered for the current scene, overriding the next scene to start.
	GoToScene(string)

	// InFocus returns whether the application is currently focused on, by whatever definition the OS has for an
	// application being in focus. For example, on linux/osx/windows a window is in focus once it is clicked on
	// and out of focus after another window is clicked on.
	InFocus() bool
	// Quit causes the app to cleanly exit. The current scene will not call it's End function.
	Quit()

	// EventHandler returns this app's active event handler.
	EventHandler() event.Handler
}
