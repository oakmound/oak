package collision

import (
	"math"

	"bitbucket.org/oakmoundstudio/oak/event"
)

// RayCast returns the set of points where a line
// from x,y going at a certain angle, for a certain length, intersects
// with existing rectangles in the rtree.
// It converts the ray into a series of points which are themselves
// used to check collision at a miniscule width and height.
func RayCast(x, y, degrees, length float64) []Point {
	results := []Point{}
	resultHash := make(map[*Space]bool)

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)

		next := rt.SearchIntersect(loc)

		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			if _, ok := resultHash[nx]; !ok {
				resultHash[nx] = true
				results = append(results, NewPoint(nx, x, y))
			}
		}
		x += c
		y += s
	}
	return results
}

// RayCastSingle acts as RayCast, but it returns only the first collision
// that the generated ray intersects, ignoring entities
// in the given invalidIDs list.
// Example Use case: shooting a bullet, hitting the first thing that isn't yourself.
func RayCastSingle(x, y, degrees, length float64, invalidIDS []event.CID) Point {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
	output:
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for e := 0; e < len(invalidIDS); e++ {
				if nx.CID == invalidIDS[e] {
					continue output
				}
			}
			return NewPoint(nx, x, y)
		}
		x += c
		y += s

	}
	return NilPoint()
}

// RayCastSingleLabels acts like RayCastSingle, but only returns elements
// that match one of the input labels
func RayCastSingleLabels(x, y, degrees, length float64, labels ...Label) Point {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for _, label := range labels {
				if nx.Label == label {
					return NewPoint(nx, x, y)
				}
			}
		}
		x += c
		y += s

	}
	return NilPoint()
}

// RayCastSingleIgnoreLabels is the opposite of Labels, in that it will return
// the first collision point that is not contained in the set of ignore labels
func RayCastSingleIgnoreLabels(x, y, degrees, length float64, labels ...Label) Point {
	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
	output:
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for _, label := range labels {
				if nx.Label == label {
					continue output
				}
			}
			return NewPoint(nx, x, y)
		}
		x += c
		y += s

	}
	return NilPoint()
}

// RayCastSingleIgnore is just like ignore labels but also ignores certain
// caller ids
func RayCastSingleIgnore(x, y, degrees, length float64, invalidIDS []event.CID, labels ...Label) Point {
	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
	output:
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for _, label := range labels {
				if nx.Label == label {
					continue output
				}
			}
			for e := 0; e < len(invalidIDS); e++ {
				if nx.CID == invalidIDS[e] {
					continue output
				}
			}
			return NewPoint(nx, x, y)
		}
		x += c
		y += s

	}
	return NilPoint()
}

// ConeCast repeatedly calls RayCast in a cone shape
// ConeCast advances COUNTER-CLOCKWISE
func ConeCast(x, y, angle, angleWidth, rays, length float64) (points []Point) {
	da := angleWidth / rays
	for a := angle; a < angle+angleWidth; a += da {
		cps := RayCast(x, y, a, length)
		if len(cps) > 0 {
			points = append(points, cps...)
		}
	}
	return
}

// ConeCastSingle repeatedly calls RayCastSignle in a cone shape
func ConeCastSingle(x, y, angle, angleWidth, rays, length float64, invalidIDS []event.CID) (points []Point) {
	da := angleWidth / rays
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingle(x, y, a, length, invalidIDS)
		if cp.Zone != nil {
			points = append(points, cp)
			//render.NewLine(x, y, cp.X(), cp.Y(), color.RGBA{255, 255, 255, 255})
			//render.Draw(sweep, 5000)
			//render.UndrawAfter(sweep, 50*time.Millisecond)
		}
	}
	return
}

// ConeCastSingleLabels repeatedly calls RayCastSingleLabels in a cone shape
func ConeCastSingleLabels(x, y, angle, angleWidth, rays, length float64, labels ...Label) (points []Point) {
	da := angleWidth / rays
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingleLabels(x, y, a, length, labels...)
		if cp.Zone != nil {
			//render.NewLine(x, y, cp.X(), cp.Y(), color.RGBA{0, 0, 255, 255})
			//l.SetLayer(60000)
			//render.DrawForTime(l, 2, time.Millisecond*50)
			points = append(points, cp)
		}
	}
	return
}
