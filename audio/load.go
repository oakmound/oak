package audio

import (
	"path/filepath"

	audio "github.com/oakmound/oak/v3/audio/klang"
	"github.com/oakmound/oak/v3/audio/mp3"
	"github.com/oakmound/oak/v3/audio/wav"
	"golang.org/x/sync/errgroup"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

// Data is an alias for an interface supporting the built in filters in our
// underlying audio library
type Data audio.FullAudio

// Get will read cached audio data from Load, or error if the given
// file is not in the cache.
func (c *Cache) Get(file string) (Data, error) {
	c.mu.RLock()
	data, ok := c.data[file]
	c.mu.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: file}
	}
	return data, nil
}

// Load loads the given file and caches it by two keys:
// the full file name given and the final element of the file's
// path. If the file cannot be found or if its extension is not
// supported an error will be returned.
func (c *Cache) Load(file string) (Data, error) {
	dlog.Verb("Loading", file)
	f, err := fileutil.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buffer audio.Audio
	switch filepath.Ext(file) {
	case ".wav":
		buffer, err = wav.Load(f)
	case ".mp3":
		buffer, err = mp3.Load(f)
	default:
		return nil, oakerr.UnsupportedFormat{Format: filepath.Ext(file)}
	}
	if err != nil {
		return nil, err
	}
	data := buffer.(audio.FullAudio)
	c.setLoaded(file, data)
	return data, nil
}

// BatchLoad attempts to load all files within a given directory
// depending on their file ending
func BatchLoad(baseFolder string) error {
	return batchLoad(baseFolder, false)
}

// BlankBatchLoad acts like BatchLoad, but replaces all loaded assets
// with empty audio constructs. This is intended to reduce start-up
// times in development.
func BlankBatchLoad(baseFolder string) error {
	return batchLoad(baseFolder, true)
}

func batchLoad(baseFolder string, blankOut bool) error {
	files, err := fileutil.ReadDir(baseFolder)

	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			switch filepath.Ext(fileName) {
			case ".wav", ".mp3":
				eg.Go(func() error {
					var err error
					if blankOut {
						dlog.Verb("blank loading file")
						err = blankLoad(fileName)
					} else {
						_, err = DefaultCache.Load(filepath.Join(baseFolder, fileName))
					}
					if err != nil {
						return err
					}
					return nil
				})
			default:
				dlog.Error("Unsupported file ending for batchLoad: ", fileName)
			}
		}
	}
	err = eg.Wait()
	return err
}

func blankLoad(filename string) error {
	mformat := audio.Format{
		SampleRate: 44000,
		Bits:       16,
		Channels:   2,
	}
	buffer, err := audio.EncodeBytes(
		audio.Encoding{
			Format: mformat,
			Data:   []byte{0, 0, 0, 0},
		})
	if err != nil {
		return err
	}
	data := buffer.(audio.FullAudio)
	DefaultCache.setLoaded(filename, data)
	return nil
}
