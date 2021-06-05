package oak

import (
	"bytes"
	"testing"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
)

type ent struct{}

func (e ent) Init() event.CID {
	return 0
}

func TestDebugConsole(t *testing.T) {
	c1 := NewController()
	c1.config.LoadBuiltinCommands = true
	triggered := false
	err := c1.AddCommand("test", func([]string) {
		triggered = true
	})
	if err != nil {
		t.Fatalf("failed to add test command")
	}

	render.UpdateDebugMap("r", render.EmptyRenderable())

	event.NextID(ent{})

	r := bytes.NewBufferString(
		"test\n" +
			"nothing\n" +
			"fade nothing\n" +
			"fade nothing 100\n" +
			"fade r\n" +
			"skip nothing\n" +
			"print nothing\n" +
			"print 2\n" +
			"print 1\n" +
			"mouse nothing\n" +
			"mouse details\n" +
			"garbage input\n" +
			"\n" +
			"skip scene\n")
	go c1.debugConsole(r)
	sleep()
	sleep()
	if !triggered {
		t.Fatalf("debug console did not trigger test command")
	}
	<-c1.skipSceneCh
}

func TestMouseDetails(t *testing.T) {
	c1 := NewController()

	ev := mouse.NewZeroEvent(0, 0)
	c1.mouseDetails(0, &ev)
	s := collision.NewUnassignedSpace(-1, -1, 2, 2)
	collision.Add(s)
	c1.mouseDetails(0, &ev)
	collision.Remove(s)

	// This should spew this nothing entity, but it doesn't.
	id := event.NextID(ent{})
	s = collision.NewSpace(-1, -1, 2, 2, id)
	c1.mouseDetails(0, &ev)
	collision.Remove(s)
}
