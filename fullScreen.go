package oak

import "errors"

type FullScreenable interface {
	SetFullScreen()
}

func SetFullScreen() error {
	if fs, ok := windowControl.(FullScreenable); ok {
		fs.SetFullScreen()
		return nil
	}
	return errors.New("Fullscreen not supported on this platform")
}

type PositionableWindow interface {
	MoveWindow(x, y, w, h int32)
}

func MoveWindow(x, y, w, h int) error {
	if mw, ok := windowControl.(PositionableWindow); ok {
		mw.MoveWindow(int32(x), int32(y), int32(w), int32(h))
		return nil
	}
	return errors.New("Window movement not supported on this platform")
}
