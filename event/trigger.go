package event

import (
	"time"

	"github.com/oakmound/oak/timing"
)

// Oak uses the following built in events:
//
// - EnterFrame: the beginning of every logical frame.
//   Payload: (int) frames passed since this scene started
//
// - CollisionStart/Stop: when a PhaseCollision entity starts/stops touching some label.
//   Payload: (collision.Label) the label the entity has started/stopped touching
//
// - MouseCollisionStart/Stop: as above, for mouse collision
//   Payload: (mouse.Event)
//
// - KeyDown/Up: when a key is pressed down / lifted up.
//   Payload: (string) the key pressed
//
// - Mouse events: MousePress, MouseRelease, MouseScrollDown, MouseScrollUp, MouseDrag
//   Payload: (mouse.Event) details on the mouse event
//
// - AnimationEnd: Triggered on animations CIDs when they loop from the last to the first frame
//   Payload: nil
//
// - ViewportUpdate: Triggered when the viewport changes.
//   Payload: []float64{viewportX, viewportY}

// Trigger an event, but only
// for one ID. Use case example:
// on onHit event
func (id CID) Trigger(eventName string, data interface{}) {

	go func(eventName string, data interface{}) {
		eb := GetBus()
		mutex.RLock()
		if idMap, ok := eb.bindingMap[eventName]; ok {
			if bs, ok := idMap[int(id)]; ok {
				for i := bs.highIndex - 1; i >= 0; i-- {
					for j, bnd := range (*bs.highPriority[i]).sl {
						handleBindable(bnd, int(id), data, j, eventName)
					}
				}
				triggerDefault((bs.defaultPriority).sl, int(id), eventName, data)

				for i := 0; i < bs.lowIndex; i++ {
					for j, bnd := range (*bs.lowPriority[i]).sl {
						handleBindable(bnd, int(id), data, j, eventName)
					}
				}
			}
		}
		mutex.RUnlock()
	}(eventName, data)
}

// TriggerAfter will trigger the given event after d time.
func (id CID) TriggerAfter(d time.Duration, eventName string, data interface{}) {
	go func() {
		timing.DoAfter(d, func() {
			id.Trigger(eventName, data)
		})
	}()
}

// Trigger is equivalent to bus.Trigger(...)
func Trigger(eventName string, data interface{}) {
	thisBus.Trigger(eventName, data)
}

// TriggerBack is equivalent to bus.TriggerBack(...)
func TriggerBack(eventName string, data interface{}) chan bool {
	return thisBus.TriggerBack(eventName, data)
}

// TriggerBack is a version of Trigger which returns a channel that
// informs on when all bindables have been called and returned from
// the input event. It is dangerous to use this unless you have a
// very good idea how things will synchronize, as if a triggered
// bindable itself makes a TriggerBack call, this will cause the engine to freeze,
// as the function will never end because the first TriggerBack has control of
// the lock for the event bus, and the first TriggerBack won't give up that lock
// until the function ends.
//
// This inherently means that when you call Trigger, the event will almost
// almost never be immediately triggered but rather will be triggered sometime
// soon in the future.
//
// TriggerBack is right now used by the primary logic loop to dictate logical
// framerate, so EnterFrame events are called through TriggerBack.
func (eb *Bus) TriggerBack(eventName string, data interface{}) chan bool {

	ch := make(chan bool)
	go func(ch chan bool, eb *Bus, eventName string, data interface{}) {
		ebtrigger(eb, eventName, data)
		ch <- true
	}(ch, eb, eventName, data)

	return ch
}

// Trigger will scan through the event bus and call all bindables found attached
// to the given event, with the passed in data.
func (eb *Bus) Trigger(eventName string, data interface{}) {
	go func(eb *Bus, eventName string, data interface{}) {
		ebtrigger(eb, eventName, data)
	}(eb, eventName, data)
}

func ebtrigger(eb *Bus, eventName string, data interface{}) {
	mutex.RLock()
	// Loop through all bindableStores for this eventName
	for id, bs := range (*eb).bindingMap[eventName] {
		// Top to bottom, high priority
		for i := bs.highIndex - 1; i >= 0; i-- {
			for j, bnd := range (*bs.highPriority[i]).sl {
				handleBindable(bnd, id, data, j, eventName)
			}
		}
	}

	for id, bs := range (*eb).bindingMap[eventName] {
		if bs != nil && bs.defaultPriority != nil {
			triggerDefault((bs.defaultPriority).sl, id, eventName, data)
		}
	}

	for id, bs := range (*eb).bindingMap[eventName] {
		// Bottom to top, low priority
		for i := 0; i < bs.lowIndex; i++ {
			for j, bnd := range (*bs.lowPriority[i]).sl {
				handleBindable(bnd, id, data, j, eventName)
			}
		}
	}
	mutex.RUnlock()
}

func triggerDefault(sl []Bindable, id int, eventName string, data interface{}) {
	progCh := make(chan bool)
	for i, bnd := range sl {
		go func(bnd Bindable, id int, eventName string, data interface{}, progCh chan bool, index int) {
			handleBindable(bnd, id, data, index, eventName)
			progCh <- true
		}(bnd, id, eventName, data, progCh, i)
	}
	for range sl {
		<-progCh
	}
}

func handleBindable(bnd Bindable, id int, data interface{}, index int, eventName string) {
	if bnd != nil {
		if id == 0 || GetEntity(id) != nil {
			response := bnd(id, data)
			switch response {
			case UnbindEvent:
				UnbindAll(BindingOption{
					Event{
						eventName,
						id,
					},
					0,
				})
			case UnbindSingle:
				Binding{
					BindingOption{
						Event{
							eventName,
							id,
						},
						0,
					},
					index,
				}.Unbind()
			}
		}
	}
}
