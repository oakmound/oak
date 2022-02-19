//go:build (!windows && !linux)

package pcm

import "github.com/oakmound/oak/v3/oakerr"

func initOS() error {
	return oakerr.UnsupportedPlatform{
		Operation: "pcm.Init",
	}
}

func newWriter(f Format) (Writer, error) {
	return nil, oakerr.UnsupportedPlatform{
		Operation: "pcm.NewWriter",
	}
}
