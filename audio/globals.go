package audio

import (
	"sync"

	"github.com/oakmound/oak/v3/audio/font"
)

var (
	loadedLock sync.RWMutex
	loaded     = make(map[string]Data)
	// DefaultFont is the font used for default functions. It can be publicly
	// modified to apply a default font to generated audios through def
	// methods. If it is not modified, it is a font of zero filters.
	DefaultFont = font.New()
)
