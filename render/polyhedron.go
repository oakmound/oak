package render

import (
	"image"
	"image/color"
	"math"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-compgeo/dcel"
)

type Polyhedron struct {
	Sprite
	dcel.DCEL
	Center physics.Vector
}

func NewCuboid(x, y, z, w, h, d float64) *Polyhedron {
	px := x
	py := y
	dc := dcel.DCEL{}
	// Is there a smart way to loop through this?
	// verts
	x = w
	y = h
	z = -d
	dc.Vertices = make([]dcel.Point, 8)
	dc.HalfEdges = make([]dcel.Edge, 24)
	dc.OutEdges = make([]*dcel.Edge, 8)
	for i := 0; i < 8; i++ {
		// Set coordinates of this vertex
		if x == 0 {
			x = w
		} else {
			x = 0
		}
		if i%2 == 0 {
			if y == 0 {
				y = h
			} else {
				y = 0
			}
		}
		if i%4 == 0 {
			z += d
		}

		dc.Vertices[i] = dcel.Point{x, y, z}
	}
	corners := []int{0, 3, 5, 6}
	// These edges, except for the ones
	// at a corner's index, are those a
	// corner's edges' twins start from
	addEdges := []int{7, 4, 2, 1}
	edge := 0
	for k, i := range corners {
		addEdgesK := []int{}
		for k2, v := range addEdges {
			if k2 != k {
				addEdgesK = append(addEdgesK, v)
			}
		}
		m := 0
		for j := edge; j < edge+6; j += 2 {
			dc.HalfEdges[j] = dcel.Edge{
				Origin: &dc.Vertices[i],
			}
			dc.OutEdges[i] = &dc.HalfEdges[j]
			dc.HalfEdges[j+1] = dcel.Edge{
				Origin: &dc.Vertices[addEdgesK[m]],
			}
			dc.OutEdges[addEdgesK[m]] = &dc.HalfEdges[j+1]
			m++
		}
		// 6 edges per corner
		edge += 6
	}
	// Set Twins
	for i := range dc.HalfEdges {
		dc.HalfEdges[i].Twin = &dc.HalfEdges[dcel.EdgeTwin(i)]
	}
	// We're ignoring prev, next, face for now
	// because this is way harder than it should be
	return NewPolyhedronFromDCEL(&dc, px, py)
}

func NewPolyhedronFromDCEL(dc *dcel.DCEL, x, y float64) *Polyhedron {
	p := new(Polyhedron)
	p.SetPos(x, y)
	p.DCEL = *dc
	p.Update()
	p.Center = physics.NewVector(p.X+(1+p.MaxX())/2, p.Y+(1+p.MaxY())/2)
	return p
}

func (p *Polyhedron) Update() {

	// Reset p's rgba
	maxX := p.MaxX() + 1
	maxY := p.MaxY() + 1
	rect := image.Rect(0, 0, int(maxX), int(maxY))
	rgba := image.NewRGBA(rect)
	// We ignore Z -- Z is used for rotations
	// There isn't an alternative to this, aside from
	// recoloring things to account for different z values,
	// without having a camera system, which is a lot of work

	// Try to maintain center
	if p.Center.X != 0 || p.Center.Y != 0 {
		cx := p.X + maxX/2
		cy := p.Y + maxY/2
		p.X -= (cx - p.Center.X)
		p.Y -= (cy - p.Center.Y)
	}

	red := color.RGBA{255, 0, 0, 255}
	gray := color.RGBA{160, 160, 160, 255}

	// Step 1: draw all edges
	// Given the edge twin mandate, we can just use
	// every other halfEdge.
	//fmt.Println(p.DCEL)
	for i := 0; i < len(p.HalfEdges); i += 2 {
		points := p.FullEdge(i)
		// draw a line from points[0] to points[1]
		// fmt.Println("Drawing from ", points[0][0], points[0][1], "to",
		// 	points[1][0], points[1][1])
		drawLineOnto(rgba, int(points[0][0]), int(points[0][1]),
			int(points[1][0]), int(points[1][1]), gray)
	}

	// Step 2: draw all vertices
	for _, v := range p.Vertices {
		rgba.Set(int(v[0]), int(v[1]), red)
	}

	p.SetRGBA(rgba)
}

func (p *Polyhedron) RotZ(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for i, v := range p.Vertices {
		p.Vertices[i][0] = v[0]*ct - v[1]*st
		p.Vertices[i][1] = v[1]*ct + v[0]*st
	}
	p.clearNegativePoints()
	p.Update()
}

func (p *Polyhedron) RotX(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for i, v := range p.Vertices {
		p.Vertices[i][1] = v[1]*ct - v[2]*st
		p.Vertices[i][2] = v[2]*ct + v[1]*st
	}
	p.clearNegativePoints()
	p.Update()
}

func (p *Polyhedron) RotY(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for i, v := range p.Vertices {
		p.Vertices[i][0] = v[0]*ct - v[2]*st
		p.Vertices[i][2] = v[2]*ct + v[0]*st
	}
	p.clearNegativePoints()
	p.Update()
}

func (p *Polyhedron) Scale(factor float64) {
	for i, v := range p.Vertices {
		p.Vertices[i][0] = v[0] * factor
		p.Vertices[i][1] = v[1] * factor
		p.Vertices[i][2] = v[2] * factor
	}
	p.clearNegativePoints()
	p.Update()
}

func (p *Polyhedron) clearNegativePoints() {
	// Anything with an x,y less than 0 needs to be increased,
	// this is a limitation so we stay in the bounds of a given rgba
	// rectangle on screen, so we increase everything by minX,minY
	x := p.MinX()
	y := p.MinY()
	for i, v := range p.Vertices {
		p.Vertices[i][0] = v[0] - x
	}
	for i, v := range p.Vertices {
		p.Vertices[i][1] = v[1] - y
	}
}

// Utilities
func (p *Polyhedron) String() string {
	return "Polyhedron"
}

func (p *Polyhedron) ShiftX(x float64) {
	p.Center.X += x
	p.Sprite.ShiftX(x)
}

func (p *Polyhedron) ShiftY(y float64) {
	p.Center.Y += y
	p.Sprite.ShiftY(y)
}
