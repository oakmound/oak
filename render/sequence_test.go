package render

import (
	"image/color"
	"testing"
	"time"

	"github.com/oakmound/oak/event"
)

type Dummy struct{}

func (d Dummy) Init() event.CID {
	return event.NextID(d)
}

func TestSequenceTrigger(t *testing.T) {
	sq := NewSequence(
		// This syntax is bad for external calls,
		// we could swap fps and make this variadic
		[]Modifiable{
			NewColorBox(10, 10, color.RGBA{255, 0, 0, 255}),
			NewColorBox(10, 10, color.RGBA{0, 255, 0, 255}),
		}, 5)
	go event.ResolvePending()
	cid := Dummy{}.Init()
	sq.SetTriggerID(cid)
	triggerCh := make(chan bool)
	cid.Bind(func(int, interface{}) int {
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
		triggerCh <- true
		return 0
	}, "AnimationEnd")
	// We sleep to trigger the sequence to want to animate to the next frame
	time.Sleep(1 * time.Second)
	// Normally update is only called inside of Draw, so this is a pretend draw
	sq.update()
	time.Sleep(1 * time.Second)
	sq.update()
	<-triggerCh
}
