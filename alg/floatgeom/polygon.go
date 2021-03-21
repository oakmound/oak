package floatgeom

// A Polygon2 is a series of points in 2D space.
type Polygon2 struct {
	// Bounding is a cached bounding box calculated from the input points
	// It is exported for convenience, but should be modified with care
	Bounding Rect2
	// The component points of the polygon. If modified, Bounding should
	// be updated with NewBoundingRect2.
	Points []Point2
}

// NewPolygon2 is a helper method to construct a valid polygon. Polygons
// cannot contain less than 3 points.
func NewPolygon2(p1, p2, p3 Point2, pn ...Point2) Polygon2 {
	pts := append([]Point2{p1, p2, p3}, pn...)
	bounding := NewBoundingRect2(pts...)
	return Polygon2{
		Bounding: bounding,
		Points:   pts,
	}
}

// Contains returns whether or not the current Polygon contains the passed in Point.
// If it is known that the polygon is convex, ConvexContains should be preferred for
// performance.
func (pg Polygon2) Contains(x, y float64) (contains bool) {

	if !pg.Bounding.Contains(Point2{x, y}) {
		return
	}

	j := len(pg.Points) - 1
	for i := 0; i < len(pg.Points); i++ {
		tp1 := pg.Points[i]
		tp2 := pg.Points[j]
		if (tp1.Y() > y) != (tp2.Y() > y) { // Three comparisons
			if x < (tp2.X()-tp1.X())*(y-tp1.Y())/(tp2.Y()-tp1.Y())+tp1.X() { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		j = i
	}
	return
}

// ConvexContains returns whether the given point is contained by the input polygon.
// It assumes the polygon is convex.
func (pg Polygon2) ConvexContains(x, y float64) bool {

	p := Point2{x, y}

	if !pg.Bounding.Contains(p) {
		return false
	}

	prev := 0
	for i := 0; i < len(pg.Points); i++ {
		tp1 := pg.Points[i]
		tp2 := pg.Points[(i+1)%len(pg.Points)]
		tp3 := tp2.Sub(tp1)
		tp4 := p.Sub(tp1)
		cur := getSide(tp3, tp4)
		if cur == 0 {
			return false
		} else if prev == 0 {
		} else if prev != cur {
			return false
		}
		prev = cur
	}
	return true
}

// TODO: rename this to its real math name, export it
func getSide(a, b Point2) int {
	x := a.X()*b.Y() - a.Y()*b.X()
	if x == 0 {
		return 0
	} else if x < 1 {
		return -1
	} else {
		return 1
	}
}
