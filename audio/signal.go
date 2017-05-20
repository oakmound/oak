package audio

// A ChannelSignal is sent to an AudioChannel to indicate when they should
// attempt to play an audio sound
type ChannelSignal interface {
	GetIndex() int
	GetPos() (bool, float64, float64)
}

// Signal is a default ChannelSignal that just indicates what audio from a set
// of audio data should be played by int index
type Signal struct {
	Index int
}

// GetIndex returns Signal.Index
func (s Signal) GetIndex() int {
	return s.Index
}

// GetPos returns that a Signal is not positional
func (s Signal) GetPos() (bool, float64, float64) {
	return false, 0, 0
}

// PosSignal is a ChannelSignal compatible with Ears
type PosSignal struct {
	Signal
	X, Y float64
}

// NewPosSignal constructs a PosSignal
func NewPosSignal(index int, x, y float64) PosSignal {
	return PosSignal{
		Signal{index},
		x,
		y,
	}
}

// GetPos returns the floating points passed into a PosSignal
// as the origin of a sound to be heard
func (ps PosSignal) GetPos() (bool, float64, float64) {
	return true, ps.X, ps.Y
}
