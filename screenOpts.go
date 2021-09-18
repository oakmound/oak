package oak

import "github.com/oakmound/oak/v3/oakerr"

type fullScreenable interface {
	SetFullScreen(bool) error
}

// SetFullScreen attempts to set the local oak window to be full screen.
// If the window does not support this functionality, it will report as such.
func (w *Window) SetFullScreen(on bool) error {
	if fs, ok := w.windowControl.(fullScreenable); ok {
		return fs.SetFullScreen(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetFullScreen",
	}
}

type movableWindow interface {
	MoveWindow(x, y, w, h int32) error
}

// MoveWindow sets the position of a window to be x,y and it's dimensions to w,h
// If the window does not support being positioned, it will report as such.
func (w *Window) MoveWindow(x, y, wd, h int) error {
	if mw, ok := w.windowControl.(movableWindow); ok {
		return mw.MoveWindow(int32(x), int32(y), int32(wd), int32(h))
	}
	return oakerr.UnsupportedPlatform{
		Operation: "MoveWindow",
	}
}

type borderlesser interface {
	SetBorderless(bool) error
}

// SetBorderless attempts to set the local oak window to have no border.
// If the window does not support this functionaltiy, it will report as such.
func (w *Window) SetBorderless(on bool) error {
	if bs, ok := w.windowControl.(borderlesser); ok {
		return bs.SetBorderless(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetBorderless",
	}
}

type topMoster interface {
	SetTopMost(bool) error
}

// SetTopMost attempts to set the local oak window to stay on top of other windows.
// If the window does not support this functionality, it will report as such.
func (w *Window) SetTopMost(on bool) error {
	if tm, ok := w.windowControl.(topMoster); ok {
		return tm.SetTopMost(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTopMost",
	}
}

type titler interface {
	SetTitle(string) error
}

// SetTitle sets this window's title.
func (w *Window) SetTitle(title string) error {
	if t, ok := w.windowControl.(titler); ok {
		return t.SetTitle(title)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTitle",
	}
}

type trayIconer interface {
	SetTrayIcon(string) error
}

// SetTrayIcon sets a application tray icon for this program.
func (w *Window) SetTrayIcon(icon string) error {
	if t, ok := w.windowControl.(trayIconer); ok {
		return t.SetTrayIcon(icon)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTrayIcon",
	}
}

type trayNotifier interface {
	ShowNotification(title, msg string, icon bool) error
}

// ShowNotification shows a text notification, optionally using a previously set
// tray icon.
func (w *Window) ShowNotification(title, msg string, icon bool) error {
	if t, ok := w.windowControl.(trayNotifier); ok {
		return t.ShowNotification(title, msg, icon)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "ShowNotification",
	}
}

type cursorHider interface {
	HideCursor() error
}

// HideCursor disables showing the cursor when it is over this window.
func (w *Window) HideCursor() error {
	if t, ok := w.windowControl.(cursorHider); ok {
		return t.HideCursor()
	}
	return oakerr.UnsupportedPlatform{
		Operation: "HideCursor",
	}
}

type getCursorPositioner interface {
	GetCursorPosition() (x, y float64)
}

// GetCursorPosition returns the cusor position relative to the top left corner of this window.
func (w *Window) GetCursorPosition() (x, y float64, err error) {
	if wp, ok := w.windowControl.(getCursorPositioner); ok {
		x, y := wp.GetCursorPosition()
		return x, y, nil
	}
	return 0, 0, oakerr.UnsupportedPlatform{
		Operation: "GetCursorPosition",
	}
}
