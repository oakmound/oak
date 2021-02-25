package render

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/oakmound/oak/v2/oakerr"

	"golang.org/x/image/bmp"
)

// Decoder functions convert arbitrary readers to images.
// The input of a decoder in oak's loader will generally
// be an image file.
type Decoder func(io.Reader) (image.Image, error)

// CfgDecoder is an equivalent to Decoder that just exports
// the color model and dimensions of the image.
type CfgDecoder func(io.Reader) (image.Config, error)

var (
	fileDecoders = map[string]Decoder{
		".jpeg": jpeg.Decode,
		".jpg":  jpeg.Decode,
		".gif":  gif.Decode,
		".png":  png.Decode,
		".bmp":  bmp.Decode,
	}
	cfgDecoders = map[string]CfgDecoder{
		".jpeg": jpeg.DecodeConfig,
		".jpg":  jpeg.DecodeConfig,
		".gif":  gif.DecodeConfig,
		".png":  png.DecodeConfig,
		".bmp":  bmp.DecodeConfig,
	}
)

// RegisterDecoder adds a decoder to the set of image decoders
// for file loading. If the extension string is already set,
// the existing decoder will not be overwritten.
func RegisterDecoder(ext string, decoder Decoder) error {
	_, ok := fileDecoders[ext]
	if ok {
		return oakerr.ExistingElement{
			InputName:   "ext",
			InputType:   "string",
			Overwritten: false,
		}
	}
	fileDecoders[ext] = decoder
	return nil
}

// RegisterCfgDecoder acts like RegisterDecoder for CfgDecoders
func RegisterCfgDecoder(ext string, decoder CfgDecoder) error {
	_, ok := cfgDecoders[ext]
	if ok {
		return oakerr.ExistingElement{
			InputName:   "ext",
			InputType:   "string",
			Overwritten: false,
		}
	}
	cfgDecoders[ext] = decoder
	return nil
}
