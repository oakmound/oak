package audio

import "github.com/200sc/klangsynthese/font"

var (
	loaded  = make(map[string]Data)
	DefFont = font.New()
)
