package event_test

import (
	"math/rand"
	"testing"

	"github.com/oakmound/oak/v3/event"
)

func TestCallerID_CID(t *testing.T) {
	t.Run("Identity", func(t *testing.T) {
		c := event.CallerID(rand.Intn(100000))
		if c != c.CID() {
			t.Fatalf("callerID did not match itself: was %v, got %v", c, c.CID())
		}
	})
}

func TestNewCallerMap(t *testing.T) {
	t.Run("NotNil", func(t *testing.T) {
		m := event.NewCallerMap()
		if m == nil {
			t.Fatalf("created caller map was nil")
		}
	})
}

func randomCallerID() *event.CallerID {
	c1 := event.CallerID(rand.Intn(10000))
	return &c1
}

func TestCallerMap_Register(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		m := event.NewCallerMap()
		c1 := randomCallerID()
		m.Register(c1)
		c2 := m.GetEntity(c1.CID())
		if c2 != c1 {
			t.Fatalf("unable to retrieve registered caller")
		}
		if !m.HasEntity(c1.CID()) {
			t.Fatalf("caller map does not have registered caller")
		}
	})
	t.Run("Remove", func(t *testing.T) {
		m := event.NewCallerMap()
		c1 := randomCallerID()
		m.Register(c1)
		m.RemoveEntity(c1.CID())
		c3 := m.GetEntity(c1.CID())
		if c3 != nil {
			t.Fatalf("get entity had registered caller after remove")
		}
		if m.HasEntity(c1.CID()) {
			t.Fatalf("caller map has registered caller after remove")
		}
	})
	t.Run("Clear", func(t *testing.T) {
		m := event.NewCallerMap()
		c1 := randomCallerID()
		m.Register(c1)
		m.Clear()
		c3 := m.GetEntity(c1.CID())
		if c3 != nil {
			t.Fatalf("get entity had registered caller after clear")
		}
		if m.HasEntity(c1.CID()) {
			t.Fatalf("caller map has registered caller after clear")
		}
	})
}
