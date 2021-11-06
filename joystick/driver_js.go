package joystick

import (
	"reflect"
	"syscall/js"
	"errors"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/timing"
)

func osinit() error { 
	// TODO: listen to joystick connected and joystick disconnected? We'd still need to 
	// list from getGamepads every frame, it seems, to get new button presses.
	return nil
}

func newJoystick(gp js.Value, id uint32) *Joystick {
	return &Joystick{
		Handler:    event.DefaultBus,
		PollRate:   timing.FPSToFrameDelay(30),
		id:         id,
		osJoystick: newOsJoystick(gp),
	}
}

func newOsJoystick(gp js.Value) osJoystick {
	return osJoystick{
		cache: State{
			Buttons: make(map[string]bool),
		},
		newButtons: make(map[string]bool),
	}
}

func refreshGamepadState(j *Joystick, gp js.Value) {
	buttons := gp.Get("buttons")
	jsAxes := gp.Get("axes")
	if j.newJSState.axes == nil {
		j.newJSState.axes = make([]float64, jsAxes.Length())
	}
	for i := 0; i < jsAxes.Length(); i++ {
		j.newJSState.axes[i] = jsAxes.Index(i).Float()
	}
	if j.newJSState.buttons == nil {
		j.newJSState.buttons = make([]jsButton, buttons.Length())
	}
	for i := 0; i < buttons.Length(); i++ {
		jsBtn := buttons.Index(i)
		j.newJSState.buttons[i] = jsButton{
			value:   jsBtn.Get("value").Float(),
			pressed: jsBtn.Get("pressed").Bool(),
		}
	}
	j.newJSState.connected = gp.Get("connected").Bool()
	j.newJSState.mapping = gp.Get("mapping").String()
}

type jsGamepadState struct {
	axes      []float64
	buttons   []jsButton
	connected bool
	// osID      string
	// index     int
	mapping   string
}

type jsButton struct {
	value   float64
	//touched bool
	pressed bool
}

type osJoystick struct {
	cache State
	jsState      jsGamepadState
	newJSState      jsGamepadState
	newButtons map[string]bool
}

var (
	standardMappingButtons = []string{
		0: "A",
		1: "B",
		2: "X",
		3: "Y",
		4: "LeftShoulder",
		5: "RightShoulder",
		//6: LeftTrigger
		//7: RightTrigger
		8:  "Back",
		9:  "Start",
		10: "LeftStick",
		11: "RightStick",
		12: "Up",
		13: "Down",
		14: "Left",
		15: "Right",
	}
)

func (j *Joystick) prepare() error {
	return nil
}

func (j *Joystick) getState() (*State, error) {
	gp := js.Global().Get("navigator").Call("getGamepads").Index(int(j.id))
	if gp.IsNull() {
		return nil, errors.New("Joystick disconnected")
	}
	refreshGamepadState(j, gp)
	if !reflect.DeepEqual(j.newJSState, j.jsState) {
		j.jsState, j.newJSState = j.newJSState, j.jsState
		switch j.jsState.mapping {
		default:
			fallthrough
		case "standard":
			for i, btn := range j.jsState.buttons {
				if i >= len(standardMappingButtons) || standardMappingButtons[i] == "" {
					continue
				}
				j.newButtons[standardMappingButtons[i]] = btn.pressed
			}
			const int16scale = (2 << 14) - 1
			const uint8scale = (2 << 7) - 1
			j.cache.StickLX = int16(j.jsState.axes[0] * (int16scale))
			j.cache.StickLY = int16(j.jsState.axes[1]*(int16scale)) * -1
			j.cache.StickRX = int16(j.jsState.axes[2] * (int16scale))
			j.cache.StickRY = int16(j.jsState.axes[3]*(int16scale)) * -1
			j.cache.TriggerL = uint8(j.jsState.buttons[6].value * uint8scale)
			j.cache.TriggerR = uint8(j.jsState.buttons[7].value * uint8scale)
		}
		j.newButtons, j.cache.Buttons = j.cache.Buttons, j.newButtons
		j.cache.Frame++
	}
	s := new(State)
	*s = j.cache
	return s, nil
}

func (j *Joystick) vibrate(left, right uint16) error {
	return oakerr.UnsupportedPlatform{Operation: "joystick-vibrate"}
}

func (j *Joystick) close() error {
	return nil
}

func getJoysticks() []*Joystick {
	gamepads := js.Global().Get("navigator").Call("getGamepads")
	joysticks := make([]*Joystick, 0, gamepads.Length())
	for i := 0; i < gamepads.Length(); i++ {
		j := gamepads.Index(i)
		if !j.Truthy() {
			continue
		}
		joysticks = append(joysticks, newJoystick(j, uint32(i)))
	}

	return joysticks
}
