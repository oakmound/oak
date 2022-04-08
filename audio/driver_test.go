package audio

import "testing"

func TestDriver_String(t *testing.T) {
	drivers := []Driver{
		DriverDefault,
		DriverDirectSound,
		DriverPulse,
		DriverALSA,
	}
	for _, d := range drivers {
		if d.String() == "" {
			t.Errorf("driver %d had no defined string", d)
		}
	}
}
