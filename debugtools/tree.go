package debugtools

import (
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"

	"github.com/oakmound/oak/v3/collision"
)

// NewRTree creates a wrapper around a tree that supports coloring the spaces
func NewRTree(ctx *scene.Context, t *collision.Tree) *Rtree {
	return NewThickRTree(ctx, t, 1)
}

// NewThickRTree creates a wrapper around tree that colors spaces up to a thickness
func NewThickRTree(ctx *scene.Context, t *collision.Tree, thickness int) *Rtree {
	return NewThickColoredRTree(ctx, t, thickness, map[collision.Label]color.RGBA{})
}

// NewThickColoredRTree creates a wrapper around tree that colors spaces up to a thickness based on a coloring map
func NewThickColoredRTree(ctx *scene.Context, t *collision.Tree, thickness int, colorMapping map[collision.Label]color.RGBA) *Rtree {
	rt := new(Rtree)
	rt.Tree = t
	rt.Context = ctx
	rt.Thickness = thickness
	rt.LayeredPoint = render.NewLayeredPoint(0, 0, -1)
	rt.OutlineColor = color.RGBA{200, 200, 200, 255}
	rt.ColorMap = colorMapping
	return rt
}

// An Rtree wraps around a collision tree and can draw debug rectangles for every entity in
// the tree.
type Rtree struct {
	*collision.Tree
	Thickness int
	render.LayeredPoint
	OutlineColor color.RGBA
	ColorMap     map[collision.Label]color.RGBA
	DrawDisabled bool
	Context      *scene.Context
}

// GetDims returns the total possible area to draw this on.
func (r *Rtree) GetDims() (int, int) {
	bds := r.Context.Window.Bounds()
	return bds.X(), bds.Y()
}

// Draw will draw the collision outlines
func (r *Rtree) Draw(buff draw.Image, xOff, yOff float64) {
	if r.DrawDisabled {
		return
	}
	vp := r.Context.Window.Viewport()
	bds := r.Context.Window.Bounds()
	// Get all spaces on screen
	screen := collision.NewUnassignedSpace(
		float64(vp.X()),
		float64(vp.Y()),
		float64(bds.X()+vp.X()),
		float64(bds.Y()+vp.Y()))
	hits := r.Tree.Hits(screen)
	// Draw spaces that are on screen (as outlines)
	for _, h := range hits {
		c := r.OutlineColor
		if found, ok := r.ColorMap[h.Label]; ok {
			c = found
		}
		for x := 0; x < int(h.GetW()); x++ {
			for i := 0; i < r.Thickness; i++ {
				buff.Set(x+int(h.X()+xOff)-vp.X(), int(h.Y()+yOff)+i-vp.Y(), c)
				buff.Set(x+int(h.X()+xOff)-vp.X(), int(h.Y()+yOff)+int(h.GetH())-i-vp.Y(), c)
			}
		}
		for y := 0; y < int(h.GetH()); y++ {
			for i := 0; i < r.Thickness; i++ {
				buff.Set(int(h.X()+xOff)+i-vp.X(), y+int(h.Y()+yOff)-vp.Y(), c)
				buff.Set(int(h.X()+xOff)+int(h.GetW())-i-vp.X(), y+int(h.Y()+yOff)-vp.Y(), c)
			}
		}
	}
}
