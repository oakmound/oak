package render

import (
	"errors"
	"image"
	"image/color"
	"math"

	"bitbucket.org/oakmoundstudio/oak/physics"
)

type Rect struct {
	MinX, MaxX, MinY, MaxY float64
}

type Polygon struct {
	*Sprite
	Rect
	points []physics.Vector
}

func ScreenPolygon(points []physics.Vector, w, h int) (*Polygon, error) {
	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}

	MinX, MinY, MaxX, MaxY, _, _ := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		Sprite: NewSprite(0, 0, rgba),
		Rect: Rect{
			MinX: MinX,
			MinY: MinY,
			MaxX: MaxX,
			MaxY: MaxY,
		},
		points: points,
	}, nil
}

func NewPolygon(points []physics.Vector) (*Polygon, error) {

	if len(points) < 3 {
		return nil, errors.New("Please give at least three points to NewPolygon calls")
	}

	// Calculate the bounding rectangle of the polygon by
	// finding the maximum and minimum x and y values of the given points
	MinX, MinY, MaxX, MaxY, w, h := BoundingRect(points)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	return &Polygon{
		Sprite: NewSprite(MinX, MinY, rgba),
		Rect: Rect{
			MinX: MinX,
			MinY: MinY,
			MaxX: MaxX,
			MaxY: MaxY,
		},
		points: points,
	}, nil
}

func (pg *Polygon) UpdatePoints(points []physics.Vector) {
	pg.points = points
	pg.MinX, pg.MinY, pg.MaxX, pg.MaxY, _, _ = BoundingRect(points)
}

func (pg *Polygon) Fill(c color.Color) {
	// Reset the rgba of the polygon
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	minx := pg.Rect.MinX
	miny := pg.Rect.MinY
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if pg.Contains(float64(x)+minx, float64(y)+miny) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

func (pg *Polygon) GetOutline(c color.Color) *Composite {
	sl := NewComposite([]Modifiable{})
	j := len(pg.points) - 1
	for i, p2 := range pg.points {
		p1 := pg.points[j]
		MinX := math.Min(p1.X(), p2.X())
		MinY := math.Min(p1.Y(), p2.Y())
		sl.AppendOffset(NewLine(p1.X(), p1.Y(), p2.X(), p2.Y(), c), physics.NewVector(MinX, MinY))
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

func BoundingRect(points []physics.Vector) (MinX, MinY, MaxX, MaxY float64, w, h int) {
	MinX = math.MaxFloat64
	MinY = math.MaxFloat64
	MaxX = MinX * -1
	MaxY = MinY * -1
	for _, p := range points {
		x := p.X()
		y := p.Y()
		if x < MinX {
			MinX = x
		}
		if x > MaxX {
			MaxX = x
		}
		if y < MinY {
			MinY = y
		}
		if y > MaxY {
			MaxY = y
		}
	}
	w = int(MaxX - MinX)
	h = int(MaxY - MinY)
	return
}

// Returns whether or not the current Polygon contains the passed in Point.
// Still need to parallelize
func (pg *Polygon) Contains(x, y float64) (contains bool) {

	if x < pg.MinX || x > pg.MaxX || y < pg.MinY || y > pg.MaxY {
		return
	}

	j := len(pg.points) - 1
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[j]
		if (tp1.Y() > y) != (tp2.Y() > y) { // Three comparisons
			if x < (tp2.X()-tp1.X())*(y-tp1.Y())/(tp2.Y()-tp1.Y())+tp1.X() { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		j = i
	}
	return
}

func (pg *Polygon) WrappingContains(x, y float64) bool {

	if x < pg.MinX || x > pg.MaxX || y < pg.MinY || y > pg.MaxY {
		return false
	}

	wn := 0

	j := len(pg.points) - 1
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[j]
		if tp1.Y() <= y && tp2.Y() > y && isLeft(tp1, tp2, x, y) > 0 { // Three comparison, Five add/sub, Two mult/div
			wn++
		}
		if tp2.Y() >= y && isLeft(tp1, tp2, x, y) < 0 { // Two Comparison, Five add/sub, Two mult/div
			wn--
		}
		j = i
	}
	return wn == 0
}

func (pg *Polygon) ConvexContains(x, y float64) bool {

	if x < pg.MinX || x > pg.MaxX || y < pg.MinY || y > pg.MaxY {
		return false
	}

	prev := 0
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[(i+1)%len(pg.points)]
		tp3 := vSub(tp2, tp1)
		tp4 := vSub(physics.NewVector(x, y), tp1)
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

func getSide(a, b physics.Vector) int {
	x := a.X()*b.Y() - a.Y()*b.X()
	if x == 0 {
		return 0
	} else if x < 1 {
		return -1
	} else {
		return 1
	}
}

func vSub(a, b physics.Vector) physics.Vector {
	return physics.NewVector(a.X()-b.X(), a.Y()-b.Y())
}

func isLeft(p1, p2 physics.Vector, x, y float64) float64 {
	return (p1.X()-x)*(p2.Y()-y) - (p2.X()-x)*(p1.Y()-y)
}
