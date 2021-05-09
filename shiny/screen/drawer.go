package screen

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/math/f64"
)

// Drawer is something you can draw Textures on.
//
// Draw is the most general purpose of this interface's methods. It supports
// arbitrary affine transformations, such as translations, scales and
// rotations.
//
// Copy and Scale are more specific versions of Draw. The affected dst pixels
// are an axis-aligned rectangle, quantized to the pixel grid. Copy copies
// pixels in a 1:1 manner, Scale is more general. They have simpler parameters
// than Draw, using ints instead of float64s.
//
// When drawing on a Window, there will not be any visible effect until Publish
// is called.
type Drawer interface {
	// Draw draws the sub-Texture defined by src and sr to the destination (the
	// method receiver). src2dst defines how to transform src coordinates to
	// dst coordinates. For example, if src2dst is the matrix
	//
	// m00 m01 m02
	// m10 m11 m12
	//
	// then the src-space point (sx, sy) maps to the dst-space point
	// (m00*sx + m01*sy + m02, m10*sx + m11*sy + m12).
	Draw(src2dst f64.Aff3, src Texture, sr image.Rectangle, op draw.Op)

	// DrawUniform is like Draw except that the src is a uniform color instead
	// of a Texture.
	DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op)

	// Copy copies the sub-Texture defined by src and sr to the destination
	// (the method receiver), such that sr.Min in src-space aligns with dp in
	// dst-space.
	Copy(dp image.Point, src Texture, sr image.Rectangle, op draw.Op)

	// Scale scales the sub-Texture defined by src and sr to the destination
	// (the method receiver), such that sr in src-space is mapped to dr in
	// dst-space.
	Scale(dr image.Rectangle, src Texture, sr image.Rectangle, op draw.Op)

	// Upload uploads the sub-Buffer defined by src and sr to the destination
	// (the method receiver), such that sr.Min in src-space aligns with dp in
	// dst-space. The destination's contents are overwritten; the draw operator
	// is implicitly draw.Src.
	//
	// It is valid to upload a Buffer while another upload of the same Buffer
	// is in progress, but a Buffer's image.RGBA pixel contents should not be
	// accessed while it is uploading. A Buffer is re-usable, in that its pixel
	// contents can be further modified, once all outstanding calls to Upload
	// have returned.
	//
	// TODO: make it optional that a Buffer's contents is preserved after
	// Upload? Undoing a swizzle is a non-trivial amount of work, and can be
	// redundant if the next paint cycle starts by clearing the buffer.
	//
	// When uploading to a Window, there will not be any visible effect until
	// Publish is called.
	Upload(dp image.Point, src Image, sr image.Rectangle)

	// Fill fills that part of the destination (the method receiver) defined by
	// dr with the given color.
	//
	// When filling a Window, there will not be any visible effect until
	// Publish is called.
	Fill(dr image.Rectangle, src color.Color, op draw.Op)
}
