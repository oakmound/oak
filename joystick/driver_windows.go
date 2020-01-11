package joystick

import (
	"sync"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/timing"
	"github.com/oakmound/w32"
)

func newJoystick(id uint32) *Joystick {
	return &Joystick{
		Handler:  event.DefaultBus,
		PollRate: timing.FPSToDuration(60),
		id:       id,
		osJoystick: osJoystick{
			wstate:    &w32.XInputState{},
			vibration: &w32.XInputVibration{},
		},
	}
}

type osJoystick struct {
	// Todo: mutex these values?
	wstate    *w32.XInputState
	vibration *w32.XInputVibration
}

// The windows driver currently uses the xinput api.
// We should consider providing alternatives.

var once sync.Once

func osinit() error {
	var err error
	once.Do(func() {
		err = w32.InitXInput()
	})
	return err
}

func (j *Joystick) prepare() error {
	return w32.XInputEnable(true)
}

type buttonName struct {
	name      string
	xinputVal uint16
}

var (
	chkButtons = []buttonName{
		{"Up", w32.XINPUT_GAMEPAD_DPAD_UP},
		{"Down", w32.XINPUT_GAMEPAD_DPAD_DOWN},
		{"Left", w32.XINPUT_GAMEPAD_DPAD_LEFT},
		{"Right", w32.XINPUT_GAMEPAD_DPAD_RIGHT},
		{"Start", w32.XINPUT_GAMEPAD_START},
		{"Back", w32.XINPUT_GAMEPAD_BACK},
		{"LeftStick", w32.XINPUT_GAMEPAD_LEFT_THUMB},
		{"RightStick", w32.XINPUT_GAMEPAD_RIGHT_THUMB},
		{"LeftShoulder", w32.XINPUT_GAMEPAD_LEFT_SHOULDER},
		{"RightShoulder", w32.XINPUT_GAMEPAD_RIGHT_SHOULDER},
		{"A", w32.XINPUT_GAMEPAD_A},
		{"B", w32.XINPUT_GAMEPAD_B},
		{"X", w32.XINPUT_GAMEPAD_X},
		{"Y", w32.XINPUT_GAMEPAD_Y},
	}
)

func (j *Joystick) getState() (*State, error) {
	err := w32.XInputGetState(j.id, j.wstate)
	if err != nil {
		return nil, err
	}
	// Convert windows state into os-regular state
	s := &State{
		Frame:    j.wstate.PacketNumber,
		StickLX:  j.wstate.Gamepad.ThumbLX,
		StickLY:  j.wstate.Gamepad.ThumbLY,
		StickRX:  j.wstate.Gamepad.ThumbRX,
		StickRY:  j.wstate.Gamepad.ThumbRY,
		TriggerL: j.wstate.Gamepad.LeftTrigger,
		TriggerR: j.wstate.Gamepad.RightTrigger,
		ID:       j.id,
		Buttons:  make(map[string]bool, len(chkButtons)),
	}

	for _, chk := range chkButtons {
		if j.wstate.Gamepad.Buttons&chk.xinputVal > 0 {
			s.Buttons[chk.name] = true
		} else {
			s.Buttons[chk.name] = false
		}
	}
	return s, nil
}

func (j *Joystick) vibrate(left, right uint16) error {
	j.vibration.LeftMotorSpeed = left
	j.vibration.RightMotorSpeed = right
	// Todo: wrap these errors?
	return w32.XInputSetState(j.id, j.vibration)
}

func (j *Joystick) close() error {
	// It seemingly makes sense to do this, but doing this disables
	// detection of future joysticks
	//return w32.XInputEnable(false)
	return nil
}

func getJoysticks() []*Joystick {
	// With xinput there are explicitly up to 4 controllers
	joys := make([]*Joystick, 0, 4)
	for i := 0; i < 4; i++ {
		err := w32.XInputGetState(uint32(i), &w32.XInputState{})
		if err == nil {
			joys = append(joys, newJoystick(uint32(i)))
		}
	}
	return joys
}
