// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

// Package mtldriver provides a Metal driver for accessing a screen.
//
// At this time, the Metal API is used only to present the final pixels
// to the screen. All rendering is performed on the CPU via the image/draw
// algorithms. Future work is to use mtl.Buffer, mtl.Texture, etc., and
// do more of the rendering work on the GPU.
package mtldriver

import (
	"image"
	"runtime"
	"unsafe"

	"dmitri.shuralyov.com/gpu/mtl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/oakmound/oak/v3/shiny/driver/internal/errscreen"
	"github.com/oakmound/oak/v3/shiny/driver/mtldriver/internal/appkit"
	"github.com/oakmound/oak/v3/shiny/driver/mtldriver/internal/coreanim"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func init() {
	runtime.LockOSThread()
}

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the Screen, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(screen.Screen)) {
	if err := main(f); err != nil {
		f(errscreen.Stub(err))
	}
}

func main(f func(screen.Screen)) error {
	device, err := mtl.CreateSystemDefaultDevice()
	if err != nil {
		return err
	}
	err = glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	{
		// TODO(dmitshur): Delete this when https://github.com/go-gl/glfw/issues/272 is resolved.
		// Post an empty event from the main thread before it can happen in a non-main thread,
		// to work around https://github.com/glfw/glfw/issues/1649.
		glfw.PostEmptyEvent()
	}
	var (
		done            = make(chan struct{})
		newWindowCh     = make(chan newWindowReq, 1)
		releaseWindowCh = make(chan releaseWindowReq, 1)
		moveWindowCh    = make(chan moveWindowReq, 1)
	)
	go func() {
		f(&screenImpl{
			newWindowCh: newWindowCh,
		})
		close(done)
		glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on done.
	}()
	for {
		select {
		case <-done:
			return nil
		case req := <-newWindowCh:
			w, err := newWindow(device, releaseWindowCh, moveWindowCh, req.opts)
			req.respCh <- newWindowResp{w, err}
		case req := <-releaseWindowCh:
			req.window.Destroy()
			req.respCh <- struct{}{}
		case req := <-moveWindowCh:
			req.window.SetPos(int(req.x), int(req.y))
			req.window.SetSize(int(req.width), int(req.height))
			req.respCh <- struct{}{}
		default:
			glfw.WaitEvents()
		}
	}
}

type newWindowReq struct {
	opts   screen.WindowGenerator
	respCh chan newWindowResp
}

type newWindowResp struct {
	w   screen.Window
	err error
}
type moveWindowReq struct {
	window              *glfw.Window
	x, y, width, height int
	respCh              chan struct{}
}

type releaseWindowReq struct {
	window *glfw.Window
	respCh chan struct{}
}

// newWindow creates a new GLFW window.
// It must be called on the main thread.
func newWindow(device mtl.Device, releaseWindowCh chan releaseWindowReq, moveWindowCh chan moveWindowReq, opts screen.WindowGenerator) (screen.Window, error) {
	width, height := optsSize(opts)
	window, err := glfw.CreateWindow(width, height, opts.Title, nil, nil)
	if err != nil {
		return nil, err
	}

	ml := coreanim.MakeMetalLayer()
	ml.SetDevice(device)
	ml.SetPixelFormat(mtl.PixelFormatBGRA8UNorm)
	ml.SetMaximumDrawableCount(3)
	ml.SetDisplaySyncEnabled(true)
	cv := appkit.NewWindow(unsafe.Pointer(window.GetCocoaWindow())).ContentView()
	cv.SetLayer(ml)
	cv.SetWantsLayer(true)
	if opts.Borderless {
		window.SetAttrib(glfw.Decorated, 0)
	}

	w := &windowImpl{
		device:          device,
		window:          window,
		releaseWindowCh: releaseWindowCh,
		moveWindowCh:    moveWindowCh,
		ml:              ml,
		cq:              device.MakeCommandQueue(),
		rgba:            image.NewRGBA(image.Rectangle{Max: image.Point{X: opts.Width, Y: opts.Height}}),
		texture: device.MakeTexture(mtl.TextureDescriptor{
			PixelFormat: mtl.PixelFormatRGBA8UNorm,
			Width:       opts.Width,
			Height:      opts.Height,
			StorageMode: mtl.StorageModeManaged,
		}),
	}

	// Set callbacks.
	framebufferSizeCallback := func(_ *glfw.Window, width, height int) {
		w.Send(size.Event{
			WidthPx:  width,
			HeightPx: height,
			// TODO(dmitshur): ppp,
		})
		w.Send(paint.Event{External: true})
	}
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		const scale = 2 // TODO(dmitshur): compute dynamically
		w.Send(mouse.Event{X: float32(x * scale), Y: float32(y * scale)})
	})
	window.SetScrollCallback(func(_ *glfw.Window, xoff float64, yoff float64) {
		// TODO horizontal scrolling
		var btn mouse.Button
		if yoff < 0 {
			btn = mouse.ButtonWheelDown
		} else {
			btn = mouse.ButtonWheelUp
		}
		w.Send(mouse.Event{
			Button:    btn,
			Direction: mouse.DirNone,
		})
	})
	window.SetMouseButtonCallback(func(_ *glfw.Window, b glfw.MouseButton, a glfw.Action, mods glfw.ModifierKey) {
		btn := glfwMouseButton(b)
		if btn == mouse.ButtonNone {
			return
		}
		const scale = 2 // TODO(dmitshur): compute dynamically
		x, y := window.GetCursorPos()
		w.Send(mouse.Event{
			X: float32(x * scale), Y: float32(y * scale),
			Button:    btn,
			Direction: glfwMouseDirection(a),
			Modifiers: glfwKeyMods(mods),
		})
	})
	// TODO: can we combine the following two callbacks into a single event? Signs point to no.
	window.SetKeyCallback(func(_ *glfw.Window, k glfw.Key, _ int, a glfw.Action, mods glfw.ModifierKey) {
		code := glfwKeyCode(k)
		if code == key.CodeUnknown {
			return
		}
		ev := key.Event{
			Code:      code,
			Direction: glfwKeyDirection(a),
			Modifiers: glfwKeyMods(mods),
		}
		w.Send(ev)
	})
	// TODO: some characters will repeat when held down, but not all of them,
	// and not any consistent type of character (e.g. 'n' will repeat, 'b' will not)
	window.SetCharCallback(func(_ *glfw.Window, char rune) {
		w.Send(key.Event{
			Rune: char,
		})
	})
	window.SetCloseCallback(func(*glfw.Window) {
		w.lifecycler.SetDead(true)
		w.lifecycler.SendEvent(w, nil)
	})

	// TODO(dmitshur): more fine-grained tracking of whether window is visible and/or focused
	w.lifecycler.SetDead(false)
	w.lifecycler.SetVisible(true)
	w.lifecycler.SetFocused(true)
	w.lifecycler.SendEvent(w, nil)

	// Send the initial size and paint events.
	width, height = window.GetFramebufferSize()
	framebufferSizeCallback(window, width, height)

	return w, nil
}

func optsSize(opts screen.WindowGenerator) (width, height int) {
	width, height = 1024/2, 768/2
	if opts.Width > 0 {
		width = opts.Width
	}
	if opts.Height > 0 {
		height = opts.Height
	}
	return width, height
}

func glfwMouseButton(button glfw.MouseButton) mouse.Button {
	switch button {
	case glfw.MouseButtonLeft:
		return mouse.ButtonLeft
	case glfw.MouseButtonRight:
		return mouse.ButtonRight
	case glfw.MouseButtonMiddle:
		return mouse.ButtonMiddle
	default:
		return mouse.ButtonNone
	}
}

func glfwMouseDirection(action glfw.Action) mouse.Direction {
	switch action {
	case glfw.Press:
		return mouse.DirPress
	case glfw.Release:
		return mouse.DirRelease
	default:
		panic("unreachable")
	}
}

var keyMap = map[glfw.Key]key.Code{
	glfw.KeyEscape:    key.CodeEscape,
	glfw.KeyEnter:     key.CodeReturnEnter,
	glfw.KeyTab:       key.CodeTab,
	glfw.KeyBackspace: key.CodeDeleteBackspace,
	glfw.KeyInsert:    key.CodeInsert,
	// note not differentiated from backspace
	glfw.KeyDelete:       key.CodeDeleteBackspace,
	glfw.KeyRight:        key.CodeRightArrow,
	glfw.KeyLeft:         key.CodeLeftArrow,
	glfw.KeyDown:         key.CodeDownArrow,
	glfw.KeyUp:           key.CodeUpArrow,
	glfw.KeyPageUp:       key.CodePageUp,
	glfw.KeyPageDown:     key.CodePageDown,
	glfw.KeyHome:         key.CodeHome,
	glfw.KeyEnd:          key.CodeEnd,
	glfw.KeyCapsLock:     key.CodeCapsLock,
	glfw.KeyNumLock:      key.CodeKeypadNumLock,
	glfw.KeyPause:        key.CodePause,
	glfw.KeyF1:           key.CodeF1,
	glfw.KeyF2:           key.CodeF2,
	glfw.KeyF3:           key.CodeF3,
	glfw.KeyF4:           key.CodeF4,
	glfw.KeyF5:           key.CodeF5,
	glfw.KeyF6:           key.CodeF6,
	glfw.KeyF7:           key.CodeF7,
	glfw.KeyF8:           key.CodeF8,
	glfw.KeyF9:           key.CodeF9,
	glfw.KeyF10:          key.CodeF10,
	glfw.KeyF11:          key.CodeF11,
	glfw.KeyF12:          key.CodeF12,
	glfw.KeyF13:          key.CodeF13,
	glfw.KeyF14:          key.CodeF14,
	glfw.KeyF15:          key.CodeF15,
	glfw.KeyF16:          key.CodeF16,
	glfw.KeyF17:          key.CodeF17,
	glfw.KeyF18:          key.CodeF18,
	glfw.KeyF19:          key.CodeF19,
	glfw.KeyF20:          key.CodeF20,
	glfw.KeyF21:          key.CodeF21,
	glfw.KeyF22:          key.CodeF22,
	glfw.KeyF23:          key.CodeF23,
	glfw.KeyF24:          key.CodeF24,
	glfw.KeyKPDecimal:    key.CodeKeypadFullStop,
	glfw.KeyKPEnter:      key.CodeKeypadEnter,
	glfw.KeyLeftShift:    key.CodeLeftShift,
	glfw.KeyLeftControl:  key.CodeLeftControl,
	glfw.KeyLeftAlt:      key.CodeLeftAlt,
	glfw.KeyLeftSuper:    key.CodeLeftGUI,
	glfw.KeyRightShift:   key.CodeRightShift,
	glfw.KeyRightControl: key.CodeRightControl,
	glfw.KeyRightAlt:     key.CodeRightAlt,
	glfw.KeyRightSuper:   key.CodeRightGUI,
	glfw.KeySpace:        key.CodeSpacebar,
	glfw.KeyApostrophe:   key.CodeApostrophe,
	glfw.KeyComma:        key.CodeComma,
	glfw.KeyMinus:        key.CodeHyphenMinus,
	glfw.KeyPeriod:       key.CodeFullStop,
	glfw.KeySlash:        key.CodeSlash,
	glfw.Key0:            key.Code0,
	glfw.Key1:            key.Code1,
	glfw.Key2:            key.Code2,
	glfw.Key3:            key.Code3,
	glfw.Key4:            key.Code4,
	glfw.Key5:            key.Code5,
	glfw.Key6:            key.Code6,
	glfw.Key7:            key.Code7,
	glfw.Key8:            key.Code8,
	glfw.Key9:            key.Code9,
	glfw.KeySemicolon:    key.CodeUnknown,
	glfw.KeyEqual:        key.CodeUnknown,
	glfw.KeyA:            key.CodeA,
	glfw.KeyB:            key.CodeB,
	glfw.KeyC:            key.CodeC,
	glfw.KeyD:            key.CodeD,
	glfw.KeyE:            key.CodeE,
	glfw.KeyF:            key.CodeF,
	glfw.KeyG:            key.CodeG,
	glfw.KeyH:            key.CodeH,
	glfw.KeyI:            key.CodeI,
	glfw.KeyJ:            key.CodeJ,
	glfw.KeyK:            key.CodeK,
	glfw.KeyL:            key.CodeL,
	glfw.KeyM:            key.CodeM,
	glfw.KeyN:            key.CodeN,
	glfw.KeyO:            key.CodeO,
	glfw.KeyP:            key.CodeP,
	glfw.KeyQ:            key.CodeQ,
	glfw.KeyR:            key.CodeR,
	glfw.KeyS:            key.CodeS,
	glfw.KeyT:            key.CodeT,
	glfw.KeyU:            key.CodeU,
	glfw.KeyV:            key.CodeV,
	glfw.KeyW:            key.CodeW,
	glfw.KeyX:            key.CodeX,
	glfw.KeyY:            key.CodeY,
	glfw.KeyZ:            key.CodeZ,
	glfw.KeyLeftBracket:  key.CodeLeftSquareBracket,
	glfw.KeyBackslash:    key.CodeBackslash,
	glfw.KeyRightBracket: key.CodeRightSquareBracket,
	glfw.KeyGraveAccent:  key.CodeGraveAccent,
	glfw.KeyScrollLock:   key.CodeUnknown,
	glfw.KeyPrintScreen:  key.CodeUnknown,
	glfw.KeyKP0:          key.CodeKeypad0,
	glfw.KeyKP1:          key.CodeKeypad1,
	glfw.KeyKP2:          key.CodeKeypad2,
	glfw.KeyKP3:          key.CodeKeypad3,
	glfw.KeyKP4:          key.CodeKeypad4,
	glfw.KeyKP5:          key.CodeKeypad5,
	glfw.KeyKP6:          key.CodeKeypad6,
	glfw.KeyKP7:          key.CodeKeypad7,
	glfw.KeyKP8:          key.CodeKeypad8,
	glfw.KeyKP9:          key.CodeKeypad9,
	glfw.KeyKPDivide:     key.CodeKeypadSlash,
	glfw.KeyKPMultiply:   key.CodeKeypadAsterisk,
	glfw.KeyKPSubtract:   key.CodeKeypadHyphenMinus,
	glfw.KeyKPAdd:        key.CodeKeypadPlusSign,
	glfw.KeyKPEqual:      key.CodeKeypadEqualSign,
}

func glfwKeyCode(k glfw.Key) key.Code {
	if kc, ok := keyMap[k]; ok {
		return kc
	}
	return key.CodeUnknown
}

func glfwKeyDirection(action glfw.Action) key.Direction {
	switch action {
	case glfw.Press:
		return key.DirPress
	case glfw.Release:
		return key.DirRelease
	case glfw.Repeat:
		return key.DirNone
	default:
		panic("unreachable")
	}
}

func glfwKeyMods(m glfw.ModifierKey) (mod key.Modifiers) {
	if m&glfw.ModAlt == glfw.ModAlt {
		mod |= key.ModAlt
	}
	if m&glfw.ModShift == glfw.ModShift {
		mod |= key.ModShift
	}
	if m&glfw.ModControl == glfw.ModControl {
		mod |= key.ModControl
	}
	if m&glfw.ModSuper == glfw.ModSuper {
		mod |= key.ModMeta
	}
	return mod
}
