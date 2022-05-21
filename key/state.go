package key

import (
	"sync"
	"time"
)

// A State tracks what keys of a keyboard are currently pressed and for how long they have been
// pressed if they are held down.
type State struct {
	state        map[Code]bool
	durations    map[Code]time.Time
	stateLock    sync.RWMutex
	durationLock sync.RWMutex
}

// NewState creates a state object for tracking keyboard state.
func NewState() State {
	return State{
		state:     make(map[Code]bool),
		durations: make(map[Code]time.Time),
	}
}

// SetUp will cause later IsDown calls to report false
// for the given key. This is called internally when
// events are sent from the real keyboard and mouse.
// Calling this can interrupt real input or cause
// unintended behavior and should be done cautiously.
func (ks *State) SetUp(key Code) {
	ks.stateLock.Lock()
	ks.durationLock.Lock()
	delete(ks.state, key)
	delete(ks.durations, key)
	ks.durationLock.Unlock()
	ks.stateLock.Unlock()
}

// SetDown will cause later IsDown calls to report true
// for the given key. This is called internally when
// events are sent from the real keyboard and mouse.
// Calling this can interrupt real input or cause
// unintended behavior and should be done cautiously.
func (ks *State) SetDown(key Code) {
	ks.stateLock.Lock()
	ks.state[key] = true
	ks.durations[key] = time.Now()
	ks.stateLock.Unlock()
}

// IsDown returns whether a key is held down
func (ks *State) IsDown(key Code) (k bool) {
	ks.stateLock.RLock()
	k = ks.state[key]
	ks.stateLock.RUnlock()
	return
}

// IsHeld returns whether a key is held down, and for how long
// it has been held.
func (ks *State) IsHeld(key Code) (k bool, d time.Duration) {
	ks.stateLock.RLock()
	k = ks.state[key]
	ks.stateLock.RUnlock()
	if k {
		ks.durationLock.RLock()
		d = time.Since(ks.durations[key])
		ks.durationLock.RUnlock()
	}
	return
}
