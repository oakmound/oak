package render

import (
	"image"
	"image/color"
	"math"

	"bitbucket.org/oakmoundstudio/oak/physics"
)

/////////////////////
// Won't be in this package
////////////////////

// This really shouldn't be
// X, Y Z but rather [0,1,2]
// but
type DCELPoint [3]float64

func (dp DCELPoint) X() float64 {
	return dp[0]
}

func (dp DCELPoint) Y() float64 {
	return dp[1]
}

func (dp DCELPoint) Z() float64 {
	return dp[2]
}

type DCELEdge struct {
	// Origin is the vertex this edge starts at
	Origin *DCELPoint
	// Face is the index within Faces that this
	// edge wraps around
	Face *DCELFace
	// Next and Prev are the edges following and
	// preceding this edge that also wrap around
	// Face
	Next *DCELEdge
	Prev *DCELEdge
}

type DCELFace struct {
	Outer, Inner *DCELEdge
}

type DCEL struct {
	Vertices []DCELPoint
	// outEdges[0] is the (an) edge in HalfEdges whose
	// orgin is Vertices[0]
	OutEdges  []*DCELEdge
	HalfEdges []DCELEdge
	// The first value in a face is the outside component
	// of the face, the second value is the inside component
	Faces []DCELFace
}

// Twin is the edge that runs along this edge
// towards Origin, around the opposing face
//
// So HalfEdges[e.Twin].Origin is where this
// edge points to.
// Mandate: twin edges come in pairs
// if i is even, then, i+1 is its pair,
// and otherwise i-i is its pair.
func (d *DCEL) EdgeTwin(i int) int {
	if i%2 == 0 {
		return i + 1
	}
	return i - 1
}

// FullEdge returns the ith edge in the form of its
// two vertices
func (d *DCEL) FullEdge(i int) [2]*DCELPoint {
	e := d.HalfEdges[i]
	e2 := d.HalfEdges[d.EdgeTwin(i)]
	return [2]*DCELPoint{
		e.Origin,
		e2.Origin}
}

// Max functions iterate through vertices
// to find the maximum value along a given axis
// in the DCEL

func (d *DCEL) MaxX() float64 {
	return d.Max(0)
}

func (d *DCEL) MaxY() float64 {
	return d.Max(1)
}

func (d *DCEL) MaxZ() float64 {
	return d.Max(2)
}

func (d *DCEL) Max(i int) (x float64) {
	for _, p := range d.Vertices {
		if p[i] > x {
			x = p[i]
		}
	}
	return x
}

func (d *DCEL) MinX() float64 {
	return d.Min(0)
}

func (d *DCEL) MinY() float64 {
	return d.Min(1)
}

func (d *DCEL) MinZ() float64 {
	return d.Min(2)
}

func (d *DCEL) Min(i int) (x float64) {
	x = math.Inf(1)
	for _, p := range d.Vertices {
		if p[i] < x {
			x = p[i]
		}
	}
	return x
}

///////////////////////

type Polyhedron struct {
	Sprite
	DCEL
	Center physics.Vector
}

func NewCube(x, y, z, w, h, d float64) *Polyhedron {
	p := new(Polyhedron)
	p.SetPos(x, y)

	dc := DCEL{}
	// Is there a smart way to loop through this?
	// verts
	x = w
	y = h
	z = -d
	dc.Vertices = make([]DCELPoint, 8)
	dc.HalfEdges = make([]DCELEdge, 24)
	dc.OutEdges = make([]*DCELEdge, 8)
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

		dc.Vertices[i] = DCELPoint{x, y, z}
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
			dc.HalfEdges[j] = DCELEdge{
				Origin: &dc.Vertices[i],
			}
			dc.OutEdges[i] = &dc.HalfEdges[j]
			dc.HalfEdges[j+1] = DCELEdge{
				Origin: &dc.Vertices[addEdgesK[m]],
			}
			dc.OutEdges[addEdgesK[m]] = &dc.HalfEdges[j+1]
			m++
		}
		// 6 edges per corner
		edge += 6
	}
	// We're ignoring prev, next, face for now
	// because this is way harder than it should be
	p.DCEL = dc
	p.Update()
	p.Center = physics.NewVector(p.X+(1+p.MaxX())/2, p.Y+(1+p.MaxY())/2)
	return p
}

func (p *Polyhedron) Update() {

	// Reset p's rgba
	// Todo: define a better buffer than 5
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
