// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package windriver

import (
	"image"
	"sync"
	"syscall"

	"github.com/oakmound/oak/v3/shiny/driver/internal/swizzle"
	"github.com/oakmound/oak/v3/shiny/driver/internal/win32"
)

type bufferImpl struct {
	hbitmap syscall.Handle
	buf     []byte
	rgba    image.RGBA
	size    image.Point

	mu        sync.Mutex
	nUpload   uint32
	released  bool
	cleanedUp bool
}

func (b *bufferImpl) Size() image.Point       { return b.size }
func (b *bufferImpl) Bounds() image.Rectangle { return image.Rectangle{Max: b.size} }
func (b *bufferImpl) RGBA() *image.RGBA       { return &b.rgba }

func (b *bufferImpl) preUpload() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.nUpload == 0 {
		swizzle.BGRA(b.buf)
	}
	b.nUpload++
}

func (b *bufferImpl) postUpload() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nUpload--
	if b.nUpload != 0 {
		return
	}

	if b.released {
		go b.cleanUp()
	}
}

func (b *bufferImpl) Release() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.released && b.nUpload == 0 {
		go b.cleanUp()
	}
	b.released = true
}

func (b *bufferImpl) cleanUp() {
	b.mu.Lock()
	b.cleanedUp = true
	b.mu.Unlock()

	b.rgba.Pix = nil
	_DeleteObject(b.hbitmap)
}

func (b *bufferImpl) blitToDC(dc syscall.Handle, dp image.Point, sr image.Rectangle) error {
	b.preUpload()
	defer b.postUpload()

	dr := sr.Add(dp.Sub(sr.Min))
	return copyBitmapToDC(win32.HDC(dc), dr, b.hbitmap, sr)
}
