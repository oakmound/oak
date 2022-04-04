// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x11driver

// TODO: implement a back buffer.

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/render"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/driver/internal/lifecycler"
	"github.com/oakmound/oak/v3/shiny/driver/internal/x11"
	"github.com/oakmound/oak/v3/shiny/driver/internal/x11key"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

type windowImpl struct {
	s *screenImpl

	xw xproto.Window
	xg xproto.Gcontext
	xp render.Picture

	event.Deque
	xevents chan xgb.Event

	// This next group of variables are mutable, but are only modified in the
	// screenImpl.run goroutine.
	width, height uint32

	lifecycler lifecycler.State

	mu sync.Mutex

	x, y     uint32
	released bool
}

func (w *windowImpl) Release() {
	w.mu.Lock()
	released := w.released
	w.released = true
	w.mu.Unlock()

	// TODO: call w.lifecycler.SetDead and w.lifecycler.SendEvent, a la
	// handling atomWMDeleteWindow?

	if released {
		return
	}
	render.FreePicture(w.s.xc, w.xp)
	xproto.FreeGC(w.s.xc, w.xg)
	xproto.DestroyWindow(w.s.xc, w.xw)
}

func (w *windowImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	src.(*bufferImpl).upload(xproto.Drawable(w.xw), w.xg, w.s.xsi.RootDepth, dp, sr)
}

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	fill(w.s.xc, w.xp, dr, src, op)
}

func (w *windowImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {
	w.s.drawUniform(w.xp, &src2dst, src, sr, op)
}

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	src.(*textureImpl).draw(w.xp, &src2dst, sr, op)
}

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Copy(w, dp, src, sr, op)
}

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, op)
}

func (w *windowImpl) Publish() screen.PublishResult {
	// TODO: implement a back buffer, and copy or flip that here to the front
	// buffer.

	// This sync isn't needed to flush the outgoing X11 requests. Instead, it
	// acts as a form of flow control. Outgoing requests can be quite small on
	// the wire, e.g. draw this texture ID (an integer) to this rectangle (four
	// more integers), but much more expensive on the server (blending a
	// million source and destination pixels). Without this sync, the Go X11
	// client could easily end up sending work at a faster rate than the X11
	// server can serve.
	w.s.xc.Sync()

	return screen.PublishResult{}
}

func (w *windowImpl) SetFullScreen(fullscreen bool) error {
	return x11.SetFullScreen(w.s.XUtil, w.xw, fullscreen)
}

func (w *windowImpl) SetBorderless(borderless bool) error {
	return x11.SetBorderless(w.s.XUtil, w.xw, borderless)
}

func (w *windowImpl) handleConfigureNotify(ev xproto.ConfigureNotifyEvent) {
	// TODO: does the order of these lifecycle and size events matter? Should
	// they really be a single, atomic event?
	w.lifecycler.SetVisible((int(ev.X)+int(ev.Width)) > 0 && (int(ev.Y)+int(ev.Height)) > 0)
	w.lifecycler.SendEvent(w, nil)

	newWidth, newHeight := uint32(ev.Width), uint32(ev.Height)
	if w.width == newWidth && w.height == newHeight {
		return
	}
	w.width, w.height = newWidth, newHeight
	w.Send(size.Event{
		WidthPx:     int(newWidth),
		HeightPx:    int(newHeight),
		WidthPt:     geom.Pt(newWidth),
		HeightPt:    geom.Pt(newHeight),
		PixelsPerPt: w.s.pixelsPerPt,
	})
}

func (w *windowImpl) handleExpose() {
	w.Send(paint.Event{})
}

func (w *windowImpl) handleKey(detail xproto.Keycode, state uint16, dir key.Direction) {
	r, c := w.s.keysyms.Lookup(uint8(detail), state, w.s.numLockMod)
	w.Send(key.Event{
		Rune:      r,
		Code:      c,
		Modifiers: x11key.KeyModifiers(state),
		Direction: dir,
	})
}

func (w *windowImpl) handleMouse(x, y int16, b xproto.Button, state uint16, dir mouse.Direction) {
	// TODO: should a mouse.Event have a separate MouseModifiers field, for
	// which buttons are pressed during a mouse move?
	btn := mouse.Button(b)
	switch btn {
	case 4:
		btn = mouse.ButtonWheelUp
	case 5:
		btn = mouse.ButtonWheelDown
	case 6:
		btn = mouse.ButtonWheelLeft
	case 7:
		btn = mouse.ButtonWheelRight
	}
	if btn.IsWheel() {
		if dir != mouse.DirPress {
			return
		}
		dir = mouse.DirStep
	}
	w.Send(mouse.Event{
		X:         float32(x),
		Y:         float32(y),
		Button:    btn,
		Modifiers: x11key.KeyModifiers(state),
		Direction: dir,
	})
}

func (w *windowImpl) MoveWindow(x, y, width, height int32) error {
	newX, newY, newW, newH := x11.MoveWindow(w.s.xc, w.xw, x, y, width, height)
	w.x = uint32(newX)
	w.y = uint32(newY)
	w.width = uint32(newW)
	w.height = uint32(newH)
	w.Send(size.Event{
		WidthPx:     int(newW),
		HeightPx:    int(newH),
		WidthPt:     geom.Pt(newW),
		HeightPt:    geom.Pt(newH),
		PixelsPerPt: w.s.pixelsPerPt,
	})
	return nil
}

func (w *Window) SetTitle(title string) error {
	xproto.ChangeProperty(w.s.xc, xproto.PropModeReplace, w.xw,
		w.s.atoms["_NET_WM_NAME"], w.s.atoms["UTF8_STRING"],
		8, uint32(len(title)), []byte(title))
	return nil
}

func (w *Window) SetTopMost(topMost bool) error {
	return x11.SetTopMost(w.s.XUtil, w.xw, topMost)
}

func (w *Window) SetIcon(icon image.Image) error {
	bds := icon.Bounds()
	wd := bds.Max.X
	h := bds.Max.Y
	u32w := uint32(wd)
	u32h := uint32(h)
	// 4 bytes, b/g/r/a, per pixel
	bgra := make([]byte, 8, 8+wd*h*4)
	// prepend width and height
	bgra[0] = byte(u32w)
	bgra[1] = byte(u32w >> 8)
	bgra[2] = byte(u32w >> 16)
	bgra[3] = byte(u32w >> 24)
	bgra[4] = byte(u32h)
	bgra[5] = byte(u32h >> 8)
	bgra[6] = byte(u32h >> 16)
	bgra[7] = byte(u32h >> 24)
	for x := 0; x < wd; x++ {
		for y := 0; y < h; y++ {
			c := icon.At(x, (h-1)-y)
			r, g, b, a := c.RGBA()
			bgra = append(bgra, byte(b>>8))
			bgra = append(bgra, byte(g>>8))
			bgra = append(bgra, byte(r>>8))
			bgra = append(bgra, byte(a>>8))
		}
	}
	const XA_CARDINAL = 6

	xproto.ChangeProperty(w.s.xc, xproto.PropModeReplace, w.xw,
		w.s.atoms["_NET_WM_ICON"], XA_CARDINAL,
		32, uint32(len(bgra))/4, bgra)
	return nil
}
