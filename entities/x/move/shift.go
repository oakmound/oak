package move

// ShiftX will ShiftX on the vector of the mover and
// set the renderable and space positions to that of the updated vector.
func ShiftX(mvr Mover, x float64) {
	vec := mvr.Vec()
	vec.ShiftX(x)
	mvr.GetRenderable().SetPos(vec.X(), vec.Y())
	sp := mvr.GetSpace()
	sp.Update(vec.X(), vec.Y(), sp.GetW(), sp.GetH())
}

// ShiftY will ShiftY on the vector of the mover and
// set the renderable and space positions to that of the updated vector.
func ShiftY(mvr Mover, y float64) {
	vec := mvr.Vec()
	vec.ShiftY(y)
	mvr.GetRenderable().SetPos(vec.X(), vec.Y())
	sp := mvr.GetSpace()
	sp.Update(vec.X(), vec.Y(), sp.GetW(), sp.GetH())
}
