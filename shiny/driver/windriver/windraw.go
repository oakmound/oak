// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package windriver

import (
	"image"
	"syscall"
	"unsafe"

	"github.com/oakmound/oak/v3/shiny/driver/internal/win32"
)

func mkbitmap(size image.Point) (syscall.Handle, *byte, error) {
	bi := _BITMAPINFO{
		Header: _BITMAPINFOHEADER{
			Size:        uint32(unsafe.Sizeof(_BITMAPINFOHEADER{})),
			Width:       int32(size.X),
			Height:      -int32(size.Y), // negative height to force top-down drawing
			Planes:      1,
			BitCount:    32,
			Compression: _BI_RGB,
			SizeImage:   uint32(size.X * size.Y * 4),
		},
	}

	var ppvBits *byte
	bitmap, err := _CreateDIBSection(0, &bi, _DIB_RGB_COLORS, &ppvBits, 0, 0)
	if err != nil {
		return 0, nil, err
	}
	return bitmap, ppvBits, nil
}

var blendOverFunc = _BLENDFUNCTION{
	BlendOp:             _AC_SRC_OVER,
	BlendFlags:          0,
	SourceConstantAlpha: 255,           // only use per-pixel alphas
	AlphaFormat:         _AC_SRC_ALPHA, // premultiplied
}

func copyBitmapToDC(dc win32.HDC, dr image.Rectangle, src syscall.Handle, sr image.Rectangle) (retErr error) {
	memdc, err := _CreateCompatibleDC(syscall.Handle(dc))
	if err != nil {
		return err
	}
	defer _DeleteDC(memdc)

	prev, err := _SelectObject(memdc, src)
	if err != nil {
		return err
	}
	defer func() {
		_, err2 := _SelectObject(memdc, prev)
		if retErr == nil {
			retErr = err2
		}
	}()

	return _StretchBlt(syscall.Handle(dc), int32(dr.Min.X), int32(dr.Min.Y), int32(dr.Dx()), int32(dr.Dy()),
		memdc, int32(sr.Min.X), int32(sr.Min.Y), int32(sr.Dx()), int32(sr.Dy()), _SRCCOPY)
}
