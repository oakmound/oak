// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

// TODO: implement a back buffer.

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
	"syscall"
	"unsafe"

	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/driver/internal/win32"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/sys/windows"
)

var (
	windowLock sync.RWMutex
	allWindows = make(map[win32.HWND]*windowImpl)
)

type windowImpl struct {
	hwnd win32.HWND

	event.Deque

	sz             size.Event
	lifecycleStage lifecycle.Stage
	// Todo: the windows api is confused about
	// whether styles are int32 or uint32s.
	style, exStyle int32
	fullscreen     bool
	borderless     bool
	maximized      bool
	windowRect     *win32.RECT
	clientRect     *win32.RECT

	// guid is set on intialization and converted to trayGUID
	// when the tray icon is created.
	guid     [16]byte
	trayGUID *win32.GUID
}

func (w *windowImpl) Release() {
	if w.trayGUID != nil {
		iconData := win32.NOTIFYICONDATA{}
		iconData.CbSize = uint32(unsafe.Sizeof(iconData))
		iconData.UFlags = win32.NIF_GUID
		iconData.HWnd = w.hwnd
		iconData.GUIDItem = *w.trayGUID
		win32.Shell_NotifyIcon(win32.NIM_DELETE, &iconData)
	}
	win32.Release(win32.HWND(w.hwnd))
}

func (w *windowImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	src.(*bufferImpl).preUpload()
	defer src.(*bufferImpl).postUpload()

	w.execCmd(&cmd{
		id:     cmdUpload,
		dp:     dp,
		buffer: src.(*bufferImpl),
		sr:     sr,
	})
}

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	w.execCmd(&cmd{
		id:    cmdFill,
		dr:    dr,
		color: src,
		op:    op,
	})
}

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	if op != draw.Src && op != draw.Over {
		// TODO:
		return
	}
	w.execCmd(&cmd{
		id:      cmdDraw,
		src2dst: src2dst,
		texture: src.(*textureImpl).bitmap,
		sr:      sr,
		op:      op,
	})
}

func (w *windowImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {
	if op != draw.Src && op != draw.Over {
		return
	}
	w.execCmd(&cmd{
		id:      cmdDrawUniform,
		src2dst: src2dst,
		color:   src,
		sr:      sr,
		op:      op,
	})
}

func (w *windowImpl) SetTitle(title string) error {
	win32.SetWindowText(w.hwnd, title)
	return nil
}

func (w *windowImpl) SetBorderless(borderless bool) error {
	// Don't set borderless if currently fullscreen.
	if !w.fullscreen && borderless != w.borderless {
		if !w.borderless {
			// We don't need to get these values when w.borderless is true
			// because scaling is impossible without a border to grab to scale.
			// Todo: except through programatic window resizing.
			w.windowRect, _ = win32.GetWindowRect(w.hwnd)
			w.clientRect, _ = win32.GetClientRect(w.hwnd)
		}
		w.borderless = borderless
		if borderless {
			win32.SetWindowLong(w.hwnd, win32.GWL_STYLE,
				w.style & ^(win32.WS_CAPTION|win32.WS_THICKFRAME))
			win32.SetWindowLong(w.hwnd, win32.GWL_EXSTYLE,
				w.exStyle&^(win32.WS_EX_DLGMODALFRAME|
					win32.WS_EX_WINDOWEDGE|win32.WS_EX_CLIENTEDGE|win32.WS_EX_STATICEDGE))

			leftOffset := ((w.windowRect.Right - w.windowRect.Left) - w.clientRect.Right) / 2
			// The -leftOffset here is assuming that the bottom border has the same
			// height as the left and right do width.
			topOffset := ((w.windowRect.Bottom - w.windowRect.Top) - w.clientRect.Bottom) - leftOffset

			win32.SetWindowPos(w.hwnd, 0, w.windowRect.Left+leftOffset, w.windowRect.Top+topOffset,
				w.clientRect.Right, w.clientRect.Bottom,
				win32.SWP_NOZORDER|win32.SWP_NOACTIVATE|win32.SWP_FRAMECHANGED)

		} else {
			win32.SetWindowLong(w.hwnd, win32.GWL_STYLE, w.style)
			win32.SetWindowLong(w.hwnd, win32.GWL_EXSTYLE, w.exStyle)

			// On restore, resize to the previous saved rect size.
			win32.SetWindowPos(w.hwnd, 0, w.windowRect.Left, w.windowRect.Top,
				w.windowRect.Right-w.windowRect.Left, w.windowRect.Bottom-w.windowRect.Top,
				win32.SWP_NOZORDER|win32.SWP_NOACTIVATE|win32.SWP_FRAMECHANGED)

		}
		return nil
	}
	if w.fullscreen {
		return errors.New("cannot combine borderless and fullscreen")
	}
	return nil
}

func (w *windowImpl) SetFullScreen(fullscreen bool) error {
	if w.borderless {
		return errors.New("cannot combine borderless and fullscreen")
	}
	// Fullscreen impl copied from chromium
	// https://src.chromium.org/viewvc/chrome/trunk/src/ui/views/win/fullscreen_handler.cc
	// Save current window state if not already fullscreen.
	if !w.fullscreen {
		// Save current window information.  We force the window into restored mode
		// before going fullscreen because Windows doesn't seem to hide the
		// taskbar if the window is in the maximized state.
		w.maximized = win32.IsZoomed(w.hwnd)
		if w.maximized {
			win32.SendMessage(w.hwnd, win32.WM_SYSCOMMAND, win32.SC_RESTORE, 0)
		}
		w.windowRect, _ = win32.GetWindowRect(w.hwnd)
	}

	w.fullscreen = fullscreen

	if w.fullscreen {
		// Set new window style and size.
		win32.SetWindowLong(w.hwnd, win32.GWL_STYLE,
			w.style & ^(win32.WS_CAPTION|win32.WS_THICKFRAME))
		win32.SetWindowLong(w.hwnd, win32.GWL_EXSTYLE,
			w.exStyle&^(win32.WS_EX_DLGMODALFRAME|
				win32.WS_EX_WINDOWEDGE|win32.WS_EX_CLIENTEDGE|win32.WS_EX_STATICEDGE))
		// On expand, if we're given a window_rect, grow to it, otherwise do
		// not resize.
		// shiny cmt: Need to look into what this for_metro argument means,
		// right now we don't use it
		// if (!for_metro) {
		monitorInfo := win32.MONITORINFO{}
		monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
		win32.GetMonitorInfo(win32.MonitorFromWindow(w.hwnd, win32.MONITOR_DEFAULTTONEAREST),
			&monitorInfo)
		windowRect := monitorInfo.RcMonitor
		win32.SetWindowPos(w.hwnd, 0, windowRect.Left, windowRect.Top,
			windowRect.Right-windowRect.Left, windowRect.Bottom-windowRect.Top,
			win32.SWP_NOZORDER|win32.SWP_NOACTIVATE|win32.SWP_FRAMECHANGED)
		// }
	} else {
		// Reset original window style and size.  The multiple window size/moves
		// here are ugly, but if SetWindowPos() doesn't redraw, the taskbar won't be
		// repainted.  Better-looking methods welcome.
		win32.SetWindowLong(w.hwnd, win32.GWL_STYLE, w.style)
		win32.SetWindowLong(w.hwnd, win32.GWL_EXSTYLE, w.exStyle)

		// if !for_metro {
		// On restore, resize to the previous saved rect size.
		newRect := w.windowRect
		win32.SetWindowPos(w.hwnd, 0, newRect.Left, newRect.Top,
			newRect.Right-newRect.Left, newRect.Bottom-newRect.Top,
			win32.SWP_NOZORDER|win32.SWP_NOACTIVATE|win32.SWP_FRAMECHANGED)
		//}
		if w.maximized {
			win32.SendMessage(w.hwnd, win32.WM_SYSCOMMAND, win32.SC_MAXIMIZE, 0)
		}
	}
	return nil
}

// HideCursor turns the OS cursor into a 1x1 transparent image.
func (w *windowImpl) HideCursor() error {
	emptyCursor := win32.GetEmptyCursor()
	success := win32.SetClassLongPtr(w.hwnd, win32.GCLP_HCURSOR, uintptr(emptyCursor))
	if !success {
		return fmt.Errorf("setClassLongPtr failed")
	}
	return nil
}

func (w *windowImpl) SetTrayIcon(iconPath string) error {
	if w.trayGUID == nil {
		if err := w.createTrayItem(); err != nil {
			return err
		}
	}
	iconData := win32.NOTIFYICONDATA{}
	iconData.CbSize = uint32(unsafe.Sizeof(iconData))
	iconData.UFlags = win32.NIF_GUID | win32.NIF_MESSAGE
	iconData.HWnd = w.hwnd
	iconData.GUIDItem = *w.trayGUID
	iconData.UFlags = win32.NIF_GUID | win32.NIF_ICON
	var err error
	iconData.HIcon, err = win32.LoadImage(
		0,
		windows.StringToUTF16Ptr(iconPath),
		win32.IMAGE_ICON,
		0, 0,
		win32.LR_DEFAULTSIZE|win32.LR_LOADFROMFILE)
	if err != nil {
		return fmt.Errorf("failed to load icon: %w", err)
	}
	if !win32.Shell_NotifyIcon(win32.NIM_MODIFY, &iconData) {
		return fmt.Errorf("failed to create notification icon")
	}
	return nil
}

func (w *windowImpl) ShowNotification(title, msg string, icon bool) error {
	if w.trayGUID == nil {
		if err := w.createTrayItem(); err != nil {
			return err
		}
	}
	iconData := win32.NOTIFYICONDATA{}
	iconData.CbSize = uint32(unsafe.Sizeof(iconData))
	iconData.UFlags = win32.NIF_GUID | win32.NIF_INFO
	iconData.HWnd = w.hwnd
	iconData.GUIDItem = *w.trayGUID
	copy(iconData.SzInfoTitle[:], windows.StringToUTF16(title))
	copy(iconData.SzInfo[:], windows.StringToUTF16(msg))
	if icon {
		iconData.DwInfoFlags = win32.NIIF_USER | win32.NIIF_LARGE_ICON
	}

	if !win32.Shell_NotifyIcon(win32.NIM_MODIFY, &iconData) {
		return fmt.Errorf("failed to create notification icon")
	}
	return nil
}

func (w *windowImpl) createTrayItem() error {
	w.trayGUID = new(win32.GUID)
	*w.trayGUID = win32.MakeGUID(w.guid)
	iconData := win32.NOTIFYICONDATA{}
	iconData.CbSize = uint32(unsafe.Sizeof(iconData))
	iconData.UFlags = win32.NIF_GUID | win32.NIF_MESSAGE
	iconData.HWnd = w.hwnd
	iconData.GUIDItem = *w.trayGUID
	iconData.UCallbackMessage = win32.WM_APP + 1
	if !win32.Shell_NotifyIcon(win32.NIM_ADD, &iconData) {
		return fmt.Errorf("failed to create notification")
	}
	return nil
}

func (w *windowImpl) MoveWindow(x, y, wd, ht int32) error {
	return win32.MoveWindow(w.hwnd, x, y, wd, ht, true)
}

func drawWindow(dc win32.HDC, src2dst f64.Aff3, src interface{}, sr image.Rectangle, op draw.Op) (retErr error) {
	var dr image.Rectangle
	if src2dst[1] != 0 || src2dst[3] != 0 {
		// general drawing
		dr = sr.Sub(sr.Min)

		prevmode, err := _SetGraphicsMode(syscall.Handle(dc), _GM_ADVANCED)
		if err != nil {
			return err
		}
		defer func() {
			_, err := _SetGraphicsMode(syscall.Handle(dc), prevmode)
			if retErr == nil {
				retErr = err
			}
		}()

		x := _XFORM{
			eM11: +float32(src2dst[0]),
			eM12: -float32(src2dst[1]),
			eM21: -float32(src2dst[3]),
			eM22: +float32(src2dst[4]),
			eDx:  +float32(src2dst[2]),
			eDy:  +float32(src2dst[5]),
		}
		err = _SetWorldTransform(syscall.Handle(dc), &x)
		if err != nil {
			return err
		}
		defer func() {
			err := _ModifyWorldTransform(syscall.Handle(dc), nil, _MWT_IDENTITY)
			if retErr == nil {
				retErr = err
			}
		}()
	} else if src2dst[0] == 1 && src2dst[4] == 1 {
		// copy bitmap
		dr = sr.Add(image.Point{int(src2dst[2]), int(src2dst[5])})
	} else {
		// scale bitmap
		dstXMin := float64(sr.Min.X)*src2dst[0] + src2dst[2]
		dstXMax := float64(sr.Max.X)*src2dst[0] + src2dst[2]
		if dstXMin > dstXMax {
			// TODO: check if this (and below) works when src2dst[0] < 0.
			dstXMin, dstXMax = dstXMax, dstXMin
		}
		dstYMin := float64(sr.Min.Y)*src2dst[4] + src2dst[5]
		dstYMax := float64(sr.Max.Y)*src2dst[4] + src2dst[5]
		if dstYMin > dstYMax {
			// TODO: check if this (and below) works when src2dst[4] < 0.
			dstYMin, dstYMax = dstYMax, dstYMin
		}
		dr = image.Rectangle{
			image.Point{int(math.Floor(dstXMin)), int(math.Floor(dstYMin))},
			image.Point{int(math.Ceil(dstXMax)), int(math.Ceil(dstYMax))},
		}
	}
	switch s := src.(type) {
	case syscall.Handle:
		return copyBitmapToDC(dc, dr, s, sr, op)
	case color.Color:
		return fill(dc, dr, s, op)
	}
	return fmt.Errorf("unsupported type %T", src)
}

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Copy(w, dp, src, sr, op)
}

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, op)
}

func (w *windowImpl) Publish() screen.PublishResult {
	// TODO
	return screen.PublishResult{}
}

func init() {
	send := func(hwnd win32.HWND, e interface{}) {
		windowLock.RLock()
		w := allWindows[hwnd]
		windowLock.RUnlock()

		w.Send(e)
	}
	win32.MouseEvent = func(hwnd win32.HWND, e mouse.Event) { send(hwnd, e) }
	win32.PaintEvent = func(hwnd win32.HWND, e paint.Event) { send(hwnd, e) }
	win32.KeyEvent = func(hwnd win32.HWND, e key.Event) { send(hwnd, e) }
	win32.LifecycleEvent = lifecycleEvent
	win32.SizeEvent = sizeEvent
}

func lifecycleEvent(hwnd win32.HWND, to lifecycle.Stage) {
	windowLock.RLock()
	w := allWindows[hwnd]
	windowLock.RUnlock()

	if w.lifecycleStage == to {
		return
	}
	w.Send(lifecycle.Event{
		From: w.lifecycleStage,
		To:   to,
	})
	w.lifecycleStage = to
	if w.lifecycleStage == lifecycle.StageDead {
		w.Release()
	}
}

func sizeEvent(hwnd win32.HWND, e size.Event) {
	windowLock.RLock()
	w := allWindows[win32.HWND(hwnd)]
	windowLock.RUnlock()

	w.Send(e)

	if e != w.sz {
		w.sz = e
		w.Send(paint.Event{})
	}
}

// cmd is used to carry parameters between user code
// and Windows message pump thread.
type cmd struct {
	id  int
	err error

	src2dst f64.Aff3
	sr      image.Rectangle
	dp      image.Point
	dr      image.Rectangle
	color   color.Color
	op      draw.Op
	texture syscall.Handle
	buffer  *bufferImpl
}

const (
	cmdDraw = iota
	cmdFill
	cmdUpload
	cmdDrawUniform
)

// msgCmd is the stored value for our handleCmd function for syscalls.
var msgCmd = win32.AddWindowMsg(handleCmd)

func (w *windowImpl) execCmd(c *cmd) {
	win32.SendMessage(win32.HWND(w.hwnd), msgCmd, 0, uintptr(unsafe.Pointer(c)))
	if c.err != nil {
		panic(fmt.Sprintf("execCmd faild for cmd.id=%d: %v", c.id, c.err)) // TODO handle errors
	}
}

func handleCmd(hwnd win32.HWND, uMsg uint32, wParam, lParam uintptr) {
	c := (*cmd)(unsafe.Pointer(lParam))

	dc, err := win32.GetDC(hwnd)
	if err != nil {
		c.err = err
		return
	}
	defer win32.ReleaseDC(hwnd, dc)

	switch c.id {
	case cmdDraw:
		c.err = drawWindow(dc, c.src2dst, c.texture, c.sr, c.op)
	case cmdDrawUniform:
		c.err = drawWindow(dc, c.src2dst, c.color, c.sr, c.op)
	case cmdFill:
		c.err = fill(dc, c.dr, c.color, c.op)
	case cmdUpload:
		// TODO: adjust if dp is outside dst bounds, or sr is outside buffer bounds.
		dr := c.sr.Add(c.dp.Sub(c.sr.Min))
		c.err = copyBitmapToDC(dc, dr, c.buffer.hbitmap, c.sr, draw.Src)
	default:
		c.err = fmt.Errorf("unknown command id=%d", c.id)
	}
}

func (w *windowImpl) GetCursorPosition() (x, y float64) {
	xint, yint, _ := win32.GetCursorPos()
	return float64(xint), float64(yint)
}
