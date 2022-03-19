package key

import (
	"sync"

	"github.com/oakmound/oak/v3/event"
	"golang.org/x/mobile/event/key"
)

var (
	// Down is sent when a key is pressed. It is sent both as
	// Down, and as Down + the key name.
	AnyDown = event.RegisterEvent[Event]()
	// Up is sent when a key is released. It is sent both as
	// Up, and as Up + the key name.
	AnyUp = event.RegisterEvent[Event]()
	// Held is sent when a key is held down. It is sent both as
	// Held, and as Held + the key name.
	AnyHeld = event.RegisterEvent[Event]()
)

// An Event is sent as the payload for all key bindings.
type Event = key.Event

// A code is a unique integer code for a given common key
type Code = key.Code

var upEventsLock sync.Mutex
var upEvents = map[Code]event.EventID[Event]{}

func Up(code Code) event.EventID[Event] {
	upEventsLock.Lock()
	defer upEventsLock.Unlock()
	if ev, ok := upEvents[code]; ok {
		return ev
	}
	ev := event.RegisterEvent[Event]()
	upEvents[code] = ev
	return ev
}

var downEventsLock sync.Mutex
var downEvents = map[Code]event.EventID[Event]{}

func Down(code Code) event.EventID[Event] {
	downEventsLock.Lock()
	defer downEventsLock.Unlock()
	if ev, ok := downEvents[code]; ok {
		return ev
	}
	ev := event.RegisterEvent[Event]()
	downEvents[code] = ev
	return ev
}

var heldEventsLock sync.Mutex
var heldEvents = map[Code]event.EventID[Event]{}

func Held(code Code) event.EventID[Event] {
	heldEventsLock.Lock()
	defer heldEventsLock.Unlock()
	if ev, ok := heldEvents[code]; ok {
		return ev
	}
	ev := event.RegisterEvent[Event]()
	heldEvents[code] = ev
	return ev
}
