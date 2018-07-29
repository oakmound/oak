package floatgeom

//Tri3 is a triangle of Point3s
type Tri3 [3]Point3

// Barycentric finds the barycentric coordinates of the given x,y cartesian
// coordinates within this triangle. If the point (x,y) is outside of the
// triangle, one of the output values will be negative.
// Credit goes to github.com/yellingintothefan for their work in gel
func (t Tri3) Barycentric(x, y float64) Point3 {
	p := Point3{x, y, 0.0}
	v0 := t[1].Sub(t[0])
	v1 := t[2].Sub(t[0])
	v2 := p.Sub(t[0])
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	v := (d11*d20 - d01*d21) / (d00*d11 - d01*d01)
	w := (d00*d21 - d01*d20) / (d00*d11 - d01*d01)
	u := 1.0 - v - w
	return Point3{v, w, u}
}

// Normal calculates the surface normal of a triangle
func (t Tri3) Normal() Point3 {
	u := t[1].Sub(t[0])
	v := t[2].Sub(t[0])
	// Check that the triangle is defined in a clockwise fashion.

	return u.Cross(v).Normalize()
}
