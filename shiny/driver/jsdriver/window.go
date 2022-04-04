//go:build js
// +build js

package jsdriver

import (
	"image"
	"image/draw"
	"syscall/js"

	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type Window struct {
	screen *screenImpl
	cvs    *Canvas2D
	event.Deque
}

func (w *Window) Release() {}
func (w *Window) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	rgba := src.(*textureImpl).rgba
	js.CopyBytesToJS(w.cvs.copybuff, rgba.Pix)
	w.cvs.imgData.Get("data").Call("set", w.cvs.copybuff)
	w.cvs.ctx.Call("putImageData", w.cvs.imgData, 0, 0)
}
func (w *Window) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {}

func (w *Window) Publish() {}

func (w *Window) sendMouseEvent(mouseEvent js.Value, dir mouse.Direction) {
	x, y := mouseEvent.Get("offsetX"), mouseEvent.Get("offsetY")
	button := mouseEvent.Get("button")
	var mButton mouse.Button
	switch button.Int() {
	case 0:
		mButton = mouse.ButtonLeft
	case 1:
		mButton = mouse.ButtonMiddle
	case 2:
		mButton = mouse.ButtonRight
	default:
		mButton = mouse.ButtonNone
	}
	w.Send(mouse.Event{
		Button:    mButton,
		X:         float32(x.Float()),
		Y:         float32(y.Float()),
		Direction: dir,
	})
}

func (w *Window) sendKeyEvent(keyEvent js.Value, dir key.Direction) {
	var mods key.Modifiers
	if keyEvent.Get("shiftKey").Bool() {
		mods |= key.ModShift
	}
	if keyEvent.Get("metaKey").Bool() {
		mods |= key.ModMeta
	}
	if keyEvent.Get("ctrlKey").Bool() {
		mods |= key.ModControl
	}
	if keyEvent.Get("altKey").Bool() {
		mods |= key.ModAlt
	}
	if keyEvent.Get("repeat").Bool() {
		dir = key.DirNone
	}
	w.Send(key.Event{
		Modifiers: mods,
		Code:      parseKeyCode(keyEvent.Get("code").String()),
		Rune:      rune(keyEvent.Get("key").String()[0]),
		Direction: dir,
	})
}

func parseKeyCode(cd string) key.Code {
	switch cd {
	case "ArrowLeft":
		return key.CodeLeftArrow
	case "ArrowRight":
		return key.CodeRightArrow
	case "ArrowUp":
		return key.CodeUpArrow
	case "ArrowDown":
		return key.CodeDownArrow
	case "KeyA":
		return key.CodeA
	case "KeyB":
		return key.CodeB
	case "KeyC":
		return key.CodeC
	case "KeyD":
		return key.CodeD
	case "KeyE":
		return key.CodeE
	case "KeyF":
		return key.CodeF
	case "KeyG":
		return key.CodeG
	case "KeyH":
		return key.CodeH
	case "KeyI":
		return key.CodeI
	case "KeyJ":
		return key.CodeJ
	case "KeyK":
		return key.CodeK
	case "KeyL":
		return key.CodeL
	case "KeyM":
		return key.CodeM
	case "KeyN":
		return key.CodeN
	case "KeyO":
		return key.CodeO
	case "KeyP":
		return key.CodeP
	case "KeyQ":
		return key.CodeQ
	case "KeyR":
		return key.CodeR
	case "KeyS":
		return key.CodeS
	case "KeyT":
		return key.CodeT
	case "KeyU":
		return key.CodeU
	case "KeyV":
		return key.CodeV
	case "KeyW":
		return key.CodeW
	case "KeyX":
		return key.CodeX
	case "KeyY":
		return key.CodeY
	case "KeyZ":
		return key.CodeZ
	case "Enter":
		return key.CodeReturnEnter
	case "Escape":
		return key.CodeEscape
	case "Space":
		return key.CodeSpacebar
	case "ShiftLeft":
		return key.CodeLeftShift
	case "ShiftRight":
		return key.CodeRightShift
	case "Backspace":
		return key.CodeDeleteBackspace
	// TODO: more keys
	default:
		return key.CodeUnknown
	}
}
