package btn

import (
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
)

// TextBox is a Box with an associated text element
type TextBox struct {
	Box
	*render.Text
}

// Init creates the CID
func (b *TextBox) Init() event.CID {
	b.CID = event.NextID(b)
	return b.CID
}

// NewTextBox creates a textbox
func NewTextBox(cid event.CID, x, y, w, h, txtX, txtY float64,
	f *render.Font, r render.Renderable, layers ...int) *TextBox {

	if f == nil {
		f = render.DefFont()
	}

	b := new(TextBox)

	cid = cid.Parse(b)

	b.Box = *NewBox(cid, x, y, w, h, r, layers...)
	b.Text = f.NewStrText("Init", 0, 0)
	b.Text.Vector = b.Text.Attach(b.Box.Vector, txtX, (-txtY)+b.H)

	// We dont want to modify the input's layers but we do want the text to show up on top of the base renderable.
	txtLayers := make([]int, len(layers))
	copy(txtLayers, layers)
	txtLayers[len(txtLayers)-1]++
	render.Draw(b.Text, txtLayers...)
	return b
}

// Y pulls the y of the composed Box (disambiguation with the y of the text component)
func (b *TextBox) Y() float64 {
	return b.Box.Y()
}

// X pulls the x of the composed Box (disambiguation with the x of the text component)
func (b *TextBox) X() float64 {
	return b.Box.X()
}

// ShiftX shifts the box by x. The associated text is attached and so will be moved along by default
func (b *TextBox) ShiftX(x float64) {
	b.Box.ShiftX(x)
}

// ShiftY shifts the box by y. The associated text is attached and so will be moved along by default
func (b *TextBox) ShiftY(y float64) {
	b.Box.ShiftY(y)
}

// SetSpace overwrites entities.Solid,
// pointing this button to use the mouse collision Rtree
// instead of the entity collision space.
func (b *TextBox) SetSpace(sp *collision.Space) {
	mouse.Remove(b.Space)
	b.Space = sp
	mouse.Add(b.Space)
}

// SetPos acts as SetSpace does, overwriting entities.Solid.
func (b *TextBox) SetPos(x, y float64) {
	b.Box.SetPos(x, y)
}

// SetOffsets changes the text position within the box
func (b *TextBox) SetOffsets(txtX, txtY float64) {
	b.Text.Vector = b.Text.Attach(b.Box.Vector, txtX, -txtY+b.H)
}
