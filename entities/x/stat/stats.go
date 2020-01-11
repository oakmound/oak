package stat

import "github.com/oakmound/oak/v2/event"

type timedStat struct {
	name string
	on   bool
}
type stat struct {
	name string
	inc  int
}

// TimedOn returns a binding that will trigger toggling on the given event
func TimedOn(eventName string) event.Bindable {
	return TimedBind(eventName, true)
}

// TimedOff returns a binding that will trigger toggling off the given event
func TimedOff(eventName string) event.Bindable {
	return TimedBind(eventName, false)
}

// TimedBind returns a binding that will trigger toggling on or off the given event
func TimedBind(eventName string, on bool) event.Bindable {
	return func(int, interface{}) int {
		event.Trigger(eventName, timedStat{eventName, on})
		return 0
	}
}

// Bind returns a binding that will increment the given event by 'inc'
func Bind(eventName string, inc int) event.Bindable {
	return func(int, interface{}) int {
		event.Trigger(eventName, stat{eventName, inc})
		return 0
	}
}

// Inc triggers an event, incrementing the given statistic by one
func (st *Statistics) Inc(eventName string) {
	st.Trigger(eventName, 1)
}

// Trigger triggers the given event with a given increment to update a statistic
func (st *Statistics) Trigger(eventName string, inc int) {
	event.Trigger(eventName, stat{eventName, inc})
}

// TriggerOn triggers the given event, toggling it on
func (st *Statistics) TriggerOn(eventName string) {
	st.TriggerTimed(eventName, true)
}

// TriggerOff triggers the given event, toggling it off
func (st *Statistics) TriggerOff(eventName string) {
	st.TriggerTimed(eventName, false)
}

// TriggerTimed triggers the given event, toggling it on or off
func (st *Statistics) TriggerTimed(eventName string, on bool) {
	event.Trigger(eventName, timedStat{eventName, on})
}
