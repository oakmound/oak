package render

import (
	"image"
	"image/color"
	"math"
	"sort"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-compgeo/dcel"
)

type Polyhedron struct {
	Sprite
	dcel.DCEL
	FaceColors []color.Color
	EdgeColors []color.Color
	Center     physics.Vector
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
	dc.Vertices = make([]*dcel.Point, 8)
	dc.HalfEdges = make([]*dcel.Edge, 24)
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

		dc.Vertices[i] = dcel.NewPoint(x, y, z)
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
			dc.HalfEdges[j] = &dcel.Edge{
				Origin: dc.Vertices[i],
			}
			dc.OutEdges[i] = dc.HalfEdges[j]
			dc.HalfEdges[j+1] = &dcel.Edge{
				Origin: dc.Vertices[addEdgesK[m]],
			}
			dc.OutEdges[addEdgesK[m]] = dc.HalfEdges[j+1]
			m++
		}
		// 6 edges per corner
		edge += 6
	}
	// Set Twins
	for i := range dc.HalfEdges {
		dc.HalfEdges[i].Twin = dc.HalfEdges[dcel.EdgeTwin(i)]
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

	p.clearNegativePoints()

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

	// Eventually:
	// For all Faces, Edges, and Vertices, sort by z value
	// and draw them high-to-low
	zOrder := make([]interface{}, len(p.HalfEdges)/2+len(p.Faces)-1+len(p.Vertices))

	// Step 1: draw all edges
	// Given the edge twin mandate, we can just use
	// every other halfEdge.
	if len(p.EdgeColors) < len(p.HalfEdges) {
		diff := len(p.HalfEdges) - len(p.EdgeColors)
		p.EdgeColors = append(p.EdgeColors, make([]color.Color, diff)...)
	}

	for i := 0; i < len(p.HalfEdges); i += 2 {
		points, err := p.FullEdge(i)
		if err != nil {
			continue
		}
		if p.EdgeColors[i] == nil {
			p.EdgeColors[i] = color.RGBA{255, 0, 255, 255}
		}
		zOrder = append(zOrder, coloredEdge{points, p.EdgeColors[i]})
	}

	// Step 2: draw all vertices
	for _, v := range p.Vertices {
		zOrder = append(zOrder, v)
	}

	if len(p.FaceColors) < len(p.Faces) {
		diff := len(p.Faces) - len(p.FaceColors)
		p.FaceColors = append(p.FaceColors, make([]color.Color, diff)...)
	}

	for i := 1; i < len(p.Faces); i++ {
		f := p.Faces[i]
		if p.FaceColors[i] == nil {
			p.FaceColors[i] = color.RGBA{0, 255, 255, 255}
		}
		verts := f.Vertices()
		max_z := math.MaxFloat64 * -1
		phys_verts := make([]physics.Vector, len(verts))
		for i, v := range verts {
			phys_verts[i] = physics.NewVector(v[0], v[1])
			if v[2] > max_z {
				max_z = v[2]
			}
		}
		poly, err := NewPolygon(phys_verts)
		if err != nil {
			continue
		}
		fpoly := facePolygon{
			poly,
			max_z,
			p.FaceColors[i],
		}
		zOrder = append(zOrder, fpoly)
	}

	// This is very hacky
	sort.Slice(zOrder, func(i, j int) bool {
		z1 := 0.0
		z2 := 0.0
		switch v := zOrder[i].(type) {
		case facePolygon:
			z1 = v.z
		case dcel.Point:
			z1 = v[2] + .002
		case coloredEdge:
			z1 = math.Max(v.ps[0][2], v.ps[1][2]) + .001
		}
		switch v := zOrder[j].(type) {
		case facePolygon:
			z2 = v.z
		case dcel.Point:
			z2 = v[2] + .002
		case coloredEdge:
			z2 = math.Max(v.ps[0][2], v.ps[1][2]) + .001
		}
		return z1 < z2
	})

	for _, item := range zOrder {
		switch v := item.(type) {
		case facePolygon:
			for x := v.Rect.minX; x < v.Rect.maxX; x++ {
				for y := v.Rect.minY; y < v.Rect.maxY; y++ {
					if v.Contains(x, y) {
						rgba.Set(int(x), int(y), v.c)
					}
				}
			}
		case dcel.Point:
			rgba.Set(int(v[0]), int(v[1]), red)
		case coloredEdge:
			drawLineOnto(rgba, int(v.ps[0][0]), int(v.ps[0][1]),
				int(v.ps[1][0]), int(v.ps[1][1]), v.c)
		}
	}

	p.SetRGBA(rgba)
}

type coloredEdge struct {
	ps [2]*dcel.Point
	c  color.Color
}

type facePolygon struct {
	*Polygon
	z float64
	c color.Color
}

func (p *Polyhedron) RotZ(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v0 := v[0]*ct - v[1]*st
		v[1] = v[1]*ct + v[0]*st
		v[0] = v0
	}
	p.Update()
}

func (p *Polyhedron) RotX(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v1 := v[1]*ct - v[2]*st
		v[2] = v[2]*ct + v[1]*st
		v[1] = v1
	}
	p.Update()
}

func (p *Polyhedron) RotY(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v0 := v[0]*ct - v[2]*st
		v[2] = v[2]*ct + v[0]*st
		v[0] = v0
	}
	p.Update()
}

func (p *Polyhedron) Scale(factor float64) {
	for i, v := range p.Vertices {
		p.Vertices[i][0] = v[0] * factor
		p.Vertices[i][1] = v[1] * factor
		p.Vertices[i][2] = v[2] * factor
	}
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
