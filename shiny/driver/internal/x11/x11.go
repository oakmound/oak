package x11

import (
	"fmt"

	"github.com/BurntSushi/xgbutil/ewmh"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/motif"
)

func MoveWindow(xc *xgb.Conn, xw xproto.Window, x, y, width, height int32) (int32, int32, int32, int32) {
	vals := []uint32{}

	flags := xproto.ConfigWindowHeight |
		xproto.ConfigWindowWidth |
		xproto.ConfigWindowX |
		xproto.ConfigWindowY

	vals = append(vals, uint32(x))
	vals = append(vals, uint32(y))

	if int16(width) <= 0 {
		width = 1
	}
	vals = append(vals, uint32(width))

	if int16(height) <= 0 {
		height = 1
	}
	vals = append(vals, uint32(height))

	cook := xproto.ConfigureWindowChecked(xc, xw, uint16(flags), vals)
	if err := cook.Check(); err != nil {
		fmt.Println("X11 configure window failed: ", err)
	}
	return x, y, width, height
}

func ToggleFullScreen(xutil *xgbutil.XUtil, win xproto.Window) error {
	return ewmh.WmStateReq(xutil, win, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")
}

func SetFullScreen(xutil *xgbutil.XUtil, win xproto.Window, fullscreen bool) error {
	action := ewmh.StateRemove
	if fullscreen {
		action = ewmh.StateAdd
	}
	return ewmh.WmStateReq(xutil, win, action, "_NET_WM_STATE_FULLSCREEN")
}

func ToggleTopMost(xutil *xgbutil.XUtil, win xproto.Window) error {
	return ewmh.WmStateReq(xutil, win, ewmh.StateToggle, "_NET_WM_STATE_ABOVE")
}

func SetTopMost(xutil *xgbutil.XUtil, win xproto.Window, topMost bool) error {
	action := ewmh.StateRemove
	if topMost {
		action = ewmh.StateAdd
	}
	return ewmh.WmStateReq(xutil, win, action, "_NET_WM_STATE_ABOVE")
}

func SetBorderless(xutil *xgbutil.XUtil, win xproto.Window, borderless bool) error {
	hints := &motif.Hints{
		Flags:      motif.HintDecorations,
		Decoration: motif.DecorationNone, 
	}
	if !borderless {
		hints.Decoration = motif.DecorationAll
	}
	return motif.WmHintsSet(xutil, win, hints)
}
