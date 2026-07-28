// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
)

type _WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HWND
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type _WINDOWPOS struct {
	HWND            syscall.Handle
	HWNDInsertAfter syscall.Handle
	X               int32
	Y               int32
	Cx              int32
	Cy              int32
	Flags           uint32
}

const (
	_WM_SETFOCUS         = 7
	_WM_KILLFOCUS        = 8
	_WM_PAINT            = 15
	_WM_CLOSE            = 16
	_WM_WINDOWPOSCHANGED = 71
	_WM_KEYDOWN          = 256
	_WM_KEYUP            = 257
	//_WM_SYSKEYDOWN       = 260
	//_WM_SYSKEYUP         = 261
	_WM_MOUSEMOVE       = 512
	_WM_MOUSEWHEEL      = 522
	_WM_LBUTTONDOWN     = 513
	_WM_LBUTTONUP       = 514
	_WM_RBUTTONDOWN     = 516
	_WM_RBUTTONUP       = 517
	_WM_MBUTTONDOWN     = 519
	_WM_MBUTTONUP       = 520
	_WM_USER            = 0x0400
	_WM_INPUTLANGCHANGE = 0x0051
)

// WM_SYSCOMMAND

const (
	_VK_SHIFT   = 16
	_VK_CONTROL = 17
	_VK_MENU    = 18
	_VK_LWIN    = 0x5B
	_VK_RWIN    = 0x5C
)

const (
	IDI_APPLICATION = 32512
	IDC_ARROW       = 32512
)

const (
	_CW_USEDEFAULT = 0x80000000 - 0x100000000

	_SW_SHOWNORMAL    = 1
	_SW_SHOWMINIMIZED = 2
	_SW_SHOWMAXIMIZED = 3
	_SW_MINIMIZE      = 6
	_SW_SHOWDEFAULT   = 10

	_SWP_NOSIZE = 0x0001
)

const (
	_WHEEL_DELTA = 120
)

func _GET_X_LPARAM(lp uintptr) int32 {
	return int32(_LOWORD(lp))
}

func _GET_Y_LPARAM(lp uintptr) int32 {
	return int32(_HIWORD(lp))
}

func _GET_WHEEL_DELTA_WPARAM(lp uintptr) int16 {
	return int16(_HIWORD(lp))
}

func _LOWORD(l uintptr) uint16 {
	return uint16(uint32(l))
}

func _HIWORD(l uintptr) uint16 {
	return uint16(uint32(l >> 16))
}

// notes to self
// UINT = uint32
// callbacks = uintptr
// strings = *uint16

//sys	GetDC(hwnd syscall.Handle) (dc syscall.Handle, err error) = user32.GetDC
//sys	ReleaseDC(hwnd syscall.Handle, dc syscall.Handle) (err error) = user32.ReleaseDC

//sys	_CreateWindowEx(exstyle uint32, className *uint16, windowText *uint16, style uint32, x int32, y int32, width int32, height int32, parent syscall.Handle, menu syscall.Handle, hInstance syscall.Handle, lpParam uintptr) (hwnd syscall.Handle, err error) = user32.CreateWindowExW
//sys	_DefWindowProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.DefWindowProcW
//sys	_GetClientRect(hwnd syscall.Handle, rect *_RECT) (err error) = user32.GetClientRect
//sys	_GetWindowRect(hwnd syscall.Handle, rect *_RECT) (err error) = user32.GetWindowRect
//sys   _GetKeyboardLayout(threadID uint32) (locale syscall.Handle) = user32.GetKeyboardLayout
//sys   _GetKeyboardState(lpKeyState *byte) (err error) = user32.GetKeyboardState
//sys	_GetKeyState(virtkey int32) (keystatus int16) = user32.GetKeyState
//sys	_GetMessage(msg *_MSG, hwnd syscall.Handle, msgfiltermin uint32, msgfiltermax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	_LoadCursor(hInstance syscall.Handle, cursorName uintptr) (cursor syscall.Handle, err error) = user32.LoadCursorW
//sys	_LoadIcon(hInstance syscall.Handle, iconName uintptr) (icon syscall.Handle, err error) = user32.LoadIconW
//sys	_MoveWindow(hwnd syscall.Handle, x int32, y int32, w int32, h int32, repaint bool) (err error) = user32.MoveWindow
//sys   _PostQuitMessage(exitCode int32) = user32.PostQuitMessage
//sys	_RegisterClass(wc *_WNDCLASS) (atom uint16, err error) = user32.RegisterClassW
//sys   _ToUnicodeEx(wVirtKey uint32, wScanCode uint32, lpKeyState *byte, pwszBuff *uint16, cchBuff int32, wFlags uint32, dwhkl syscall.Handle) (ret int32) = user32.ToUnicodeEx
