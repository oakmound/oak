package joystick

import (
	"errors"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/timing"

	"encoding/binary"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/oakmound/libudev"
	"github.com/oakmound/libudev/types"
)

// This has all been tested with wired xbox 360 controllers.
// Todo: get more controllers, test with more controllers.

func newJoystick(devName string, id uint32) *Joystick {
	return &Joystick{
		Handler:  event.DefaultBus,
		PollRate: timing.FPSToDuration(60),
		id:       id,
		osJoystick: osJoystick{
			cache: State{
				Buttons: make(map[string]bool),
			},
			devName: devName,
			quit:    make(chan struct{}),
		},
	}
}

type osJoystick struct {
	devName string
	fh      *os.File
	cache   State
	sync.Mutex
	quit         chan struct{}
	disconnected bool
}

func osinit() error {
	return nil
}

type jevent struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

const (
	axisType   = 2
	buttonType = 1
)

var (
	buttons = []string{
		0: "A",
		1: "B",
		2: "X",
		3: "Y",
		4: "LeftShoulder",
		5: "RightShoulder",
		6: "Back",
		7: "Start",
		// 8 is the "Xbox" button in the center
		9:  "LeftStick",
		10: "RightStick",
	}
)

func (j *Joystick) prepare() error {
	var err error
	j.fh, err = os.Open(j.devName)
	if err == nil {
		go func(j *Joystick) {
			// Read events continually
			e := &jevent{}
			for {
				select {
				case <-j.quit:
					return
				default:
				}
				err := binary.Read(j.fh, binary.LittleEndian, e)
				if err != nil {
					j.disconnected = true
					return
				}
				j.Lock()
				switch e.Type {
				case axisType:
					switch e.Number {
					case 0:
						j.cache.StickLX = e.Value
					case 1:
						j.cache.StickLY = e.Value * -1
					case 2:
						// The controller offers int16 fidelity of the
						// triggers. We're lowering it to Xinput's uint8
						// Todo: Flip that around?
						j.cache.TriggerL = uint8(uint16(e.Value) / 16)
					case 3:
						j.cache.StickRX = e.Value
					case 4:
						j.cache.StickRY = e.Value * -1
					case 5:
						j.cache.TriggerR = uint8(uint16(e.Value) / 16)
					case 6:
						if e.Value < 0 {
							j.cache.Buttons["Left"] = true
							j.cache.Buttons["Right"] = false
						} else if e.Value > 0 {
							j.cache.Buttons["Right"] = true
							j.cache.Buttons["Left"] = false
						} else {
							j.cache.Buttons["Right"] = false
							j.cache.Buttons["Left"] = false
						}
					case 7:
						if e.Value < 0 {
							j.cache.Buttons["Up"] = true
							j.cache.Buttons["Down"] = false
						} else if e.Value > 0 {
							j.cache.Buttons["Down"] = true
							j.cache.Buttons["Up"] = false
						} else {
							j.cache.Buttons["Down"] = false
							j.cache.Buttons["Up"] = false
						}
					}
				case buttonType:
					j.cache.Buttons[buttons[e.Number]] = (e.Value == 1)
				}
				// No mutex here could cause a frame delay on inputs
				j.cache.Frame = e.Time
				j.Unlock()
			}
		}(j)
	}
	return err
}

func (j *Joystick) getState() (*State, error) {
	if j.disconnected {
		return nil, errors.New("Joystick disconnected")
	}
	s := new(State)
	*s = j.cache
	s.Buttons = make(map[string]bool)
	j.Lock()
	for k, b := range j.cache.Buttons {
		s.Buttons[k] = b
	}
	j.Unlock()
	return s, nil
}

func (j *Joystick) vibrate(left, right uint16) error {
	return errors.New("Vibration not supported")
}

func (j *Joystick) close() error {
	go func() {
		j.quit <- struct{}{}
	}()
	return j.fh.Close()
}

func getJoysticks() []*Joystick {
	sc := libudev.NewScanner()
	err, dvs := sc.ScanDevices()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Joysticks contain "js%d"
	rgx, err := regexp.Compile("js(\\d)+")
	if err != nil {
		dlog.Error(err)
		return nil
	}

	filtered := []*types.Device{}

	for _, d := range dvs {
		// Find joysticks
		if !rgx.MatchString(d.Devpath) {
			continue
		}
		// Ignore mice
		if v, ok := d.Env["ID_INPUT_MOUSE"]; ok && v == "1" {
			continue
		}
		// Todo: what else do we ignore?
		filtered = append(filtered, d)
	}

	joys := make([]*Joystick, len(filtered))
	for i, f := range filtered {
		var id uint32 = math.MaxUint32
		matches := rgx.FindStringSubmatch(f.Devpath)
		if len(matches) > 1 {
			idint, err := strconv.Atoi(matches[1])
			id = uint32(idint)
			dlog.ErrorCheck(err)
		}
		joys[i] = newJoystick(filepath.Join("/", "dev", f.Env["DEVNAME"]), id)
	}
	return joys
}
