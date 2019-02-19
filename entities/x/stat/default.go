package stat

var (
	// DefStatistics is a base set of statistics used by package-level calls
	// When using multiple statistics, avoid using overlapping event names
	DefStatistics = NewStatistics()
)

// Inc triggers an event, incrementing the given statistic by one
func Inc(eventName string) {
	DefStatistics.Inc(eventName)
}

// Trigger triggers the given event with a given increment to update a statistic
func Trigger(eventName string, inc int) {
	DefStatistics.Trigger(eventName, inc)
}

// TriggerOn triggers the given event, toggling it on
func TriggerOn(eventName string) {
	DefStatistics.TriggerOn(eventName)
}

// TriggerOff triggers the given event, toggling it off
func TriggerOff(eventName string) {
	DefStatistics.TriggerOff(eventName)
}

// TriggerTimed triggers the given event, toggling it on or off
func TriggerTimed(eventName string, on bool) {
	DefStatistics.TriggerTimed(eventName, on)
}

// TrackStats records a stat event to the Statistics map and creates the statistic if it does not already exist
func TrackStats(no int, data interface{}) int {
	return DefStatistics.TrackStats(no, data)
}

// TrackTimeStats acts like TrackStats, but tracks durations of events. If the
// event has not started, it logs a start time, and then when the event ends
// it will log the delta since the start.
func TrackTimeStats(no int, data interface{}) int {
	return DefStatistics.TrackTimeStats(no, data)
}

// IsTimedStat returns whether the given stat name is a part of this statistics'
// set of timed stats
func IsTimedStat(s string) bool {
	return DefStatistics.IsTimedStat(s)
}
