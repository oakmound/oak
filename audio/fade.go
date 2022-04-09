package audio

import (
	"time"

	"github.com/oakmound/oak/v3/audio/pcm"
)

func FadeIn(dur time.Duration, in pcm.Reader) pcm.Reader {
	perSec := in.PCMFormat().BytesPerSecond()
	bytesToFadeIn := int((time.Duration(perSec) / 1000) * (dur / time.Millisecond))

	return &fadeInReader{
		Reader:        in,
		toFadeIn:      bytesToFadeIn,
		totalToFadeIn: bytesToFadeIn,
	}
}

type fadeInReader struct {
	pcm.Reader
	toFadeIn, totalToFadeIn int
}

func (fir *fadeInReader) ReadPCM(b []byte) (n int, err error) {
	if fir.toFadeIn == 0 {
		return fir.Reader.ReadPCM(b)
	}
	read, err := fir.Reader.ReadPCM(b)
	if err != nil {
		return read, err
	}
	format := fir.PCMFormat()
	switch format.Bits {
	case 8:
		for i, byt := range b[:read] {
			fadeInPercent := (float64(fir.totalToFadeIn) - float64(fir.toFadeIn)) / float64(fir.totalToFadeIn)
			if fadeInPercent >= 1 {
				fadeInPercent = 1
			}
			b[i] = byte(int8(float64(int8(byt)) * fadeInPercent))
			fir.toFadeIn--
		}
	case 16:
		for i := 0; i+2 <= read; i += 2 {
			fadeInPercent := (float64(fir.totalToFadeIn) - float64(fir.toFadeIn)) / float64(fir.totalToFadeIn)
			if fadeInPercent >= 1 {
				fadeInPercent = 1
			}
			i16 := int16(b[i]) + (int16(b[i+1]) << 8)
			new16 := int16(float64(i16) * fadeInPercent)
			b[i] = byte(new16)
			b[i+1] = byte(new16 >> 8)
			fir.toFadeIn -= 2
		}
	case 32:
		for i := 0; i+4 <= read; i += 4 {
			fadeInPercent := (float64(fir.totalToFadeIn) - float64(fir.toFadeIn)) / float64(fir.totalToFadeIn)
			if fadeInPercent >= 1 {
				fadeInPercent = 1
			}
			i32 := int32(b[i]) +
				(int32(b[i+1]) << 8) +
				(int32(b[i+2]) << 16) +
				(int32(b[i+3]) << 24)
			new32 := int32(float64(i32) * fadeInPercent)
			b[i] = byte(new32)
			b[i+1] = byte(new32 >> 8)
			b[i+2] = byte(new32 >> 16)
			b[i+3] = byte(new32 >> 24)
			fir.toFadeIn -= 4
		}
	}
	return read, nil
}

func FadeOut(dur time.Duration, in pcm.Reader) pcm.Reader {
	perSec := in.PCMFormat().BytesPerSecond()
	bytestoFadeOut := int((time.Duration(perSec) / 1000) * (dur / time.Millisecond))

	return &fadeOutReader{
		Reader:         in,
		toFadeOut:      bytestoFadeOut,
		totaltoFadeOut: bytestoFadeOut,
	}
}

type fadeOutReader struct {
	pcm.Reader
	toFadeOut, totaltoFadeOut int
}

func (fir *fadeOutReader) ReadPCM(b []byte) (n int, err error) {
	if fir.toFadeOut == 0 {
		return fir.Reader.ReadPCM(b)
	}
	read, err := fir.Reader.ReadPCM(b)
	if err != nil {
		return read, err
	}
	format := fir.PCMFormat()
	switch format.Bits {
	case 8:
		for i, byt := range b[:read] {
			fadeOutPercent := float64(fir.toFadeOut) / float64(fir.totaltoFadeOut)
			if fadeOutPercent <= 0 {
				fadeOutPercent = 0
			}
			b[i] = byte(int8(float64(int8(byt)) * fadeOutPercent))
			fir.toFadeOut--
		}
	case 16:
		for i := 0; i+2 <= read; i += 2 {
			fadeOutPercent := float64(fir.toFadeOut) / float64(fir.totaltoFadeOut)
			if fadeOutPercent <= 0 {
				fadeOutPercent = 0
			}
			i16 := int16(b[i]) + (int16(b[i+1]) << 8)
			new16 := int16(float64(i16) * fadeOutPercent)
			b[i] = byte(new16)
			b[i+1] = byte(new16 >> 8)
			fir.toFadeOut -= 2
		}
	case 32:
		for i := 0; i+4 <= read; i += 4 {
			fadeOutPercent := float64(fir.toFadeOut) / float64(fir.totaltoFadeOut)
			if fadeOutPercent <= 0 {
				fadeOutPercent = 0
			}
			i32 := int32(b[i]) +
				(int32(b[i+1]) << 8) +
				(int32(b[i+2]) << 16) +
				(int32(b[i+3]) << 24)
			new32 := int32(float64(i32) * fadeOutPercent)
			b[i] = byte(new32)
			b[i+1] = byte(new32 >> 8)
			b[i+2] = byte(new32 >> 16)
			b[i+3] = byte(new32 >> 24)
			fir.toFadeOut -= 4
		}
	}
	return read, nil
}

var _ pcm.Reader = &fadeOutReader{}
var _ pcm.Reader = &fadeInReader{}
