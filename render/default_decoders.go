//go:build !noimages
// +build !noimages

package render

import (
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
)

func init() {
	// Register standard image decoders. If provided with the build tag 'noimages', this is skipped.
	RegisterDecoder(".jpeg", jpeg.Decode)
	RegisterDecoder(".jpg", jpeg.Decode)
	RegisterDecoder(".gif", gif.Decode)
	RegisterDecoder(".png", png.Decode)
	RegisterDecoder(".bmp", bmp.Decode)
	RegisterCfgDecoder(".jpeg", jpeg.DecodeConfig)
	RegisterCfgDecoder(".jpg", jpeg.DecodeConfig)
	RegisterCfgDecoder(".gif", gif.DecodeConfig)
	RegisterCfgDecoder(".png", png.DecodeConfig)
	RegisterCfgDecoder(".bmp", bmp.DecodeConfig)
}
