package key

import (
	"github.com/oakmound/oak/v3/event"
	"golang.org/x/mobile/event/key"
)

const (
	// Down is sent when a key is pressed. It is sent both as
	// Down, and as Down + the key name.
	Down = "KeyDown"
	// Up is sent when a key is released. It is sent both as
	// Up, and as Up + the key name.
	Up = "KeyUp"
	// Held is sent when a key is held down. It is sent both as
	// Held, and as Held + the key name.
	Held = "KeyHeld"
)

// An Event is sent as the payload for all key bindings.
type Event = key.Event

// A code is a unique integer code for a given common key
type Code = key.Code

// Binding will convert a function that accepts a typecast key.Event into a generic event binding
//
// Example:
// 		bus.Bind(key.Down, key.Binding(keyHandler))
func Binding(fn func(event.CID, Event) int) func(event.CID, interface{}) int {
	return func(cid event.CID, iface interface{}) int {
		ke, ok := iface.(Event)
		if !ok {
			return event.UnbindSingle
		}
		return fn(cid, ke)
	}
}
