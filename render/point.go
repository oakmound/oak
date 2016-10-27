package render

type Point struct {
	X, Y float64
}

func (p *Point) ShiftX(x float64) {
	p.X += x
	SetDirty(p.X, p.Y)
}
func (p *Point) ShiftY(y float64) {
	p.Y += y
	SetDirty(p.X, p.Y)
}
func (p *Point) GetX() float64 {
	return p.X
}
func (p *Point) GetY() float64 {
	return p.Y
}

func (p *Point) SetPos(x, y float64) {
	p.X = x
	p.Y = y
	SetDirty(p.X, p.Y)
}

func (p *Point) AlwaysDirty() bool {
	return false
}
