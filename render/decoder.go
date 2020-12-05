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

var (
	fileDecoders = map[string]Decoder{
		".jpeg": jpeg.Decode,
		".jpg":  jpeg.Decode,
		".gif":  gif.Decode,
		".png":  png.Decode,
		".bmp":  bmp.Decode,
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
