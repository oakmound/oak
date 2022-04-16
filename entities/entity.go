package entities

import (
	"image/color"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
)

type Generator struct {
	Position   floatgeom.Point2
	Dimensions floatgeom.Point2
	Speed      floatgeom.Point2

	Parent event.Caller

	Color      color.Color
	Renderable render.Renderable

	Mod mod.Mod

	Label collision.Label

	DrawLayers []int

	UseMouseTree     bool
	WithoutCollision bool

	Children [][]Option
}

func And(opts ...Option) Option {
	return func(g Generator) Generator {
		for _, o := range opts {
			g = o(g)
		}
		return g
	}
}

func WithChild(opts ...Option) Option {
	return func(s Generator) Generator {
		s.Children = append(s.Children, opts)
		return s
	}
}

func WithRect(v floatgeom.Rect2) Option {
	return func(s Generator) Generator {
		s.Position = v.Min
		s.Dimensions = v.Max.Sub(v.Min)
		return s
	}
}

func WithOffset(p floatgeom.Point2) Option {
	return func(g Generator) Generator {
		g.Position = g.Position.Add(p)
		return g
	}
}

var defaultGenerator = Generator{
	Dimensions: floatgeom.Point2{1, 1},
	DrawLayers: []int{0},
}

type Entity struct {
	event.CallerID

	ctx *scene.Context

	Rect  floatgeom.Rect2
	Speed floatgeom.Point2
	Delta floatgeom.Point2

	Renderable render.Renderable

	collision.Phase

	Space *collision.Space
	Tree  *collision.Tree

	metadata map[string]string

	Children []*Entity
}

func (e Entity) CID() event.CallerID {
	return e.CallerID.CID()
}

func (e Entity) X() float64 {
	return e.Rect.Min.X()
}
func (e Entity) Y() float64 {
	return e.Rect.Min.Y()
}
func (e Entity) W() float64 {
	return e.Rect.W()
}
func (e Entity) H() float64 {
	return e.Rect.H()
}
func (e Entity) Top() float64 {
	return e.Y()
}
func (e Entity) Bottom() float64 {
	return e.Y() + e.H()
}
func (e Entity) Left() float64 {
	return e.X()
}
func (e Entity) Right() float64 {
	return e.X() + e.W()
}

func (e *Entity) ShiftDelta() {
	e.Shift(e.Delta)
}

func (e *Entity) Shift(delta floatgeom.Point2) {
	// TODO: attachment?
	// TODO: helper
	e.Renderable.ShiftX(delta.X())
	e.Renderable.ShiftY(delta.Y())
	e.Rect = e.Rect.Shift(delta)
	if e.Tree != nil {
		e.Tree.UpdateSpace(
			e.X(), e.Y(), e.W(), e.H(), e.Space,
		)
	}
}

func (e *Entity) SetX(x float64) {
	e.ShiftX(x - e.X())
}

func (e *Entity) SetY(y float64) {
	e.ShiftY(y - e.Y())
}

func (e *Entity) ShiftX(x float64) {
	e.Renderable.ShiftX(x)
	e.Rect = e.Rect.Shift(floatgeom.Point2{x, 0})
	if e.Tree != nil {
		e.Tree.UpdateSpace(
			e.X(), e.Y(), e.W(), e.H(), e.Space,
		)
	}
}

func (e *Entity) ShiftY(y float64) {
	e.Renderable.ShiftY(y)
	e.Rect = e.Rect.Shift(floatgeom.Point2{0, y})
	if e.Tree != nil {
		e.Tree.UpdateSpace(
			e.X(), e.Y(), e.W(), e.H(), e.Space,
		)
	}
}

func (e *Entity) SetPos(p floatgeom.Point2) {
	w, h := e.W(), e.H()
	e.Rect = floatgeom.NewRect2WH(p.X(), p.Y(), w, h)
	e.Renderable.SetPos(p.X(), p.Y())
	if e.Tree != nil {
		e.Tree.UpdateSpace(
			e.X(), e.Y(), e.W(), e.H(), e.Space,
		)
	}
}

// TODO: take a point, not floats
func (e *Entity) ShiftPos(x, y float64) {
	p := e.Rect.Min
	e.SetPos(p.Add(floatgeom.Point2{x, y}))
}

func (e *Entity) HitLabel(label collision.Label) *collision.Space {
	return e.Tree.HitLabel(e.Space, label)
}

func (e *Entity) Destroy() {
	e.Renderable.Undraw()
	e.Tree.Remove(e.Space)
	e.ctx.UnbindAllFrom(e.CallerID)
}

// SetMetadata sets the metadata for some key to some value. Empty value strings
// will not be stored.
func (e *Entity) SetMetadata(k, v string) {
	if v == "" {
		delete(e.metadata, k)
	} else {
		e.metadata[k] = v
	}
}

// Metadata accesses the value, and whether it existed, for a given metadata key
func (e *Entity) Metadata(k string) (v string, ok bool) {
	v, ok = e.metadata[k]
	return v, ok
}

func New(ctx *scene.Context, opts ...Option) *Entity {
	g := defaultGenerator
	for _, o := range opts {
		g = o(g)
	}

	children := make([]*Entity, len(g.Children))
	for i, childOpts := range g.Children {
		childOpts = append(childOpts, WithOffset(g.Position))
		children[i] = New(ctx, childOpts...)
	}

	e := &Entity{
		ctx: ctx,
		Rect: floatgeom.NewRect2WH(
			g.Position[0],
			g.Position[1],
			g.Dimensions[0],
			g.Dimensions[1],
		),
		Renderable: g.Renderable,
		Speed:      g.Speed,
		Children:   children,
	}

	if g.Renderable == nil && g.Color != nil {
		e.Renderable = render.NewColorBox(int(e.W()), int(e.H()), g.Color)
	}

	if m, isMod := e.Renderable.(render.Modifiable); g.Mod != nil && isMod {
		e.Renderable = m.Modify(g.Mod)
	}

	e.Renderable.SetPos(e.X(), e.Y())

	if g.Parent == nil {
		cid := ctx.CallerMap.Register(e)
		e.CallerID = cid
	} else {
		e.CallerID = g.Parent.CID()
		if e.CallerID == 0 {
			dlog.Error("entity created with uninitialized parent caller ID")
		}
	}

	if !g.WithoutCollision {
		e.Tree = ctx.CollisionTree
		if g.UseMouseTree {
			e.Tree = ctx.MouseTree
		}
		e.Space = collision.NewSpace(
			e.X(), e.Y(), e.W(), e.H(), e.CallerID,
		)
		e.Space.Label = g.Label
		e.Tree.Add(e.Space)
	}

	if len(g.DrawLayers) != 0 {
		ctx.Draw(e.Renderable, g.DrawLayers...)
	}

	return e
}
