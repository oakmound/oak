package oak

import (
	"github.com/oakmound/oak/oakerr"
)

type FullScreenable interface {
	SetFullScreen(bool)
}

func SetFullScreen(on bool) error {
	if fs, ok := windowControl.(FullScreenable); ok {
		fs.SetFullScreen(on)
		return nil
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetFullScreen",
	}
}

type PositionableWindow interface {
	MoveWindow(x, y, w, h int32)
}

func MoveWindow(x, y, w, h int) error {
	if mw, ok := windowControl.(PositionableWindow); ok {
		mw.MoveWindow(int32(x), int32(y), int32(w), int32(h))
		return nil
	}
	return oakerr.UnsupportedPlatform{
		Operation: "MoveWindow",
	}
}
