package stat

import "github.com/oakmound/oak/v3/event"

var (
	// DefStatistics is a base set of statistics used by package-level calls
	// When using multiple statistics, avoid using overlapping event names
	DefStatistics = NewStatistics()
)

// Inc triggers an event, incrementing the given statistic by one
func Inc(ev statEvent) {
	DefStatistics.Inc(ev)
}

// Trigger triggers the given event with a given increment to update a statistic
func Trigger(ev statEvent, inc int) {
	DefStatistics.Trigger(ev, inc)
}

// TriggerOn triggers the given event, toggling it on
func TriggerOn(ev timedStatEvent) {
	DefStatistics.TriggerOn(ev)
}

// TriggerOff triggers the given event, toggling it off
func TriggerOff(ev timedStatEvent) {
	DefStatistics.TriggerOff(ev)
}

// TriggerTimed triggers the given event, toggling it on or off
func TriggerTimed(ev timedStatEvent, on bool) {
	DefStatistics.TriggerTimed(ev, on)
}

// TrackStats records a stat event to the Statistics map and creates the statistic if it does not already exist
func TrackStats(no int, data interface{}) event.Response {
	return DefStatistics.TrackStats(no, data)
}

// TrackTimeStats acts like TrackStats, but tracks durations of events. If the
// event has not started, it logs a start time, and then when the event ends
// it will log the delta since the start.
func TrackTimeStats(no int, data interface{}) event.Response {
	return DefStatistics.TrackTimeStats(no, data)
}

// IsTimedStat returns whether the given stat name is a part of this statistics'
// set of timed stats
func IsTimedStat(s string) bool {
	return DefStatistics.IsTimedStat(s)
}
