package radar

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v2/render"
)

// Point is a utility function for location
type Point struct {
	X, Y *float64
}

// Radar helps store and present information around interesting entities on a radar map
type Radar struct {
	render.LayeredPoint
	points        map[Point]color.Color
	center        Point
	width, height int
	r             *image.RGBA
	outline       *render.Sprite
	ratio         float64
}

var (
	centerColor = color.RGBA{255, 255, 0, 255}
)

// NewRadar creates a radar that will display at 0,0 with the given dimensions.
// The points given will be displayed on the radar relative to the center point,
// With the absolute distance reduced by the given ratio
func NewRadar(w, h int, points map[Point]color.Color, center Point, ratio float64) *Radar {
	r := new(Radar)
	r.LayeredPoint = render.NewLayeredPoint(0, 0, 0)
	r.points = points
	r.width = w
	r.height = h
	r.center = center
	r.r = image.NewRGBA(image.Rect(0, 0, w, h))
	r.outline = render.NewColorBox(w, h, color.RGBA{0, 0, 125, 125})
	r.ratio = ratio
	return r
}

// SetPos sets the position of the radar on the screen
func (r *Radar) SetPos(x, y float64) {
	r.LayeredPoint.SetPos(x, y)
	r.outline.SetPos(x, y)
}

// GetRGBA returns this radar's image
func (r *Radar) GetRGBA() *image.RGBA {
	return r.r
}

// Draw draws the radar, satisfying render.Renderable
func (r *Radar) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}

// DrawOffset draws the radar at a given offset
func (r *Radar) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Draw each point p in r.points
	// at r.X() + center.X() - p.X(), r.Y() + center.Y() - p.Y()
	// IF that value is < r.width/2, > -r.width/2, < r.height/2, > -r.height/2
	for p, c := range r.points {
		x := int((*p.X-*r.center.X)/r.ratio) + r.width/2
		y := int((*p.Y-*r.center.Y)/r.ratio) + r.height/2
		for x2 := x - 1; x2 < x+1; x2++ {
			for y2 := y - 1; y2 < y+1; y2++ {
				r.r.Set(x2, y2, c)
			}
		}
	}
	r.r.Set(r.width/2, r.height/2, centerColor)
	render.ShinyDraw(buff, r.r, int(xOff+r.X()), int(yOff+r.Y()))
	r.outline.DrawOffset(buff, xOff, yOff)
	r.r = image.NewRGBA(image.Rect(0, 0, r.width, r.height))
}

// AddPoint adds an additional point to the radar to be tracked
func (r *Radar) AddPoint(loc Point, c color.Color) {
	r.points[loc] = c
}
