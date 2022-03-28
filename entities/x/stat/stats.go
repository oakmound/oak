package stat

import (
	"fmt"

	"github.com/oakmound/oak/v3/event"
)

// TODO: these functions are useless unless the types are exported, and
// if the types are exported the api is bad

type timedStatEvent struct {
	event event.EventID[timedStat]
	fmt.Stringer
}

type timedStat struct {
	name string
	on   bool
}

type statEvent struct {
	event event.EventID[stat]
	fmt.Stringer
}

type stat struct {
	name string
	inc  int
}

// TimedOn returns a binding that will trigger toggling on the given event
func TimedOn(ev timedStatEvent) event.UnsafeBindable {
	return TimedBind(ev, true)
}

// TimedOff returns a binding that will trigger toggling off the given event
func TimedOff(ev timedStatEvent) event.UnsafeBindable {
	return TimedBind(ev, false)
}

// TimedBind returns a binding that will trigger toggling on or off the given event
func TimedBind(ev timedStatEvent, on bool) event.UnsafeBindable {
	return func(event.CallerID, event.Handler, interface{}) event.Response {
		event.TriggerOn(event.DefaultBus, ev.event, timedStat{ev.String(), on})
		return 0
	}
}

// Bind returns a binding that will increment the given event by 'inc'
func Bind(ev statEvent, inc int) event.UnsafeBindable {
	return func(event.CallerID, event.Handler, interface{}) event.Response {
		event.TriggerOn(event.DefaultBus, ev.event, stat{ev.String(), inc})
		return 0
	}
}

// Inc triggers an event, incrementing the given statistic by one
func (st *Statistics) Inc(ev statEvent) {
	st.Trigger(ev, 1)
}

// Trigger triggers the given event with a given increment to update a statistic
func (st *Statistics) Trigger(ev statEvent, inc int) {
	event.TriggerOn(event.DefaultBus, ev.event, stat{ev.String(), inc})
}

// TriggerOn triggers the given event, toggling it on
func (st *Statistics) TriggerOn(ev timedStatEvent) {
	st.TriggerTimed(ev, true)
}

// TriggerOff triggers the given event, toggling it off
func (st *Statistics) TriggerOff(ev timedStatEvent) {
	st.TriggerTimed(ev, false)
}

// TriggerTimed triggers the given event, toggling it on or off
func (st *Statistics) TriggerTimed(ev timedStatEvent, on bool) {
	event.TriggerOn(event.DefaultBus, ev.event, timedStat{ev.String(), on})
}
