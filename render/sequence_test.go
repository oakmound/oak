package render

import (
	"image"
	"image/color"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render/mod"
)

type Dummy struct{}

func (d Dummy) Init() event.CallerID {
	return event.NextID(d)
}

func TestSequenceTrigger(t *testing.T) {
	sq := NewSequence(5,
		NewColorBox(10, 10, color.RGBA{255, 0, 0, 255}),
		NewColorBox(10, 10, color.RGBA{0, 255, 0, 255}))
	cid := Dummy{}.Init()
	sq.SetTriggerID(cid)
	triggerCh := make(chan struct{})
	cid.Bind(event.AnimationEnd, func(event.CallerID, interface{}) int {
		// This is a bad idea in real code, this will lock up
		// unbindings because the function that triggered this owns
		// the lock on the event bus until this function exits.
		// It is for this reason that all triggers, bindings,
		// and unbindings do nothing when they are called, just put
		// off work to be done-- to make sure no one is expecting a
		// result from one of those functions, from within a triggered
		// function, causing a deadlock.
		//
		// For this test this is the easiest way to do this though
		triggerCh <- struct{}{}
		return 0
	})
	// We sleep to trigger the sequence to want to animate to the next frame
	time.Sleep(1 * time.Second)
	// Normally update is only called inside of Draw, so this is a pretend draw
	sq.update()
	time.Sleep(1 * time.Second)
	sq.update()
	<-triggerCh
}

func TestSequenceFunctions(t *testing.T) {
	rgba1 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	rgba2 := image.NewRGBA(image.Rect(0, 0, 5, 5))
	sq := NewSequence(5,
		NewSprite(0, 0, rgba1),
		NewSprite(0, 0, rgba2))
	sq2 := sq.Copy().(*Sequence)
	if sq.Get(0).GetRGBA() != rgba1 {
		t.Fatalf("rgba mismatch")
	}
	if !reflect.DeepEqual(sq2.Get(0).GetRGBA(), rgba1) {
		t.Fatalf("rgba mismatch")
	}
	if sq.GetRGBA() != rgba1 {
		t.Fatalf("rgba mismatch")
	}
	if !reflect.DeepEqual(sq2.Get(0).GetRGBA(), rgba1) {
		t.Fatalf("rgba mismatch")
	}
	if sq.Get(1).GetRGBA() != rgba2 {
		t.Fatalf("rgba mismatch")
	}
	if !reflect.DeepEqual(sq2.Get(1).GetRGBA(), rgba2) {
		t.Fatalf("rgba mismatch")
	}
	time.Sleep(1 * time.Second)
	sq.update()
	if sq.GetRGBA() != rgba2 {
		t.Fatalf("rgba mismatch")
	}
	sq.Pause()
	time.Sleep(1 * time.Second)
	sq.update()
	if sq.GetRGBA() != rgba2 {
		t.Fatalf("rgba mismatch")
	}
	sq.Unpause()
	time.Sleep(1 * time.Second)
	sq.update()
	if sq.GetRGBA() != rgba1 {
		t.Fatalf("rgba mismatch")
	}
	sq.SetFPS(.5)
	time.Sleep(1 * time.Second)
	sq.update()
	if sq.GetRGBA() != rgba1 {
		t.Fatalf("rgba mismatch")
	}
	time.Sleep(1 * time.Second)
	sq.update()
	if sq.GetRGBA() != rgba2 {
		t.Fatalf("rgba mismatch")
	}

	w, h := sq.GetDims()
	if w != 5 || h != 5 {
		t.Fatalf("get dims mismatch")
	}

	if sq.IsStatic() {
		t.Fatalf("sequence should not have been static")
	}

	if sq.Get(-1) != nil {
		t.Fatalf("get -1 should return nil")
	}
	if sq.Get(math.MaxInt32) != nil {
		t.Fatalf("get math max should return nil")
	}
}

func TestSequenceModify(t *testing.T) {
	rgba1 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	rgba2 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	sq := NewSequence(5,
		NewSprite(0, 0, rgba1),
		NewSprite(0, 0, rgba2))
	sq.Modify(mod.CutRel(.5, .5))
	w, h := sq.Get(0).GetDims()
	if w != 5 || h != 5 {
		t.Fatalf("get dims mismatch")
	}

	sq.Filter(mod.Brighten(100))
}

func TestTweenSequence(t *testing.T) {
	start := NewColorBox(10, 10, color.RGBA{0, 0, 0, 0})
	end := NewColorBox(10, 10, color.RGBA{255, 255, 255, 255})
	TweenSequence(start.GetRGBA(), end.GetRGBA(), 2, 5)
	// Tween behavior is tested elsewhere, this is just a "this doesn't crash" test
}
