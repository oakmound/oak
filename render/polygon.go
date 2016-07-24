package render

import (
	"errors"
	geo "github.com/kellydunn/golang-geo"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"math"
)

type Polygon struct {
	gPolygon *geo.Polygon
	x, y     float64
	r        *image.RGBA
	layer    int
}

func ScreenPolygon(points []*geo.Point, w, h int) (*Polygon, error) {
	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}
	gPolygon := geo.NewPolygon(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		gPolygon,
		0,
		0,
		rgba,
		0,
	}, nil
}

func NewPolygon(points []*geo.Point) (*Polygon, error) {

	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}
	gPolygon := geo.NewPolygon(points)

	// Calculate the bounding rectangle of the polygon by
	// finding the maximum and minimum x and y values of the given points
	minX, minY, w, h := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		gPolygon,
		minX,
		minY,
		rgba,
		0,
	}, nil
}

func (pg *Polygon) UpdatePoints(points []*geo.Point) {
	pg.gPolygon = geo.NewPolygon(points)
}

func (pg *Polygon) Fill(c color.Color) {
	// Reset the rgba of the polygon
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			p := geo.NewPoint(float64(x), float64(y))
			if pg.gPolygon.Contains(p) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

func (pg *Polygon) FillInverse(c color.Color) {
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			p := geo.NewPoint(float64(x), float64(y))
			if !pg.gPolygon.Contains(p) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

func BoundingRect(points []*geo.Point) (minX, minY float64, w, h int) {
	minX = math.MaxFloat64
	minY = math.MaxFloat64
	maxX := minX * -1
	maxY := minY * -1
	for _, p := range points {
		lat := p.Lat()
		lng := p.Lng()
		if lat < minX {
			minX = lat
		}
		if lat > maxX {
			maxX = lat
		}
		if lng < minY {
			minY = lng
		}
		if lng > maxY {
			maxY = lng
		}
	}
	w = int(maxX - minX)
	h = int(maxY - minY)
	return
}

func (pg *Polygon) GetRGBA() *image.RGBA {
	return pg.r
}

func (pg *Polygon) Draw(buff screen.Buffer) {
	ShinyDraw(buff, pg.r, int(pg.x), int(pg.y))
}

func (pg *Polygon) ShiftX(x float64) {
	pg.x += x
}
func (pg *Polygon) ShiftY(y float64) {
	pg.y += y
}

func (pg *Polygon) GetLayer() int {
	return pg.layer
}

func (pg *Polygon) SetLayer(l int) {
	pg.layer = l
}

func (pg *Polygon) UnDraw() {
	pg.layer = -1
}

func (pg *Polygon) SetPos(x, y float64) {
	pg.x = x
	pg.y = y
}

// kellydunn/golang-geo
// we need to parallelize and duplicate this functionality
// Returns the points of the current Polygon.
// func (p *Polygon) Points() []*Point {
// 	return p.points
// }

// // Appends the passed in contour to the current Polygon.
// func (p *Polygon) Add(point *Point) {
// 	p.points = append(p.points, point)
// }

// // Returns whether or not the polygon is closed.
// // TODO:  This can obviously be improved, but for now,
// //        this should be sufficient for detecting if points
// //        are contained using the raycast algorithm.
// func (p *Polygon) IsClosed() bool {
// 	if len(p.points) < 3 {
// 		return false
// 	}

// 	return true
// }

// // Returns whether or not the current Polygon contains the passed in Point.
// func (p *Polygon) Contains(point *Point) bool {
// 	if !p.IsClosed() {
// 		return false
// 	}

// 	start := len(p.points) - 1
// 	end := 0

// 	contains := p.intersectsWithRaycast(point, p.points[start], p.points[end])

// 	for i := 1; i < len(p.points); i++ {
// 		if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
// 			contains = !contains
// 		}
// 	}

// 	return contains
// }

// // Using the raycast algorithm, this returns whether or not the passed in point
// // Intersects with the edge drawn by the passed in start and end points.
// // Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
// func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
// 	// Always ensure that the the first point
// 	// has a y coordinate that is less than the second point
// 	if start.lng > end.lng {

// 		// Switch the points if otherwise.
// 		start, end = end, start

// 	}

// 	// Move the point's y coordinate
// 	// outside of the bounds of the testing region
// 	// so we can start drawing a ray
// 	for point.lng == start.lng || point.lng == end.lng {
// 		newLng := math.Nextafter(point.lng, math.Inf(1))
// 		point = NewPoint(point.lat, newLng)
// 	}

// 	// If we are outside of the polygon, indicate so.
// 	if point.lng < start.lng || point.lng > end.lng {
// 		return false
// 	}

// 	if start.lat > end.lat {
// 		if point.lat > start.lat {
// 			return false
// 		}
// 		if point.lat < end.lat {
// 			return true
// 		}

// 	} else {
// 		if point.lat > end.lat {
// 			return false
// 		}
// 		if point.lat < start.lat {
// 			return true
// 		}
// 	}

// 	raySlope := (point.lng - start.lng) / (point.lat - start.lat)
// 	diagSlope := (end.lng - start.lng) / (end.lat - start.lat)

// 	return raySlope >= diagSlope
// }
