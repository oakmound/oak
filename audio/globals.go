package audio

import "github.com/200sc/klangsynthese/font"

var (
	loaded  = make(map[string]Data, 0)
	DefFont = font.New()
)
