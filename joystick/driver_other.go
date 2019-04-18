// +build !windows,!linux,!darwin

package joystick

func osinit() {}

func (j *Joystick) prepare() error {
	return errors.New("OS not supported")
}

func (j *Joystick) getState() (*State, error) {
	return nil, return errors.New("OS not supported")
}

func (j *Joystick) vibrate(left, right uint16) error {
	return errors.New("OS not supported")
}

func (j *Joystick) close() error {
	return errors.New("OS not supported")
}

func getJoysticks() []*Joystick {
	return nil
}
