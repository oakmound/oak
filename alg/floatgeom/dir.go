package floatgeom

// Dir2 is a helper type for representing points as directions
type Dir2 Point2

// Dir2 values
var (
	Up        = Dir2(Point2{0, -1})
	Down      = Dir2(Point2{0, 1})
	Left      = Dir2(Point2{-1, 0})
	Right     = Dir2(Point2{1, 0})
	UpRight   = Up.And(Right)
	DownRight = Down.And(Right)
	DownLeft  = Down.And(Left)
	UpLeft    = Up.And(Left)
)

// And combines two directions
func (d Dir2) And(d2 Dir2) Dir2 {
	return Dir2(Point2(d).Add(Point2(d2)))
}

// X retrieves the horizontal component of a Dir2
func (d Dir2) X() float64 {
	return Point2(d).X()
}

// Y retrieves the vertical component for a Dir2
func (d Dir2) Y() float64 {
	return Point2(d).Y()
}
