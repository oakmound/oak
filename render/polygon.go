package render

import (
	"errors"
	"image"
	"image/color"
	"math"
)

type Rect struct {
	minX, maxX, minY, maxY float64
}

type Polygon struct {
	Sprite
	Rect
	points []Point
}

func ScreenPolygon(points []Point, w, h int) (*Polygon, error) {
	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}

	minX, minY, maxX, maxY, _, _ := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		Sprite: Sprite{
			LayeredPoint: LayeredPoint{
				Point: Point{
					X: 0.0,
					Y: 0.0,
				},
			},
			r: rgba,
		},
		Rect: Rect{
			minX: minX,
			minY: minY,
			maxX: maxX,
			maxY: maxY,
		},
		points: points,
	}, nil
}

func NewPolygon(points []Point) (*Polygon, error) {

	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}

	// Calculate the bounding rectangle of the polygon by
	// finding the maximum and minimum x and y values of the given points
	minX, minY, maxX, maxY, w, h := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		Sprite: Sprite{
			LayeredPoint: LayeredPoint{
				Point: Point{
					X: 0.0,
					Y: 0.0,
				},
			},
			r: rgba,
		},
		Rect: Rect{
			minX: minX,
			minY: minY,
			maxX: maxX,
			maxY: maxY,
		},
		points: points,
	}, nil
}

func (pg *Polygon) UpdatePoints(points []Point) {
	pg.points = points
	pg.minX, pg.minY, pg.maxX, pg.maxY, _, _ = BoundingRect(points)
}

func (pg *Polygon) Fill(c color.Color) {
	// Reset the rgba of the polygon
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if !pg.Contains(float64(x), float64(y)) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

func (pg *Polygon) GetOutline(c color.Color) *CompositeSlice {
	sl := new(CompositeSlice)
	j := len(pg.points) - 1
	for i, p2 := range pg.points {
		p1 := pg.points[j]
		sl.Append(NewLine(p1.X, p1.Y, p2.X, p2.Y, c))
		j = i
	}
	return sl
}

func (pg *Polygon) FillInverse(c color.Color) {
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if !pg.ConvexContains(float64(x), float64(y)) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

func BoundingRect(points []Point) (minX, minY, maxX, maxY float64, w, h int) {
	minX = math.MaxFloat64
	minY = math.MaxFloat64
	maxX = minX * -1
	maxY = minY * -1
	for _, p := range points {
		x := p.X
		y := p.Y
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	w = int(maxX - minX)
	h = int(maxY - minY)
	return
}

// Returns whether or not the current Polygon contains the passed in Point.
// Still need to parallelize
func (pg *Polygon) Contains(x, y float64) (contains bool) {

	if x < pg.minX || x > pg.maxX || y < pg.minY || y > pg.maxY {
		return
	}

	j := len(pg.points) - 1
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[j]
		if (tp1.Y > y) != (tp2.Y > y) { // Three comparisons
			if x < (tp2.X-tp1.X)*(y-tp1.Y)/(tp2.Y-tp1.Y)+tp1.X { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		j = i
	}
	return
}

func (pg *Polygon) WrappingContains(x, y float64) bool {

	if x < pg.minX || x > pg.maxX || y < pg.minY || y > pg.maxY {
		return false
	}

	wn := 0

	j := len(pg.points) - 1
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[j]
		if tp1.Y <= y && tp2.Y > y && isLeft(tp1, tp2, x, y) > 0 { // Three comparison, Five add/sub, Two mult/div
			wn++
		}
		if tp2.Y >= y && isLeft(tp1, tp2, x, y) < 0 { // Two Comparison, Five add/sub, Two mult/div
			wn--
		}
		j = i
	}
	return wn == 0
}

func (pg *Polygon) ConvexContains(x, y float64) bool {

	if x < pg.minX || x > pg.maxX || y < pg.minY || y > pg.maxY {
		return false
	}

	prev := 0
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[(i+1)%len(pg.points)]
		tp3 := vSub(tp2, tp1)
		tp4 := vSub(Point{x, y}, tp1)
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

func getSide(a, b Point) int {
	x := a.X*b.Y - a.Y*b.X
	if x == 0 {
		return 0
	} else if x < 1 {
		return -1
	} else {
		return 1
	}
}

func vSub(a, b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}

func isLeft(p1, p2 Point, x, y float64) float64 {
	return (p1.X-x)*(p2.Y-y) - (p2.X-x)*(p1.Y-y)
}
