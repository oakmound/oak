package audio

import (
	"io"
	"sync"

	"github.com/oakmound/oak/v3/audio/format/flac"
	"github.com/oakmound/oak/v3/audio/format/mp3"
	"github.com/oakmound/oak/v3/audio/format/wav"
	"github.com/oakmound/oak/v3/audio/pcm"
)

type fileLoader func(r io.Reader) (pcm.Reader, error)

var fileLoadersLock sync.RWMutex
var fileLoaders = map[string]func(r io.Reader) (pcm.Reader, error){
	"mp3":  mp3.Load,
	"wav":  wav.Load,
	"flac": flac.Load,
}

func RegisterFormat(extension string, fn fileLoader) {
	fileLoadersLock.Lock()
	fileLoaders[extension] = fn
	fileLoadersLock.Unlock()
}

func LoaderForExtension(extension string) (fileLoader, bool) {
	fileLoadersLock.RLock()
	defer fileLoadersLock.RUnlock()
	loader, ok := fileLoaders[extension]
	return loader, ok
}
