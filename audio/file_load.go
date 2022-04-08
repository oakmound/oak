package audio

import (
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/oakmound/oak/v3/audio/pcm"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

// Get will read cached audio data from Load, or error if the given
// file is not in the cache.
func (c *Cache) Get(file string) (pcm.Reader, error) {
	c.mu.RLock()
	data, ok := c.data[file]
	c.mu.RUnlock()
	if !ok {
		return nil, oakerr.NotFound{InputName: file}
	}
	return data.Copy(), nil
}

// Load loads the given file and caches it by two keys:
// the full file name given and the final element of the file's
// path. If the file cannot be found or if its extension is not
// supported an error will be returned.
func (c *Cache) Load(file string) (pcm.Reader, error) {
	dlog.Verb("Loading", file)
	f, err := fileutil.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := filepath.Ext(file)
	ext = strings.ToLower(ext)
	reader, ok := LoaderForExtension(ext)
	if !ok {
		return nil, oakerr.UnsupportedFormat{Format: filepath.Ext(file)}
	}
	r, err := reader(f)
	if err != nil {
		return nil, err
	}
	c.setLoaded(file, r)
	return r, nil
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
			eg.Go(func() error {
				if blankOut {
					blankLoad(fileName)
				} else {
					_, err := DefaultCache.Load(filepath.Join(baseFolder, fileName))
					if err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	err = eg.Wait()
	return err
}

func blankLoad(filename string) {
	dlog.Verb("blank loading file %v", filename)
	DefaultCache.setLoaded(filename, &BytesReader{
		Format: pcm.Format{
			SampleRate: 44000,
			Bits:       16,
			Channels:   2,
		},
		Buffer: []byte{0, 0, 0, 0},
	})
}
