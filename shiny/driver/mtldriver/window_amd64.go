//go:build darwin && amd64
// +build darwin,amd64

package mtldriver

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func (w *windowImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	draw.Draw(w.bgra, sr.Sub(sr.Min).Add(dp), src.RGBA(), sr.Min, draw.Src)
}

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	draw.Draw(w.bgra, dr, &image.Uniform{src}, image.Point{}, op)
}

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	draw.NearestNeighbor.Transform(w.bgra, src2dst, src.(*textureImpl).rgba, sr, op, nil)
}

func (w *windowImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {
	draw.NearestNeighbor.Transform(w.bgra, src2dst, &image.Uniform{src}, sr, op, nil)
}

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Copy(w, dp, src, sr, op)
}

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, op)
}

type BGRA = image.RGBA

var NewBGRA = image.NewRGBA
