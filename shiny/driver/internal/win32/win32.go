// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

// Package win32 implements a partial shiny screen driver using the Win32 API.
// It provides window, lifecycle, key, and mouse management, but no drawing.
// That is left to windriver (using GDI) or gldriver (using DirectX via ANGLE).
package win32

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/oakmound/oak/v4/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

// screenHWND is the handle to the "Screen window".
// The Screen window encapsulates all screen.Screen operations
// in an actual Windows window so they all run on the main thread.
// Since any messages sent to a window will be executed on the
// main thread, we can safely use the messages below.
var screenHWND HWND

const (
	msgCreateWindow = _WM_USER + iota
	msgShow
	msgQuit
	msgLast // WM_USER value https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-user
)

var msgCallbacks = func() *uint32 {
	u := new(uint32)
	*u = msgLast + 1
	return u
}()

// userWM is used to generate private (WM_USER and above) window message IDs
// for use by screenWindowWndProc and windowWndProc.
type userWM struct {
	sync.Mutex
	id uint32
}

// next id for the given userWM (which is a construct purely used to generate unique ids).
func (m *userWM) next() uint32 {
	m.Lock()
	if m.id == 0 {
		m.id = msgLast
	}
	r := m.id
	m.id++
	m.Unlock()
	return r
}

// currentUserM gives a quick handle to globally mess with userWM.
var currentUserWM userWM

func newWindow(opts screen.WindowGenerator, class string) (HWND, error) {
	wcname, err := syscall.UTF16PtrFromString(class)
	if err != nil {
		return 0, err
	}
	title, err := syscall.UTF16PtrFromString(opts.Title)
	if err != nil {
		return 0, err
	}
	style, exStyle := WindowsStyle(opts)
	// This should be a feature, putting windows on the top layer
	if opts.TopMost {
		exStyle = exStyle | WS_EX_TOPMOST
	}
	hwnd, err := CreateWindowEx(exStyle,
		wcname, title,
		style,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		0, 0, hThisInstance, 0)
	if err != nil {
		return 0, err
	}

	// This is interesting and we'll use it eventually
	//SetWindowLongPtr(hwnd, GWL_STYLE, 0)
	// TODO(andlabs): use proper nCmdShow
	// TODO(andlabs): call UpdateWindow()

	return hwnd, nil
}

// WindowsStyle converts a screen.BorderStyle into a style and
// exStyle for a Windows window
func WindowsStyle(gen screen.WindowGenerator) (uint32, uint32) {
	return WS_OVERLAPPEDWINDOW, 0
}

// ResizeClientRect makes hwnd client rectangle opts.Width by opts.Height in size.
func ResizeClientRect(hwnd HWND, opts screen.WindowGenerator) error {
	if opts.Width <= 0 || opts.Height <= 0 {
		return errors.New("Invalid inputs to ResizeClientRect")
	}
	cr, err := GetClientRect(hwnd)
	if err != nil {
		return err
	}
	wr, err := GetWindowRect(hwnd)
	if err != nil {
		return err
	}
	w := (wr.Right - wr.Left) - (cr.Right - int32(opts.Width))
	h := (wr.Bottom - wr.Top) - (cr.Bottom - int32(opts.Height))
	x := wr.Left
	if opts.X != 0 {
		x = int32(opts.X)
	}
	y := wr.Top
	if opts.Y != 0 {
		y = int32(opts.Y)
	}
	return MoveWindow(hwnd, x, y, w, h, false)
}

// Show shows a newly created window.
// It sends the appropriate lifecycle events, makes the window appear
// on the screen, and sends an initial size event.
//
// This is a separate step from NewWindow to give the driver a chance
// to setup its internal state for a window before events start being
// delivered.
func Show(hwnd HWND) {
	SendMessage(hwnd, msgShow, 0, 0)
}

// Release sends the close message to the specified window.
// https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-close
func Release(hwnd HWND) {
	SendMessage(hwnd, WM_CLOSE, 0, 0)
}

// sendFocus change to the specified window.
// There is some value here but the panic is not safe for consumption.
// Consider: wrapper func or rewrite.
func sendFocus(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	switch uMsg {
	case _WM_SETFOCUS:
		LifecycleEvent(hwnd, lifecycle.StageFocused)
	case _WM_KILLFOCUS:
		LifecycleEvent(hwnd, lifecycle.StageVisible)
	default:
		panic(fmt.Sprintf("unexpected focus message: %d", uMsg))
	}
	lResult, _ = DefWindowProc(hwnd, uMsg, wParam, lParam)
	return lResult
}

func sendShow(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	LifecycleEvent(hwnd, lifecycle.StageVisible)
	ShowWindow(hwnd, _SW_SHOWDEFAULT)
	sendSize(hwnd)
	return 0
}

func sendSizeEvent(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	wp := (*_WINDOWPOS)(unsafe.Pointer(lParam))
	if wp.Flags&_SWP_NOSIZE != 0 {
		return 0
	}
	sendSize(hwnd)
	return 0
}

func sendSize(hwnd HWND) {
	r, err := GetClientRect(hwnd)
	if err != nil {
		panic(err) // TODO(andlabs)
	}

	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)

	// TODO(andlabs): don't assume that PixelsPerPt == 1
	SizeEvent(hwnd, size.Event{
		WidthPx:     width,
		HeightPx:    height,
		WidthPt:     geom.Pt(width),
		HeightPt:    geom.Pt(height),
		PixelsPerPt: 1,
	})
}

func sendClose(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	LifecycleEvent(hwnd, lifecycle.StageDead)
	ptr, _ := DefWindowProc(hwnd, uMsg, wParam, lParam)
	return ptr
}

func sendMouseEvent(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	e := mouse.Event{
		X:         float32(_GET_X_LPARAM(lParam)),
		Y:         float32(_GET_Y_LPARAM(lParam)),
		Modifiers: keyModifiers(),
	}

	switch uMsg {
	case _WM_MOUSEMOVE:
		e.Direction = mouse.DirNone
	case _WM_LBUTTONDOWN, _WM_MBUTTONDOWN, _WM_RBUTTONDOWN:
		e.Direction = mouse.DirPress
	case _WM_LBUTTONUP, _WM_MBUTTONUP, _WM_RBUTTONUP:
		e.Direction = mouse.DirRelease
	case _WM_MOUSEWHEEL:
		// TODO: On a trackpad, a scroll can be a drawn-out affair with a
		// distinct beginning and end. Should the intermediate events be
		// DirNone?
		e.Direction = mouse.DirStep

		x, y, _ := ScreenToClient(hwnd, int(e.X), int(e.Y))
		e.X = float32(x)
		e.Y = float32(y)
	default:
		panic("sendMouseEvent() called on non-mouse message")
	}

	switch uMsg {
	case _WM_MOUSEMOVE:
		// No-op.
	case _WM_LBUTTONDOWN, _WM_LBUTTONUP:
		e.Button = mouse.ButtonLeft
	case _WM_MBUTTONDOWN, _WM_MBUTTONUP:
		e.Button = mouse.ButtonMiddle
	case _WM_RBUTTONDOWN, _WM_RBUTTONUP:
		e.Button = mouse.ButtonRight
	case _WM_MOUSEWHEEL:
		// TODO: handle horizontal scrolling
		delta := _GET_WHEEL_DELTA_WPARAM(wParam) / _WHEEL_DELTA
		switch {
		case delta > 0:
			e.Button = mouse.ButtonWheelUp
		case delta < 0:
			e.Button = mouse.ButtonWheelDown
			delta = -delta
		default:
			return
		}
		for delta > 0 {
			MouseEvent(hwnd, e)
			delta--
		}
		return
	}

	MouseEvent(hwnd, e)

	return 0
}

// Precondition: this is called in immediate response to the message that triggered the event (so not after w.Send).
func keyModifiers() (m key.Modifiers) {
	down := func(x int32) bool {
		// GetKeyState gets the key state at the time of the message, so this is what we want.
		return _GetKeyState(x)&0x80 != 0
	}

	if down(_VK_CONTROL) {
		m |= key.ModControl
	}
	if down(_VK_MENU) {
		m |= key.ModAlt
	}
	if down(_VK_SHIFT) {
		m |= key.ModShift
	}
	if down(_VK_LWIN) || down(_VK_RWIN) {
		m |= key.ModMeta
	}
	return m
}

var (
	MouseEvent     func(hwnd HWND, e mouse.Event)
	PaintEvent     func(hwnd HWND, e paint.Event)
	SizeEvent      func(hwnd HWND, e size.Event)
	KeyEvent       func(hwnd HWND, e key.Event)
	LifecycleEvent func(hwnd HWND, e lifecycle.Stage)

	// TODO: use the golang.org/x/exp/shiny/driver/internal/lifecycler package
	// instead of or together with the LifecycleEvent callback?
)

func sendPaint(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr) {
	PaintEvent(hwnd, paint.Event{})
	lResult, _ = DefWindowProc(hwnd, uMsg, wParam, lParam)
	return lResult
}

var screenMsgs = map[uint32]func(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr){}

func AddScreenMsg(fn func(hwnd HWND, uMsg uint32, wParam, lParam uintptr)) uint32 {
	uMsg := currentUserWM.next()
	screenMsgs[uMsg] = func(hwnd HWND, uMsg uint32, wParam, lParam uintptr) uintptr {
		fn(hwnd, uMsg, wParam, lParam)
		return 0
	}
	return uMsg
}

func screenWindowWndProc(hwnd HWND, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	switch uMsg {
	case msgCreateWindow:
		p := (*newWindowParams)(unsafe.Pointer(lParam))
		p.w, p.err = newWindow(p.opts, p.class)
	case msgQuit:
		_PostQuitMessage(0)
	}
	callbacksLock.RLock()
	if callback, ok := callbacks[uMsg]; ok {
		go func() {
			callback()
			SendScreenMessage(hwnd, msgQuit, 0, 0)
		}()
	}
	callbacksLock.RUnlock()
	fn := screenMsgs[uMsg]
	if fn != nil {
		return fn(hwnd, uMsg, wParam, lParam)
	}
	lResult, _ = DefWindowProc(hwnd, uMsg, wParam, lParam)
	return lResult
}

//go:uintptrescapes

// SendScreenMessage is a perhaps poorly named wrapper for SendMessage where we know that lParam has a pointer in its call.
// Ends up calling https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagew so thats cool.
func SendScreenMessage(screen HWND, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	return SendMessage(screen, uMsg, wParam, lParam)
}

var windowMsgs = map[uint32]func(hwnd HWND, uMsg uint32, wParam, lParam uintptr) (lResult uintptr){
	_WM_SETFOCUS:         sendFocus,
	_WM_KILLFOCUS:        sendFocus,
	_WM_PAINT:            sendPaint,
	msgShow:              sendShow,
	_WM_WINDOWPOSCHANGED: sendSizeEvent,
	_WM_CLOSE:            sendClose,

	_WM_LBUTTONDOWN: sendMouseEvent,
	_WM_LBUTTONUP:   sendMouseEvent,
	_WM_MBUTTONDOWN: sendMouseEvent,
	_WM_MBUTTONUP:   sendMouseEvent,
	_WM_RBUTTONDOWN: sendMouseEvent,
	_WM_RBUTTONUP:   sendMouseEvent,
	_WM_MOUSEMOVE:   sendMouseEvent,
	_WM_MOUSEWHEEL:  sendMouseEvent,

	_WM_KEYDOWN:         sendKeyEvent,
	_WM_KEYUP:           sendKeyEvent,
	_WM_INPUTLANGCHANGE: updateKeyboardLayout,
	// TODO case _WM_SYSKEYDOWN, _WM_SYSKEYUP:
}

// AddWindowMsg stores a given window manipulator so it can be accessed via syscalls.
// Stores a reference to the reference argument for the the given id.
func AddWindowMsg(fn func(hwnd HWND, uMsg uint32, wParam, lParam uintptr)) uint32 {
	uMsg := currentUserWM.next()
	windowMsgs[uMsg] = func(hwnd HWND, uMsg uint32, wParam, lParam uintptr) uintptr {
		fn(hwnd, uMsg, wParam, lParam)
		return 0
	}
	return uMsg
}

// src: https://wiki.winehq.org/List_Of_Windows_Messages
// var unusedMessages = map[uint32]string{
// 	2:   "DESTROY",
// 	6:   "ACTIVATE",
// 	28:  "ACTIVATE_APP",
// 	32:  "SETCURSOR",
// 	70:  "WINDOWPOSCHANGING",
// 	130: "NCDESTROY",
// 	132: "NCHITTEST",
// 	134: "NCACTIVATE",
// 	144: "", // we get this, but its not documented in the source list
// 	160: "NCMOUSEMOVE",
// 	161: "NCLBUTTONDOWN",
// 	273: "COMMAND",
// 	274: "SYSCOMMAND",
// 	533: "CAPTURECHANGED",
// 	641: "IME_SETCONTEXT",
// 	642: "IME_NOTIFY",
// 	674: "NCMOUSELEAVE",
// }

func windowWndProc(hwnd HWND, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	fn := windowMsgs[uMsg]
	if fn != nil {
		return fn(hwnd, uMsg, wParam, lParam)
	}
	//fmt.Printf("unused message %d, 0x%x, %v\n", uMsg, uMsg, unusedMessages[uMsg])
	lResult, _ = DefWindowProc(hwnd, uMsg, wParam, lParam)
	return lResult
}

type newWindowParams struct {
	opts  screen.WindowGenerator
	w     HWND
	class string
	err   error
}

var nextWindow = new(int32)

// NewWindow attempts to register a screen on the given handle.
func NewWindow(screenHWND HWND, opts screen.WindowGenerator) (HWND, error) {
	var p newWindowParams
	p.opts = opts
	p.class = "shiny_Window" + strconv.Itoa(int(atomic.AddInt32(nextWindow, 1)))
	err := initWindowClass(p.class)
	if err != nil {
		return 0, fmt.Errorf("failed to register window: %w", err)
	}

	SendScreenMessage(screenHWND, msgCreateWindow, 0, uintptr(unsafe.Pointer(&p)))
	return p.w, p.err
}

func initWindowClass(class string) (err error) {
	wcname, err := syscall.UTF16PtrFromString(class)
	if err != nil {
		return err
	}
	_, err = RegisterClass(&_WNDCLASS{
		LpszClassName: wcname,
		LpfnWndProc:   syscall.NewCallback(windowWndProc),
		HIcon:         hDefaultIcon,
		HCursor:       hDefaultCursor,
		HInstance:     hThisInstance,
		HbrBackground: COLOR_BTNSHADOW,
	})
	return err
}

var nextScreenWindow = new(int32)

func initScreenWindow() (HWND, error) {
	screenWindowClass := "shiny_ScreenWindow" + strconv.Itoa(int(atomic.AddInt32(nextScreenWindow, 1)))
	swc, err := syscall.UTF16PtrFromString(screenWindowClass)
	if err != nil {
		return 0, err
	}
	emptyString, err := syscall.UTF16PtrFromString("")
	if err != nil {
		return 0, err
	}
	wc := _WNDCLASS{
		LpszClassName: swc,
		LpfnWndProc:   syscall.NewCallback(screenWindowWndProc),
		HIcon:         hDefaultIcon,
		HCursor:       hDefaultCursor,
		HInstance:     hThisInstance,
		HbrBackground: HWND(COLOR_BTNSHADOW),
	}
	_, err = RegisterClass(&wc)
	if err != nil {
		return 0, err
	}
	screenHWND, err = CreateWindowEx(0,
		swc, emptyString,
		windowStyle,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		HWND_MESSAGE, 0, hThisInstance, 0)
	if err != nil {
		return 0, err
	}
	return screenHWND, nil
}

var (
	windowStyle uint32 = WS_OVERLAPPEDWINDOW
)

var (
	hDefaultIcon   HICON
	hDefaultCursor HCURSOR
	hThisInstance  HINSTANCE
)

// initCommon attempts to set up some standard icons.
// TODO: Consider running this only once if successful.
func initCommon() (err error) {
	hDefaultIcon, err = LoadIcon(0, IDI_APPLICATION)
	if err != nil {
		return err
	}
	hDefaultCursor, err = LoadCursor(0, IDC_ARROW)
	if err != nil {
		return err
	}
	// TODO(andlabs) hThisInstance
	return nil
}

// Todo: this (and other globals) forces this package to only be able to run one window.
// Can we change this?
var (
	callbacksLock sync.RWMutex
	callbacks     = map[uint32]func(){}
)

// NewScreen sets up common infos and then attempts to create a new window.
func NewScreen() (HWND, error) {
	if err := initCommon(); err != nil {
		return 0, fmt.Errorf("init common failed: %w", err)
	}

	screenHWND, err := initScreenWindow()
	if err != nil {
		return 0, fmt.Errorf("init screen window failed: %w", err)
	}

	return screenHWND, nil
}

func Main(screenHWND HWND, f func()) error {
	keyboardLayout = _GetKeyboardLayout(0)
	defer func() {
		// TODO(andlabs): log an error if this fails?
		DestroyWindow(screenHWND)
		// TODO(andlabs): unregister window class
	}()
	// It does not matter which OS thread we are on.
	// All that matters is that we confine all UI operations
	// to the thread that created the respective window.
	runtime.LockOSThread()

	cb := atomic.AddUint32(msgCallbacks, 1)
	// Prime the pump.
	callbacksLock.Lock()
	callbacks[cb] = f
	callbacksLock.Unlock()
	PostMessage(screenHWND, cb, 0, 0)

	// Main message pump.
	var m MSG
	for {
		done, err := GetMessage(&m, 0, 0, 0)
		if err != nil {
			return fmt.Errorf("win32 GetMessage failed: %v %d", err, uintptr(err.(syscall.Errno)))
		}
		if done == 0 { // WM_QUIT
			break
		}
		TranslateMessage(&m)
		DispatchMessage(&m)
	}

	return nil
}
