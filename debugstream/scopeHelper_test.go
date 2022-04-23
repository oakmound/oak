package debugstream

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/debugtools"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/window"
)

type fakeWindow struct {
	window.Window

	skipCalls       int
	quitCalls       int
	fullscreenCalls int

	moveWindow func(x, y, w, h int)
}

func (f *fakeWindow) NextScene() {
	f.skipCalls++
}

func (f *fakeWindow) Quit() {
	f.quitCalls++
}

func (f *fakeWindow) EventHandler() event.Handler {
	return event.NewBus(nil)
}

func (f *fakeWindow) SetFullScreen(bool) error {
	f.fullscreenCalls++
	return nil
}

func (f *fakeWindow) MoveWindow(x, y, w, h int) error {
	f.moveWindow(x, y, w, h)
	return nil
}

func (f *fakeWindow) Bounds() intgeom.Point2 {
	return intgeom.Point2{1, 1}
}

func (f *fakeWindow) Viewport() intgeom.Point2 {
	return intgeom.Point2{}
}

func TestSkipCommands(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	in := bytes.NewBufferString("skip-scene\n1 skip-scene\nscope 1\nskip-scene")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := `assumed scope 1
`

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
	if fw.skipCalls != 3 {
		t.Fatal("expected 3 skips, got:", fw.skipCalls)
	}
}

func TestQuitCommands(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	in := bytes.NewBufferString("quit\n")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := ``

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
	if fw.quitCalls != 1 {
		t.Fatal("expected 1 quit, got:", fw.quitCalls)
	}
}

func TestMouseCommands(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	in := bytes.NewBufferString("mouse-details\n")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := ``

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
}

func TestFullScreen(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	in := bytes.NewBufferString("fullscreen\nfullscreen off")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := ``

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
	if fw.fullscreenCalls != 2 {
		t.Fatal("expected 2 fullscreens, got:", fw.fullscreenCalls)
	}
}

func TestMoveWindow(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	fw.moveWindow = func(x, y, w, h int) {
		if x != 1 {
			t.Fatal("x was not 1:", x)
		}
		if y != 2 {
			t.Fatal("y was not 2:", y)
		}
		if w != 3 {
			t.Fatal("w was not 3:", w)
		}
		if h != 4 {
			t.Fatal("h was not 4:", h)
		}
	}

	in := bytes.NewBufferString("move 1 2 3 4\nmove")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := `Must supply at least 2 coordinates`

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
}

func TestFadeCommands(t *testing.T) {
	sc := NewScopedCommands()
	fw := &fakeWindow{}
	sc.AddDefaultsForScope(1, fw)

	debugtools.SetDebugRenderable("test-r", render.EmptyRenderable())

	in := bytes.NewBufferString("fade test-r\nfade test-r 200\nfade\nfade bad-r")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(100 * time.Millisecond)

	expected := "Must supply at least 1 arguments\nCould not fade input bad-r\nPossible inputs are 'test-r'\n"

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
}
