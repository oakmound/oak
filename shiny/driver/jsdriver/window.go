package jsdriver

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"syscall/js"

	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type windowImpl struct {
	screen *screenImpl
	cvs    *Canvas2D
	event.Deque
}

func (w *windowImpl) Release()                                                                      {}
func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op)     {}
func (w *windowImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {}
func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op)       {}
func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	rgba := src.(*textureImpl).rgba
	js.CopyBytesToJS(w.cvs.copybuff, rgba.Pix)
	w.cvs.imgData.Get("data").Call("set", w.cvs.copybuff)
	w.cvs.ctx.Call("putImageData", w.cvs.imgData, 0, 0)
}
func (w *windowImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {}
func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)        {}

func (w *windowImpl) Publish() screen.PublishResult {
	return screen.PublishResult{}
}

func (w *windowImpl) sendMouseEvent(mouseEvent js.Value, dir mouse.Direction) {
	x, y := mouseEvent.Get("offsetX"), mouseEvent.Get("offsetY")
	w.Send(mouse.Event{
		X:         float32(x.Float()),
		Y:         float32(y.Float()),
		Direction: dir,
	})
}

func (w *windowImpl) sendKeyEvent(keyEvent js.Value, dir key.Direction) {
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
	// TODO: more keys
	default:
		fmt.Println("unknown key", cd)
		return key.CodeUnknown
	}
}
