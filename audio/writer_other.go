//go:build !windows && !linux && !darwin

package audio

import "github.com/oakmound/oak/v3/oakerr"

func initOS(driver Driver) error {
	return oakerr.UnsupportedPlatform{
		Operation: "pcm.Init",
	}
}

func newWriter(f Format) (Writer, error) {
	return nil, oakerr.UnsupportedPlatform{
		Operation: "pcm.NewWriter",
	}
}
