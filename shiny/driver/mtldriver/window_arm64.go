//go:build arm64 && darwin
// +build arm64,darwin

package mtldriver

import (
	"image"

	"github.com/oakmound/oak/v3/shiny/driver/common"
	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func (w *Window) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	common.CopyToBGRAFromRGBA(w.bgra, dp, src.RGBA(), sr)
}

func (w *Window) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	common.TransformToBGRAFromRGBA(w.bgra, src2dst, src.(*common.Image).RGBA(), sr)
}

func (w *Window) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, draw.Over)
}
