package screen

import "image"

// Spanner types span some distance. This distance can either be returned as
// a rectangle from 0,0 or as a size point.
type Spanner interface {
	// Size returns the size of this Spanner.
	Size() image.Point

	// Bounds returns the bounds of this Spanner's span.
	Bounds() image.Rectangle
}
