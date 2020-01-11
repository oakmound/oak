package oak

import (
	"github.com/oakmound/oak/v2/oakerr"
)

// FullScreenable defines windows that can be set to full screen.
type FullScreenable interface {
	SetFullScreen(bool) error
}

// SetFullScreen attempts to set the local oak window to be full screen.
// If the window does not support this functionality, it will report as such.
func SetFullScreen(on bool) error {
	if fs, ok := windowControl.(FullScreenable); ok {
		return fs.SetFullScreen(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetFullScreen",
	}
}

// MovableWindow defines windows that can have their position set
type MovableWindow interface {
	MoveWindow(x, y, w, h int32) error
}

// MoveWindow sets the position of a window to be x,y and it's dimensions to w,h
// If the window does not support being positioned, it will report as such.
func MoveWindow(x, y, w, h int) error {
	if mw, ok := windowControl.(MovableWindow); ok {
		return mw.MoveWindow(int32(x), int32(y), int32(w), int32(h))
	}
	return oakerr.UnsupportedPlatform{
		Operation: "MoveWindow",
	}
}

// A Borderlesser is a window that can have its border removed or replaced after
// removal.
type Borderlesser interface {
	SetBorderless(bool) error
}

// SetBorderless attempts to set the local oak window to have no border.
// If the window does not support this functionaltiy, it will report as such.
func SetBorderless(on bool) error {
	if bs, ok := windowControl.(Borderlesser); ok {
		return bs.SetBorderless(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetBorderless",
	}
}

// A TopMoster is a window that can be configured to stay on top of other windows.
type TopMoster interface {
	SetTopMost(bool) error
}

// SetTopMost attempts to set the local oak window to stay on top of other windows.
// If the window does not support this functionality, it will report as such.
func SetTopMost(on bool) error {
	if tm, ok := windowControl.(TopMoster); ok {
		return tm.SetTopMost(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTopMost",
	}
}
