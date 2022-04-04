//go:build (windows || linux || osx) && !js && !android
// +build windows linux osx
// +build !js
// +build !android

package oak

import (
	"image"
)

// MoveWindow calls MoveWindow on the default window.
func MoveWindow(x, y, w, h int) error {
	initDefaultWindow()
	return defaultWindow.MoveWindow(x, y, w, h)
}

// SetFullScreen calls SetFullScreen on the default window.
func SetFullScreen(fs bool) error {
	initDefaultWindow()
	return defaultWindow.SetFullScreen(fs)
}

// SetBorderless calls SetBorderless on the default window.
func SetBorderless(bs bool) error {
	initDefaultWindow()
	return defaultWindow.SetBorderless(bs)
}

// SetTopMost calls SetTopMost on the default window.
func SetTopMost(on bool) error {
	initDefaultWindow()
	return defaultWindow.SetTopMost(on)
}

// SetTitle calls SetTitle on the default window.
func SetTitle(title string) error {
	initDefaultWindow()
	return defaultWindow.SetTitle(title)
}

// SetIcon calls SetIcon on the default window.
func SetIcon(icon image.Image) error {
	initDefaultWindow()
	return defaultWindow.SetIcon(icon)
}

// HideCursor calls HideCursor on the default window.
func HideCursor() error {
	initDefaultWindow()
	return defaultWindow.HideCursor()
}

// GetCursorPosition calls GetCursorPosition on the default window.
func GetCursorPosition() (x, y float64) {
	initDefaultWindow()
	return defaultWindow.GetCursorPosition()
}
