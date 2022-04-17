package oak

import (
	"image"
	"os"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
)

func TestMouseClicks(t *testing.T) {
	c1 := NewWindow()
	c1.MouseTree = collision.NewTree()
	ch := make(chan struct{})
	c1.eventHandler = event.NewBus(event.NewCallerMap())
	bnd := event.GlobalBind(c1.eventHandler, mouse.Click, func(_ *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("click binding never bound")
	case <-bnd.Bound:
	}
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	c1.MouseTree.Add(sp)
	c1.Propagate(mouse.PressOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.ReleaseOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger click binding")
	case <-ch:
	}
}

func TestMouseClicksRelative(t *testing.T) {
	c1 := NewWindow()
	c1.MouseTree = collision.NewTree()
	ch := make(chan struct{})
	c1.eventHandler = event.NewBus(event.NewCallerMap())
	bnd := event.GlobalBind(c1.eventHandler, mouse.RelativeClickOn, func(_ *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("click binding never bound")
	case <-bnd.Bound:
	}
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	c1.MouseTree.Add(sp)
	defer c1.MouseTree.Clear()
	c1.Propagate(mouse.RelativePressOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.RelativeReleaseOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger click binding")
	case <-ch:
	}
}

type ent struct {
	event.CallerID
}

func TestPropagate(t *testing.T) {
	c1 := NewWindow()
	c1.eventHandler = event.NewBus(event.NewCallerMap())

	thisEnt := ent{}
	thisEnt.CallerID = c1.eventHandler.GetCallerMap().Register(thisEnt)
	ch := make(chan struct{})
	s := collision.NewSpace(10, 10, 10, 10, thisEnt.CallerID)
	event.Bind(c1.eventHandler, mouse.PressOn, thisEnt, func(ent, *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	c1.MouseTree = collision.NewTree()
	c1.MouseTree.Add(s)
	c1.Propagate(mouse.ReleaseOn, mouse.NewEvent(15, 15, mouse.ButtonLeft, mouse.Release))
	select {
	case <-ch:
		t.Fatalf("release propagation triggered press binding")
	case <-time.After(1 * time.Second):
	}
	c1.Propagate(mouse.PressOn, mouse.NewEvent(15, 15, mouse.ButtonLeft, mouse.Press))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger press binding")
	case <-ch:
	}
}

func TestPropagate_StopPropagation(t *testing.T) {
	c1 := NewWindow()
	c1.eventHandler = event.NewBus(event.NewCallerMap())
	c1.MouseTree = collision.NewTree()

	e1 := ent{}
	e1.CallerID = c1.eventHandler.GetCallerMap().Register(e1)
	e2 := ent{}
	e2.CallerID = c1.eventHandler.GetCallerMap().Register(e2)

	s1 := collision.NewSpace(10, 10, 10, 10, e1.CallerID)
	s1.SetZLayer(10)
	c1.MouseTree.Insert(s1)
	s2 := collision.NewSpace(10, 10, 10, 10, e2.CallerID)
	s2.SetZLayer(1)
	c1.MouseTree.Insert(s2)
	var failed bool
	<-event.Bind(c1.eventHandler, mouse.PressOn, e1, func(_ ent, ev *mouse.Event) event.Response {
		ev.StopPropagation = true
		return 0
	}).Bound
	<-event.Bind(c1.eventHandler, mouse.PressOn, e2, func(_ ent, ev *mouse.Event) event.Response {
		failed = true
		return 0
	}).Bound
	<-event.Bind(c1.eventHandler, mouse.ClickOn, e1, func(_ ent, ev *mouse.Event) event.Response {
		ev.StopPropagation = true
		return 0
	}).Bound
	<-event.Bind(c1.eventHandler, mouse.ClickOn, e2, func(_ ent, ev *mouse.Event) event.Response {
		failed = true
		return 0
	}).Bound
	<-event.Bind(c1.eventHandler, mouse.RelativeClickOn, e1, func(_ ent, ev *mouse.Event) event.Response {
		ev.StopPropagation = true
		return 0
	}).Bound
	<-event.Bind(c1.eventHandler, mouse.RelativeClickOn , e2, func(_ ent, ev *mouse.Event) event.Response {
		failed = true
		return 0
	}).Bound
	c1.TriggerMouseEvent(mouse.Event{
		Point2: floatgeom.Point2{
			15, 15,
		},
		Button:    mouse.ButtonLeft,
		EventType: mouse.Press,
	})
	c1.TriggerMouseEvent(mouse.Event{
		Point2: floatgeom.Point2{
			15, 15,
		},
		Button:    mouse.ButtonLeft,
		EventType: mouse.Release,
	})
	if failed {
		t.Fatal("stop propagation failed")
	}
}

func TestWindowGetters(t *testing.T) {
	c1 := NewWindow()
	c1.debugConsole(os.Stdin, os.Stdout)
	if c1.InFocus() {
		t.Errorf("new windows should not be in focus")
	}
	if c1.EventHandler() != event.DefaultBus {
		t.Errorf("new windows should have the default event bus")
	}
	if c1.GetBackgroundImage() != image.Black {
		t.Errorf("new windows should have a black background")
	}
	c1.SetColorBackground(image.White)
	if c1.GetBackgroundImage() != image.White {
		t.Errorf("set color background failed")
	}
	rend := render.EmptyRenderable()
	c1.SetLoadingRenderable(rend)
	if c1.LoadingR != rend {
		t.Errorf("Set loading renderable failed")
	}
	c1.SetBackground(rend)
	r, g, b, a := c1.bkgFn().At(0, 0).RGBA()
	if r != 0 || g != 0 || b != 0 || a != 0 {
		t.Errorf("background was not set to empty renderable")
	}
}
