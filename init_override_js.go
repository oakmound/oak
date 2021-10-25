//go:build js

package oak

import (
	"github.com/oakmound/oak/v3/dlog"
	"syscall/js"
)

func overrideInit(w *Window) {
	w.DrawTicker.Stop()
	if w.DrawFrameRate != 60 {
		dlog.Info("Ignoring draw frame rate in JS")
	}
	w.animationFrame = make(chan struct{})
	js.Global().Call("requestAnimationFrame", js.FuncOf(w.requestFrame))
}

func (w *Window) requestFrame(this js.Value, args []js.Value) interface{} {
	w.animationFrame <- struct{}{}
	js.Global().Call("requestAnimationFrame", js.FuncOf(w.requestFrame))
	return nil
}
