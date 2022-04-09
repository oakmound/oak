//go:build arm64 && darwin
// +build arm64,darwin

package mtldriver

import (
	"image"
	"image/color"
	"math"
	"math/bits"

	"golang.org/x/image/math/f64"
)

// This file is a copy of much of x/image/draw and draw/image
// To enable fast conversions from RGBA (which Oak uses everywhere internally)
// and BGRA (which metal refuses not to use for windows)

var _ image.Image = &BGRA{}

// BGRA is an in-memory image whose At method returns BGRA values.
type BGRA struct {
	// Pix holds the image's pixels, in B, G, R, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *BGRA) ColorModel() color.Model { return color.RGBAModel }

func (p *BGRA) Bounds() image.Rectangle { return p.Rect }

func (p *BGRA) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *BGRA) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	r := uint16(s[2])
	g := uint16(s[1])
	b := uint16(s[0])
	a := uint16(s[3])
	return color.RGBA64{
		(r << 8) | r,
		(g << 8) | g,
		(b << 8) | b,
		(a << 8) | a,
	}
}

func (p *BGRA) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return color.RGBA{s[2], s[1], s[0], s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *BGRA) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *BGRA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[2] = c1.R
	s[1] = c1.G
	s[0] = c1.B
	s[3] = c1.A
}

func (p *BGRA) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[2] = uint8(c.R >> 8)
	s[1] = uint8(c.G >> 8)
	s[0] = uint8(c.B >> 8)
	s[3] = uint8(c.A >> 8)
}

func (p *BGRA) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *BGRA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &BGRA{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &BGRA{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *BGRA) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// NewBGRA returns a new RGBA image with the given bounds.
func NewBGRA(r image.Rectangle) *BGRA {
	return &BGRA{
		Pix:    make([]uint8, pixelBufferLength(4, r, "BGRA")),
		Stride: 4 * r.Dx(),
		Rect:   r,
	}
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}

// clip clips r against each image's bounds (after translating into the
// destination image's coordinate space) and shifts the points sp and mp by
// the same amount as the change in r.Min.
func clip(dst *BGRA, r *image.Rectangle, src *image.RGBA, sp *image.Point, mask image.Image, mp *image.Point) {
	orig := r.Min
	*r = r.Intersect(dst.Bounds())
	*r = r.Intersect(src.Bounds().Add(orig.Sub(*sp)))
	if mask != nil {
		*r = r.Intersect(mask.Bounds().Add(orig.Sub(*mp)))
	}
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	if dx == 0 && dy == 0 {
		return
	}
	sp.X += dx
	sp.Y += dy
	if mp != nil {
		mp.X += dx
		mp.Y += dy
	}
}

type nnInterpolator struct{}

func (z nnInterpolator) Transform(dst *BGRA, s2d f64.Aff3, src *image.RGBA, sr image.Rectangle) {
	// Try to simplify a Transform to a Copy.
	// if s2d[0] == 1 && s2d[1] == 0 && s2d[3] == 0 && s2d[4] == 1 {
	// 	dx := int(s2d[2])
	// 	dy := int(s2d[5])
	// 	if float64(dx) == s2d[2] && float64(dy) == s2d[5] {
	// 		Copy(dst, image.Point{X: sr.Min.X + dx, Y: sr.Min.X + dy}, src, sr, op, opts)
	// 		return
	// 	}
	// }

	dr := transformRect(&s2d, &sr)
	// adr is the affected destination pixels.
	adr := dst.Bounds().Intersect(dr)
	if adr.Empty() || sr.Empty() {
		return
	}
	d2s := invert(&s2d)
	// bias is a translation of the mapping from dst coordinates to src
	// coordinates such that the latter temporarily have non-negative X
	// and Y coordinates. This allows us to write int(f) instead of
	// int(math.Floor(f)), since "round to zero" and "round down" are
	// equivalent when f >= 0, but the former is much cheaper. The X--
	// and Y-- are because the TransformLeaf methods have a "sx -= 0.5"
	// adjustment.
	bias := transformRect(&d2s, &adr).Min
	bias.X--
	bias.Y--
	d2s[2] -= float64(bias.X)
	d2s[5] -= float64(bias.Y)
	// Make adr relative to dr.Min.
	adr = adr.Sub(dr.Min)
	// sr is the source pixels. If it extends beyond the src bounds,
	// we cannot use the type-specific fast paths, as they access
	// the Pix fields directly without bounds checking.
	//
	// Similarly, the fast paths assume that the masks are nil.
	z.transform_BGRA_RGBA_Over(dst, dr, adr, &d2s, src, sr, bias)
}

func (nnInterpolator) transform_BGRA_RGBA_Over(dst *BGRA, dr, adr image.Rectangle, d2s *f64.Aff3, src *image.RGBA, sr image.Rectangle, bias image.Point) {
	for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
		dyf := float64(dr.Min.Y+int(dy)) + 0.5
		d := (dr.Min.Y+int(dy)-dst.Rect.Min.Y)*dst.Stride + (dr.Min.X+adr.Min.X-dst.Rect.Min.X)*4
		for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx, d = dx+1, d+4 {
			dxf := float64(dr.Min.X+int(dx)) + 0.5
			sx0 := int(d2s[0]*dxf+d2s[1]*dyf+d2s[2]) + bias.X
			sy0 := int(d2s[3]*dxf+d2s[4]*dyf+d2s[5]) + bias.Y
			if !(image.Point{sx0, sy0}).In(sr) {
				continue
			}
			pi := (sy0-src.Rect.Min.Y)*src.Stride + (sx0-src.Rect.Min.X)*4
			pr := uint32(src.Pix[pi+0]) * 0x101
			pg := uint32(src.Pix[pi+1]) * 0x101
			pb := uint32(src.Pix[pi+2]) * 0x101
			pa := uint32(src.Pix[pi+3]) * 0x101
			pa1 := (0xffff - pa) * 0x101
			dst.Pix[d+2] = uint8((uint32(dst.Pix[d+0])*pa1/0xffff + pr) >> 8)
			dst.Pix[d+1] = uint8((uint32(dst.Pix[d+1])*pa1/0xffff + pg) >> 8)
			dst.Pix[d+0] = uint8((uint32(dst.Pix[d+2])*pa1/0xffff + pb) >> 8)
			dst.Pix[d+3] = uint8((uint32(dst.Pix[d+3])*pa1/0xffff + pa) >> 8)
		}
	}
}

// transformRect returns a rectangle dr that contains sr transformed by s2d.
func transformRect(s2d *f64.Aff3, sr *image.Rectangle) (dr image.Rectangle) {
	ps := [...]image.Point{
		{sr.Min.X, sr.Min.Y},
		{sr.Max.X, sr.Min.Y},
		{sr.Min.X, sr.Max.Y},
		{sr.Max.X, sr.Max.Y},
	}
	for i, p := range ps {
		sxf := float64(p.X)
		syf := float64(p.Y)
		dx := int(math.Floor(s2d[0]*sxf + s2d[1]*syf + s2d[2]))
		dy := int(math.Floor(s2d[3]*sxf + s2d[4]*syf + s2d[5]))

		// The +1 adjustments below are because an image.Rectangle is inclusive
		// on the low end but exclusive on the high end.

		if i == 0 {
			dr = image.Rectangle{
				Min: image.Point{dx + 0, dy + 0},
				Max: image.Point{dx + 1, dy + 1},
			}
			continue
		}

		if dr.Min.X > dx {
			dr.Min.X = dx
		}
		dx++
		if dr.Max.X < dx {
			dr.Max.X = dx
		}

		if dr.Min.Y > dy {
			dr.Min.Y = dy
		}
		dy++
		if dr.Max.Y < dy {
			dr.Max.Y = dy
		}
	}
	return dr
}

func clipAffectedDestRect(adr image.Rectangle, dstMask image.Image, dstMaskP image.Point) (image.Rectangle, image.Image) {
	if dstMask == nil {
		return adr, nil
	}
	if r, ok := dstMask.(image.Rectangle); ok {
		return adr.Intersect(r.Sub(dstMaskP)), nil
	}
	// TODO: clip to dstMask.Bounds() if the color model implies that out-of-bounds means 0 alpha?
	return adr, dstMask
}

func invert(m *f64.Aff3) f64.Aff3 {
	m00 := +m[3*1+1]
	m01 := -m[3*0+1]
	m02 := +m[3*1+2]*m[3*0+1] - m[3*1+1]*m[3*0+2]
	m10 := -m[3*1+0]
	m11 := +m[3*0+0]
	m12 := +m[3*1+0]*m[3*0+2] - m[3*1+2]*m[3*0+0]

	det := m00*m11 - m10*m01

	return f64.Aff3{
		m00 / det,
		m01 / det,
		m02 / det,
		m10 / det,
		m11 / det,
		m12 / det,
	}
}
