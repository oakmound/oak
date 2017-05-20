package collision

import (
	"image/color"
	"math"
	"time"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// RayCast returns the set of points where a line
// from x,y going at a certain angle, for a certain length, intersects
// with existing rectangles in the rtree.
// It converts the ray into a series of points which are themselves
// used to check collision at a miniscule width and height.
func RayCast(x, y, degrees, length float64) []CollisionPoint {
	results := []CollisionPoint{}
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
				results = append(results, CollisionPoint{nx, x, y})
			}
		}
		x += c
		y += s
	}
	return results
}

// RatCastSingle acts as RayCast, but it returns only the first collision
// that the generated ray intersects, ignoring entities
// in the given invalidIDs list.
// Example Use case: shooting a bullet, hitting the first thing that isn't yourself.
func RayCastSingle(x, y, degrees, length float64, invalidIDS []event.CID) CollisionPoint {

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
			return CollisionPoint{nx, x, y}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}

func RayCastSingleLabel(x, y, degrees, length float64, label int) CollisionPoint {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			if nx.Label == label {
				return CollisionPoint{nx, x, y}
			}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}

func RayCastSingleLabels(x, y, degrees, length float64, labels ...int) CollisionPoint {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for _, label := range labels {
				if nx.Label == label {
					return CollisionPoint{nx, x, y}
				}
			}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}

func RayCastSingleIgnoreLabels(x, y, degrees, length float64, labels ...int) CollisionPoint {
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
			return CollisionPoint{nx, x, y}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}

func RayCastSingleIgnore(x, y, degrees, length float64, invalidIDS []event.CID, labels ...int) CollisionPoint {
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
			return CollisionPoint{nx, x, y}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}

// ConeCast advances COUNTER-CLOCKWISE
func ConeCast(x, y, angle, angleWidth, length float64) (points []CollisionPoint) {
	da := angleWidth / 10
	for a := angle; a < angle+angleWidth; a += da {
		cps := RayCast(x, y, a, length)
		if len(cps) > 0 {
			points = append(points, cps...)
		}
	}
	return
}

func ConeCastSingle(x, y, angle, angleWidth, length float64, invalidIDS []event.CID) (points []CollisionPoint) {
	da := angleWidth / 10
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingle(x, y, a, length, invalidIDS)
		if cp.Zone != nil {
			points = append(points, cp)
			//sweep := render.NewLine(x, y, cp.X, cp.Y, color.RGBA{255, 255, 255, 255})
			//render.Draw(sweep, 5000)
			//render.UndrawAfter(sweep, 50*time.Millisecond)
		}
	}
	return
}

func ConeCastSingleLabel(x, y, angle, angleWidth, length float64, label int) (points []CollisionPoint) {
	da := angleWidth / 10
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingleLabel(x, y, a, length, label)
		if cp.Zone != nil {
			points = append(points, cp)
		}
	}
	return
}

func ConeCastSingleLabels(x, y, angle, angleWidth, length float64, labels ...int) (points []CollisionPoint) {
	da := angleWidth / 10
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingleLabels(x, y, a, length, labels...)
		if cp.Zone != nil {
			l := render.NewLine(x, y, cp.X, cp.Y, color.RGBA{0, 0, 255, 255})
			l.SetLayer(60000)
			render.DrawForTime(l, 2, time.Millisecond*50)
			points = append(points, cp)
		}
	}
	return
}

func ConeCastSingleLabelsCt(x, y, angle, angleWidth, rays, length float64, labels ...int) (points []CollisionPoint) {
	da := angleWidth / rays
	for a := angle; a < angle+angleWidth; a += da {
		cp := RayCastSingleLabels(x, y, a, length, labels...)
		if cp.Zone != nil {
			l := render.NewLine(x, y, cp.X, cp.Y, color.RGBA{0, 0, 255, 255})
			l.SetLayer(60000)
			render.DrawForTime(l, 2, time.Millisecond*50)
			points = append(points, cp)
		}
	}
	return
}
