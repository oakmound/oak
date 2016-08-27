package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
)

type Moving struct {
	Solid
	moving
}

func (m *Moving) Init() event.CID {
	cID := event.NextID(m)
	m.CID = cID
	return cID
}

type moving struct {
	DX, DY, SpeedX, SpeedY float64
}

func (m *moving) GetDX() float64 {
	return m.DX
}
func (m *moving) GetDY() float64 {
	return m.DY
}
func (m *moving) SetDXY(x, y float64) {
	m.DX = x
	m.DY = y
}
func (m *moving) GetSpeedX() float64 {
	return m.SpeedX
}
func (m *moving) GetSpeedY() float64 {
	return m.SpeedY
}
func (m *moving) SetSpeedXY(x, y float64) {
	m.SpeedX = x
	m.SpeedY = y
}
