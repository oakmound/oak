package render

import (
	"errors"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"math"
)

type Polygon struct {
	points                 []Point
	x, y                   float64
	r                      *image.RGBA
	layer                  int
	minX, maxX, minY, maxY float64
}

type Point struct {
	X, Y float64
}

func ScreenPolygon(points []Point, w, h int) (*Polygon, error) {
	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}

	minX, minY, maxX, maxY, _, _ := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		points,
		0,
		0,
		rgba,
		0,
		minX, minY, maxX, maxY,
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
		points,
		minX,
		minY,
		rgba,
		0,
		minX, minY, maxX, maxY,
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

func (pg *Polygon) FillInverse(c color.Color) {
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
		if (tp1.Y > y) != (tp2.Y > y) {
			if x < (tp2.X-tp1.X)*(y-tp1.Y)/(tp2.Y-tp1.Y)+tp1.X {
				contains = !contains
			}
		}
		j = i
	}
	return
}
