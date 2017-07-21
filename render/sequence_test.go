package render

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
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

func TestSequenceFunctions(t *testing.T) {
	rgba1 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	rgba2 := image.NewRGBA(image.Rect(0, 0, 5, 5))
	sq := NewSequence(
		[]Modifiable{
			NewSprite(0, 0, rgba1),
			NewSprite(0, 0, rgba2),
		}, 5)
	sq2 := sq.Copy().(*Sequence)
	assert.Equal(t, sq.Get(0).GetRGBA(), rgba1)
	assert.Equal(t, sq2.Get(0).GetRGBA(), rgba1)
	assert.Equal(t, sq.GetRGBA(), rgba1)
	assert.Equal(t, sq2.GetRGBA(), rgba1)
	assert.Equal(t, sq.Get(1).GetRGBA(), rgba2)
	assert.Equal(t, sq2.Get(1).GetRGBA(), rgba2)
	time.Sleep(1 * time.Second)
	sq.update()
	assert.Equal(t, sq.GetRGBA(), rgba2)
	sq.Pause()
	time.Sleep(1 * time.Second)
	sq.update()
	assert.Equal(t, sq.GetRGBA(), rgba2)
	sq.Unpause()
	time.Sleep(1 * time.Second)
	sq.update()
	assert.Equal(t, sq.GetRGBA(), rgba1)

	// Sequences have 6,6 dimensions right now. See LayeredPoint.
	w, h := sq.GetDims()
	assert.Equal(t, w, 6)
	assert.Equal(t, h, 6)

	assert.Equal(t, sq.IsStatic(), false)
}

func TestSequenceModify(t *testing.T) {
	rgba1 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	rgba2 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	sq := NewSequence(
		[]Modifiable{
			NewSprite(0, 0, rgba1),
			NewSprite(0, 0, rgba2),
		}, 5)
	sq.Modify(CutRel(.5, .5))
	w, h := sq.Get(0).GetDims()
	assert.Equal(t, w, 5)
	assert.Equal(t, h, 5)
}

func TestTweenSequence(t *testing.T) {
	start := NewColorBox(10, 10, color.RGBA{0, 0, 0, 0})
	end := NewColorBox(10, 10, color.RGBA{255, 255, 255, 255})
	TweenSequence(start.GetRGBA(), end.GetRGBA(), 2, 5)
	// Tween behavior is tested elsewhere, this is just a "this doesn't crash" test
}
