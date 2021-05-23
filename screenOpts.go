package oak

import "github.com/oakmound/oak/v3/oakerr"

type fullScreenable interface {
	SetFullScreen(bool) error
}

// SetFullScreen attempts to set the local oak window to be full screen.
// If the window does not support this functionality, it will report as such.
func (c *Controller) SetFullScreen(on bool) error {
	if fs, ok := c.windowControl.(fullScreenable); ok {
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
func (c *Controller) MoveWindow(x, y, w, h int) error {
	if mw, ok := c.windowControl.(movableWindow); ok {
		return mw.MoveWindow(int32(x), int32(y), int32(w), int32(h))
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
func (c *Controller) SetBorderless(on bool) error {
	if bs, ok := c.windowControl.(borderlesser); ok {
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
func (c *Controller) SetTopMost(on bool) error {
	if tm, ok := c.windowControl.(topMoster); ok {
		return tm.SetTopMost(on)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTopMost",
	}
}

type titler interface {
	SetTitle(string) error
}

func (c *Controller) SetTitle(title string) error {
	if t, ok := c.windowControl.(titler); ok {
		return t.SetTitle(title)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTitle",
	}
}

type trayIconer interface {
	SetTrayIcon(string) error
}

func (c *Controller) SetTrayIcon(icon string) error {
	if t, ok := c.windowControl.(trayIconer); ok {
		return t.SetTrayIcon(icon)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "SetTrayIcon",
	}
}

type trayNotifier interface {
	ShowNotification(title, msg string, icon bool) error
}

func (c *Controller) ShowNotification(title, msg string, icon bool) error {
	if t, ok := c.windowControl.(trayNotifier); ok {
		return t.ShowNotification(title, msg, icon)
	}
	return oakerr.UnsupportedPlatform{
		Operation: "ShowNotification",
	}
}

type cursorHider interface {
	HideCursor() error
}

func (c *Controller) HideCursor() error {
	if t, ok := c.windowControl.(cursorHider); ok {
		return t.HideCursor()
	}
	return oakerr.UnsupportedPlatform{
		Operation: "HideCursor",
	}
}
