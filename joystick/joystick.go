// Package joystick provides utilities for querying and reacting to joystick or
// gamepad inputs.
package joystick

import (
	"math"
	"time"

	"github.com/oakmound/oak/v2/dlog"
)

// Events. All events include a *State payload.
const (
	Change          = "JoystickChange"
	ButtonDown      = "ButtonDown"
	ButtonUp        = "ButtonUp"
	RtTriggerChange = "RtTriggerChange"
	LtTriggerChange = "LtTriggerChange"
	RtStickChange   = "RtStickChange"
	LtStickChange   = "LtStickChange"
	Disconnected    = "JoystickDisconnected"
)

// Init calls any os functions necessary to detect joysticks
func Init() error {
	return osinit()
}

// A Triggerer can either be an event bus or event CID, allowing
// joystick triggers to be listened to globally or sent to particular entities.
type Triggerer interface {
	Trigger(string, interface{})
}

// A Joystick represents a (usually) physical controller connected to the machine.
// It can be listened to to obtain button / trigger states of the controller
// and propagate changes through an event handler.
type Joystick struct {
	Handler  Triggerer
	PollRate time.Duration
	id       uint32
	osJoystick
}

// State represents a snapshot of a joystick's inputs
type State struct {
	// Frame advances every time the State changes,
	// Making it easier to tell if the state has not changed
	// since the last poll.
	Frame uint32
	// ID is the joystick ID this state is associated with
	ID uint32
	// There isn't a canonical list of buttons (yet), because
	// this can be so dependent on controller type. To get a list of expected
	// buttons for a joystick, this map can be checked after a single GetState().
	// The package will ensure that the keyset of this map will never change for
	// a given joystick variable. True = Pressed/Down. False = Released/Up.
	// Todo: ensure the above guarantee
	Buttons  map[string]bool
	TriggerL uint8
	TriggerR uint8
	StickLX  int16
	StickLY  int16
	StickRX  int16
	StickRY  int16
}

// ListenOptions can be passed into a joystick's Listen method to
// change what events will be propagated out.
type ListenOptions struct {
	// Each boolean dictates given event type(s) will be sent during Listen
	// "JoystickChange": *State
	JoystickChanges bool
	// Whenever any button is changed:
	// "ButtonDown": *State
	// "ButtonUp": *State
	GenericButtonPresses bool
	// "($buttonName)Up/Down": *State
	// e.g.
	// "AButtonUp": *State
	// "XButtonDown:" *State
	ButtonPresses bool
	// "RtTriggerChange": *State
	// "LtTriggerChange": *State
	TriggerChanges bool
	// "RtStickChange": *State
	// "LtStickChange": *State
	StickChanges bool
}

func (lo *ListenOptions) sendFn() func(Triggerer, *State, *State) {
	// Todo: benchmark that this is an effective way of reducing time spent
	// sending unwanted events
	var fn func(Triggerer, *State, *State)
	if lo.JoystickChanges {
		fn = func(h Triggerer, cur, last *State) {
			h.Trigger(Change, cur)
		}
	}
	if lo.GenericButtonPresses {
		prevFn := fn
		if prevFn != nil {
			fn = func(h Triggerer, cur, last *State) {
				prevFn(h, cur, last)
				var downTriggered bool
				var upTriggered bool
				for k, v := range cur.Buttons {
					if v != last.Buttons[k] {
						if v && !downTriggered {
							h.Trigger(ButtonDown, cur)
							downTriggered = true
							if upTriggered {
								return
							}
						} else if !v && !upTriggered {
							h.Trigger(ButtonUp, cur)
							upTriggered = true
							if downTriggered {
								return
							}
						}
					}
				}
			}
		} else {
			fn = func(h Triggerer, cur, last *State) {
				var downTriggered bool
				var upTriggered bool
				for k, v := range cur.Buttons {
					if v != last.Buttons[k] {
						if v && !downTriggered {
							h.Trigger(ButtonDown, cur)
							downTriggered = true
							if upTriggered {
								return
							}
						} else if !v && !upTriggered {
							h.Trigger(ButtonUp, cur)
							upTriggered = true
							if downTriggered {
								return
							}
						}
					}
				}
			}
		}
	}
	if lo.ButtonPresses {
		prevFn := fn
		if prevFn != nil {
			fn = func(h Triggerer, cur, last *State) {
				prevFn(h, cur, last)
				for k, v := range cur.Buttons {
					if v != last.Buttons[k] {
						if v {
							h.Trigger(k+ButtonDown, cur)
						} else {
							h.Trigger(k+ButtonUp, cur)
						}
					}
				}
			}
		} else {
			fn = func(h Triggerer, cur, last *State) {
				for k, v := range cur.Buttons {
					if v != last.Buttons[k] {
						if v {
							h.Trigger(k+ButtonDown, cur)
						} else {
							h.Trigger(k+ButtonUp, cur)
						}
					}
				}
			}
		}
	}
	// Todo: support a stick deadzone where we don't report very tiny changes
	// near the center of the stick
	if lo.StickChanges {
		prevFn := fn
		if prevFn != nil {
			fn = func(h Triggerer, cur, last *State) {
				prevFn(h, cur, last)
				if cur.StickLX != last.StickLX ||
					cur.StickLY != last.StickLY {
					h.Trigger(LtStickChange, cur)
				}
				if cur.StickRX != last.StickRX ||
					cur.StickRY != last.StickRY {
					h.Trigger(RtStickChange, cur)
				}
			}
		} else {
			fn = func(h Triggerer, cur, last *State) {
				if cur.StickLX != last.StickLX ||
					cur.StickLY != last.StickLY {
					h.Trigger(LtStickChange, cur)
				}
				if cur.StickRX != last.StickRX ||
					cur.StickRY != last.StickRY {
					h.Trigger(RtStickChange, cur)
				}
			}
		}
	}
	if lo.TriggerChanges {
		prevFn := fn
		if prevFn != nil {
			fn = func(h Triggerer, cur, last *State) {
				prevFn(h, cur, last)
				if cur.TriggerL != last.TriggerL {
					h.Trigger(LtTriggerChange, cur)
				}
				if cur.TriggerR != last.TriggerR {
					h.Trigger(RtTriggerChange, cur)
				}
			}
		} else {
			fn = func(h Triggerer, cur, last *State) {
				if cur.TriggerL != last.TriggerL {
					h.Trigger(LtTriggerChange, cur)
				}
				if cur.TriggerR != last.TriggerR {
					h.Trigger(RtTriggerChange, cur)
				}
			}
		}
	}
	return fn
}

// Listen causes the joystick to send its inputs to its Handler, by regularly
// querying GetState. The type of events returned can be specified by options.
// If the options are nil, only JoystickChange events will be sent.
func (j *Joystick) Listen(opts *ListenOptions) (cancel func()) {
	if opts == nil {
		opts = &ListenOptions{
			JoystickChanges: true,
		}
	}
	stop := make(chan struct{})
	sendFn := opts.sendFn()
	go func() {
		// Perform required initialization to receive inputs from OS
		dlog.ErrorCheck(j.Prepare())
		t := time.NewTicker(j.PollRate)
		lastState := &State{Frame: math.MaxUint32}
		for {
			// Wait on inputs
			select {
			case <-t.C:
			case <-stop:
				t.Stop()
				err := j.Close()
				dlog.ErrorCheck(err)
				close(stop)
				return
			}
			state, err := j.GetState()
			if err != nil {
				j.Handler.Trigger(Disconnected, j.id)
				dlog.Error(err)
				t.Stop()
				j.Close()
				return
			}
			if lastState.Frame == state.Frame {
				continue
			}
			sendFn(j.Handler, state, lastState)
			lastState = state
		}
	}()
	cancel = func() {
		stop <- struct{}{}
	}
	return cancel
}

// ID returns which player this joystick is associated with
func (j *Joystick) ID() uint32 {
	return j.id
}

// Vibrate triggers vibration on a joystick (if it is supported).
func (j *Joystick) Vibrate(left, right uint16) error {
	return j.vibrate(left, right)
}

// Prepare allocates any operating system resources needed to read signals
// for the joystick.
func (j *Joystick) Prepare() error {
	return j.prepare()
}

// GetState returns the current button, trigger, and analog stick
// state of the joystick
func (j *Joystick) GetState() (*State, error) {
	return j.getState()
}

// Close closes any operating system resources backing this joystick's signals
func (j *Joystick) Close() error {
	return j.close()
}

// GetJoysticks returns all known active joysticks.
func GetJoysticks() []*Joystick {
	return getJoysticks()
}

// WaitForJoysticks will regularly call GetJoysticks so to send signals
// on the output channel when new joysticks are connected to the system.
// Call `cancel' to close the channel and stop polling.
func WaitForJoysticks(pollRate time.Duration) (joyCh <-chan *Joystick, cancel func()) {
	ch := make(chan *Joystick)
	stop := make(chan struct{})
	go func() {
		lastJoysticks := getJoysticks()
		// Send all existing joysticks
		for _, j := range lastJoysticks {
			ch <- j
		}
		for {
			t := time.NewTicker(pollRate)
			select {
			case <-t.C:
			case <-stop:
				t.Stop()
				close(ch)
				close(stop)
				return
			}
			joys := getJoysticks()
		OUTER:
			for _, j := range joys {
				for _, j2 := range lastJoysticks {
					if j.id == j2.id {
						continue OUTER
					}
				}
				// j is new
				ch <- j
			}
			lastJoysticks = joys
		}
	}()
	cancel = func() {
		stop <- struct{}{}
	}
	return ch, cancel
}
