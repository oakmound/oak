package key

import "golang.org/x/mobile/event/key"

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
