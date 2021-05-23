package window

import "github.com/oakmound/oak/v3/alg/intgeom"

type Window interface {
	SetFullScreen(bool) error
	SetBorderless(bool) error
	SetTopMost(bool) error
	SetTitle(string) error
	SetTrayIcon(string) error
	ShowNotification(title, msg string, icon bool) error
	MoveWindow(x, y, w, h int) error
	HideCursor() error
	//GetMonitorSize() (int, int)
	//Close() error
	Width() int
	Height() int
	Viewport() intgeom.Point2
	Quit()
	SetViewportBounds(intgeom.Rect2)
	NextScene()
	GoToScene(string)
	InFocus() bool
}
