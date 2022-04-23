//go:build js
// +build js

package oak

import (
	"github.com/oakmound/oak/v4/dlog"
	"syscall/js"
)

func overrideInit(w *Window) {
	w.DrawTicker.Stop()
	if w.DrawFrameRate != 60 {
		dlog.Info("Ignoring draw frame rate in JS")
	}
	if w.config.EnableDebugConsole {
		dlog.Info("Debug console is not supported in JS")
		w.config.EnableDebugConsole = false
	}
	if w.config.UnlimitedDrawFrameRate {
		dlog.Info("Unlimited draw frame rate is not supported in JS")
		w.config.UnlimitedDrawFrameRate = false
	}
	w.animationFrame = make(chan struct{})
	js.Global().Call("requestAnimationFrame", js.FuncOf(w.requestFrame))
}

func (w *Window) requestFrame(this js.Value, args []js.Value) interface{} {
	w.animationFrame <- struct{}{}
	js.Global().Call("requestAnimationFrame", js.FuncOf(w.requestFrame))
	return nil
}
