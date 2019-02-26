package radar

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/render"
)

type RadarPoint struct {
	X, Y *float64
}

type Radar struct {
	render.LayeredPoint
	points        map[RadarPoint]color.Color
	center        RadarPoint
	width, height int
	r             *image.RGBA
	outline       *render.Sprite
}

const (
	ratio = 10.0
)

var (
	centerColor = color.RGBA{255, 255, 0, 255}
)

/* Sets up the radar display */
func NewRadar(w, h int, points map[RadarPoint]color.Color, center RadarPoint) *Radar {
	r := new(Radar)
	r.LayeredPoint = render.NewLayeredPoint(0, 0, 0)
	r.points = points
	r.width = w
	r.height = h
	r.center = center
	r.r = image.NewRGBA(image.Rect(0, 0, w, h))
	r.outline = render.NewColorBox(400, 400, color.RGBA{0, 0, 200, 0})
	return r
}

func (r *Radar) SetPos(x, y float64) {
	r.LayeredPoint.SetPos(x, y)
	r.outline.SetPos(x, y)
}

func (r *Radar) GetRGBA() *image.RGBA {
	return r.r
}

func (r *Radar) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}

func (r *Radar) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Draw each point p in r.points
	// at r.X() + center.X() - p.X(), r.Y() + center.Y() - p.Y()
	// IF that value is < r.width/2, > -r.width/2, < r.height/2, > -r.height/2
	for p, c := range r.points {
		x := int((*p.X-*r.center.X)/ratio) + r.width/2
		y := int((*p.Y-*r.center.Y)/ratio) + r.height/2
		r.r.Set(x, y, c)
	}
	r.r.Set(r.width/2, r.height/2, centerColor)
	render.ShinyDraw(buff, r.r, int(xOff+r.X()), int(yOff+r.Y()))
	r.outline.DrawOffset(buff, xOff, yOff)
	r.r = image.NewRGBA(image.Rect(0, 0, r.width, r.height))
}

func (r *Radar) AddPoint(loc RadarPoint, c color.Color) {
	r.points[loc] = c
}
