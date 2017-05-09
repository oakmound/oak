package audio

type ChannelSignal interface {
	GetIndex() int
	GetPos() (bool, float64, float64)
}

type Signal struct {
	Index int
}

func (s Signal) GetIndex() int {
	return s.Index
}

func (s Signal) GetPos() (bool, float64, float64) {
	return false, 0, 0
}

type PosSignal struct {
	Index int
	X, Y  float64
}

func (ps PosSignal) GetIndex() int {
	return ps.Index
}

func (ps PosSignal) GetPos() (bool, float64, float64) {
	return true, ps.X, ps.Y
}
