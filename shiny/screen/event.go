package screen

// EventDeque is an infinitely buffered double-ended queue of events.
type EventDeque interface {
	// Send adds an event to the end of the deque. They are returned by
	// NextEvent in FIFO order.
	Send(event interface{})

	// SendFirst adds an event to the start of the deque. They are returned by
	// NextEvent in LIFO order, and have priority over events sent via Send.
	SendFirst(event interface{})

	// NextEvent returns the next event in the deque. It blocks until such an
	// event has been sent.
	//
	// Typical event types include:
	//	- lifecycle.Event
	//	- size.Event
	//	- paint.Event
	//	- key.Event
	//	- mouse.Event
	//	- touch.Event
	// from the golang.org/x/mobile/event/... packages. Other packages may send
	// events, of those types above or of other types, via Send or SendFirst.
	NextEvent() interface{}
}
