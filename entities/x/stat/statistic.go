package stat

import (
	"sync"
	"time"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
)

// Statistics stores the ongoing results of TrackStats and TrackTimeStats
type Statistics struct {
	stats    map[string]*History
	statLock sync.Mutex

	statTimes    map[string]time.Time
	statTimeLock sync.Mutex
}

// NewStatistics creates an empty statistics set
func NewStatistics() *Statistics {
	return &Statistics{
		stats:        make(map[string]*History),
		statLock:     sync.Mutex{},
		statTimes:    make(map[string]time.Time),
		statTimeLock: sync.Mutex{},
	}
}

// A History keeps track of any recorded occurrences of this statstic and their magnitude
type History struct {
	Name   string
	Events []Event
}

// An Event ties a value to a timestamp
type Event struct {
	Timestamp time.Time
	Val       int
}

// NewHistory creates a stat
func NewHistory(statName string, time time.Time) *History {
	return &History{Name: statName, Events: []Event{{time, 0}}}
}

// track adds a tracked event to the stat's history
func (h *History) track(t time.Time, v int) *History {
	if len(h.Events) > 0 {
		v += h.Events[len(h.Events)-1].Val
	}
	h.Events = append(h.Events, Event{t, v})
	return h
}

// Total takes a statistics history and finds the sum.
func (h *History) Total() int {
	return h.Events[len(h.Events)-1].Val
}

func (st *Statistics) trackStats(name string, val int) {
	st.statLock.Lock()
	stat, ok := st.stats[name]
	if !ok {
		stat = NewHistory(name, time.Now())
		st.stats[name] = stat
	}
	stat.track(time.Now(), val)
	st.statLock.Unlock()
}

// TrackStats records a stat event to the Statistics map and creates the statistic if it does not already exist
func (st *Statistics) TrackStats(no int, data interface{}) int {
	stat, ok := data.(stat)
	if !ok {
		dlog.Error("TrackStats called with a non-stat payload")
		return event.UnbindEvent
	}
	st.trackStats(stat.name, stat.inc)
	return 0
}

// TrackTimeStats acts like TrackStats, but tracks durations of events. If the
// event has not started, it logs a start time, and then when the event ends
// it will log the delta since the start.
func (st *Statistics) TrackTimeStats(no int, data interface{}) int {
	timed, ok := data.(timedStat)
	if !ok {
		dlog.Error("TrackTimeStats called with a non-timedStat payload")
		return event.UnbindEvent
	}
	if timed.on { //Turning on a thing to time track
		st.statTimeLock.Lock()
		st.statTimes[timed.name] = time.Now()
		st.statTimeLock.Unlock()
	} else {
		st.statTimeLock.Lock()
		timeDiff := int(time.Since(st.statTimes[timed.name]))
		st.statTimeLock.Unlock()
		if timeDiff < 0 {
			return 0
		}
		st.trackStats(timed.name, timeDiff)
	}
	return 0
}

// IsTimedStat returns whether the given stat name is a part of this statistics'
// set of timed stats
func (st *Statistics) IsTimedStat(s string) bool {
	_, ok := st.statTimes[s]
	return ok
}
