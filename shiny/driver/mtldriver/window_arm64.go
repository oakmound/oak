//go:build arm64 && darwin
// +build arm64,darwin

package mtldriver

import (
	"image"

	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func (w *Window) Upload(dp image.Point, srcImg screen.Image, sr image.Rectangle) {
	dst := w.bgra
	r := sr.Sub(sr.Min).Add(dp)
	src := srcImg.RGBA()
	sp := sr.Min
	clip(dst, &r, src, &sp, nil, &image.Point{})
	if r.Empty() {
		return
	}

	i0 := (r.Min.X - dst.Rect.Min.X) * 4
	i1 := (r.Max.X - dst.Rect.Min.X) * 4
	si0 := (sp.X - src.Rect.Min.X) * 4
	yMax := r.Max.Y - dst.Rect.Min.Y

	y := r.Min.Y - dst.Rect.Min.Y
	sy := sp.Y - src.Rect.Min.Y
	for ; y != yMax; y, sy = y+1, sy+1 {
		dpix := dst.Pix[y*dst.Stride:]
		spix := src.Pix[sy*src.Stride:]

		for i, si := i0, si0; i < i1; i, si = i+4, si+4 {
			s := spix[si : si+4 : si+4] // Small cap improves performance, see https://golang.org/issue/27857
			d := dpix[i : i+4 : i+4]
			d[0] = s[2]
			d[1] = s[1]
			d[2] = s[0]
			d[3] = s[3]
		}
	}
}

func (w *Window) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	nnInterpolator{}.Transform(w.bgra, src2dst, src.(*textureImpl).rgba, sr)
}

func (w *Window) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, draw.Over)
}
