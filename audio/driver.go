package audio

// A Driver defines the underlying interface that should be used for initializing PCM audio writers.
type Driver int

const (
	// DriverDefault indicates to this package to use a default driver based on the OS.
	// Currently, for windows the default is DirectSound and for unix the default is PulseAudio.
	DriverDefault Driver = iota
	DriverPulse
	DriverDirectSound
	DriverALSA
)

var driverNames = map[Driver]string{
	DriverPulse:       "pulseaudio",
	DriverDirectSound: "directsound",
	DriverDefault:     "default",
	DriverALSA:        "alsa",
}

func (d Driver) String() string {
	return driverNames[d]
}
