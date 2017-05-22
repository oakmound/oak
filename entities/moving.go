package entities

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Moving struct {
	Solid
	moving
}

func NewMoving(x, y, w, h float64, r render.Renderable, cid event.CID) Moving {
	return Moving{
		Solid:  NewSolid(x, y, w, h, r, cid),
		moving: moving{},
	}
}

func (m *Moving) Init() event.CID {
	cID := event.NextID(m)
	m.CID = cID
	return cID
}

func (m *Moving) String() string {
	st := "Moving: \n{"
	st += m.Solid.String()
	st += "} \n" + m.moving.String()
	return st
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

func (m *moving) String() string {
	dx := strconv.FormatFloat(m.DX, 'f', 2, 32)
	dy := strconv.FormatFloat(m.DY, 'f', 2, 32)
	sx := strconv.FormatFloat(m.SpeedX, 'f', 2, 32)
	sy := strconv.FormatFloat(m.SpeedY, 'f', 2, 32)
	return "DX: " + dx + ", DY: " + dy + ", SX: " + sx + ", SY: " + sy
}
