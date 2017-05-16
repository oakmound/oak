package render

type Point struct {
	X, Y float64
}

func (p *Point) ShiftX(x float64) {
	p.X += x
}
func (p *Point) ShiftY(y float64) {
	p.Y += y
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
}

func (p *Point) GetPos() (float64, float64) {
	return p.X, p.Y
}

func (p *Point) GetDims() (int, int) {
	return 6, 6
}
