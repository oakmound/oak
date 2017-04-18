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
			phys_verts[i] = physics.NewVector(v.X(), v.Y())
			if v.Z() > max_z {
				max_z = v.Z()
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
			z1 = math.Max(v.ps[0].Z(), v.ps[1].Z()) + .001
		}
		switch v := zOrder[j].(type) {
		case facePolygon:
			z2 = v.z
		case dcel.Point:
			z2 = v[2] + .002
		case coloredEdge:
			z2 = math.Max(v.ps[0].Z(), v.ps[1].Z()) + .001
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
			drawLineOnto(rgba, int(v.ps[0].X()), int(v.ps[0].Y()),
				int(v.ps[1].X()), int(v.ps[1].Y()), v.c)
		}
	}

	p.SetRGBA(rgba)
}

type coloredEdge struct {
	ps dcel.FullEdge
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
		v0 := v.X()*ct - v.Y()*st
		v.Point[1] = v.Y()*ct + v.X()*st
		v.Point[0] = v0
	}
	p.Update()
}

func (p *Polyhedron) RotX(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v1 := v.Y()*ct - v.Z()*st
		v.Point[2] = v.Z()*ct + v.Y()*st
		v.Point[1] = v1
	}
	p.Update()
}

func (p *Polyhedron) RotY(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v0 := v.X()*ct - v.Z()*st
		v.Point[2] = v.Z()*ct + v.X()*st
		v.Point[0] = v0
	}
	p.Update()
}

func (p *Polyhedron) Scale(factor float64) {
	for _, v := range p.Vertices {
		v.Point = dcel.Point{
			v.X() * factor,
			v.Y() * factor,
			v.Z() * factor,
		}
	}
	p.Update()
}

func (p *Polyhedron) clearNegativePoints() {
	// Anything with an x,y less than 0 needs to be increased,
	// this is a limitation so we stay in the bounds of a given rgba
	// rectangle on screen, so we increase everything by minX,minY
	x := p.MinX()
	y := p.MinY()
	for _, v := range p.Vertices {
		v.Point[0] = v.X() - x
	}
	for _, v := range p.Vertices {
		v.Point[1] = v.Y() - y
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
