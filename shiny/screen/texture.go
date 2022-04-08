package screen

import (
	"image"
)

// Texture is a pixel buffer, but not one that is directly accessible as a
// []byte. Conceptually, it could live on a GPU, in another process or even be
// across a network, instead of on a CPU in this process.
//
// Images can be uploaded to Textures, and Textures can be drawn on Windows.
//
// When specifying a sub-Texture via Draw, a Texture's top-left pixel is always
// (0, 0) in its own coordinate space.
type Texture interface {
	Spanner

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

	// Release releases the Texture's resources, after all pending uploads and
	// draws resolve.
	//
	// The behavior of the Texture after Release, whether calling its methods
	// or passing it as an argument, is undefined.
	Release()
}
