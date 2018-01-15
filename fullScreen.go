package oak

import (
	"github.com/oakmound/oak/oakerr"
)

// FullScreenable defines windows that can be set to full screen.
type FullScreenable interface {
	SetFullScreen(bool)
}

// SetFullScreen attempts to set the local oak window to be full screen.
// If the window does not support this functionality, it will report as such.
func SetFullScreen(on bool) error {
	if fs, ok := windowControl.(FullScreenable); ok {
		fs.SetFullScreen(on)
		return nil
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetFullScreen",
	}
}

// MovableWindow defines windows that can have their position set
type MovableWindow interface {
	MoveWindow(x, y, w, h int32)
}

// MoveWindow sets the position of a window to be x,y and it's dimensions to w,h
// If the window does not support being positioned, it will report as such.
func MoveWindow(x, y, w, h int) error {
	if mw, ok := windowControl.(MovableWindow); ok {
		mw.MoveWindow(int32(x), int32(y), int32(w), int32(h))
		return nil
	}
	return oakerr.UnsupportedPlatform{
		Operation: "MoveWindow",
	}
}
