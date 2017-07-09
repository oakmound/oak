package audio

import "github.com/200sc/klangsynthese/font"

var (
	loadedWavs = make(map[string]Data, 0)
	DefFont    = font.New()
)
