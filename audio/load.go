package audio

import (
	"path/filepath"
	"strings"

	audio "github.com/oakmound/oak/v3/audio/klang"
	"github.com/oakmound/oak/v3/audio/mp3"
	"github.com/oakmound/oak/v3/audio/wav"
	"golang.org/x/sync/errgroup"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

// Data is an alias for an interface supporting the built in filters in our
// external audio playing library
type Data audio.FullAudio

// GetSounds returns a set of Data for a set of input filenames
func GetSounds(fileNames ...string) ([]Data, error) {
	var err error
	sounds := make([]Data, len(fileNames))
	for i, f := range fileNames {
		sounds[i], err = Get(f)
		if err != nil {
			return nil, err
		}
	}
	if len(sounds) == 0 {
		return sounds, oakerr.InsufficientInputs{AtLeast: 1, InputName: "fileName"}
	}
	return sounds, nil
}

// Get without a font will just return the raw audio data
func Get(filename string) (Data, error) {
	if IsLoaded(filename) {
		return loaded[filename], nil
	}
	return nil, oakerr.NotFound{InputName: filename}
}

// Load joins the directory and filename, attempts to find
// the input file, and stores it as filename in the set of
// loaded files.
// This can cause a conflict when multiple files have the same
// name but different directories-- the first file loaded wil be the
// one stored in the loeaded map.
func Load(directory, filename string) (Data, error) {
	dlog.Verb("Loading", directory, filename)
	if data, ok := getLoaded(filename); ok {
		return data, nil
	}
	f, err := fileutil.Open(filepath.Join(directory, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buffer audio.Audio
	end := strings.ToLower(filename[len(filename)-4:])
	switch end {
	case ".wav":
		buffer, err = wav.Load(f)
	case ".mp3":
		buffer, err = mp3.Load(f)
	default:
		return nil, oakerr.UnsupportedFormat{Format: end}
	}
	if err != nil {
		return nil, err
	}
	data := buffer.(audio.FullAudio)
	setLoaded(filename, data)
	return data, nil
}

// Unload removes an element from the loaded map. If the element does not
// exist, it does nothing.
func Unload(filename string) {
	loadedLock.Lock()
	delete(loaded, filename)
	loadedLock.Unlock()
}

// IsLoaded is shorthand for (if _, ok := loaded[filename]; ok)
func IsLoaded(filename string) bool {
	loadedLock.RLock()
	_, ok := loaded[filename]
	loadedLock.RUnlock()
	return ok
}

func getLoaded(filename string) (Data, bool) {
	loadedLock.RLock()
	data, ok := loaded[filename]
	loadedLock.RUnlock()
	return data, ok
}

func setLoaded(filename string, data Data) {
	loadedLock.Lock()
	loaded[filename] = data
	loadedLock.Unlock()
}

// BatchLoad attempts to load all files within a given directory
// depending on their file ending (currently supporting .wav and .mp3)
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
			n := file.Name()
			dlog.Verb(n)
			switch strings.ToLower(n[len(n)-4:]) {
			case ".wav", ".mp3":
				dlog.Verb("loading file ", n)
				eg.Go(func() error {
					var err error
					if blankOut {
						dlog.Verb("blank loading file")
						err = blankLoad(n)
					} else {
						_, err = Load(baseFolder, n)
					}
					if err != nil {
						return err
					}
					return nil
				})
			default:
				dlog.Error("Unsupported file ending for batchLoad: ", n)
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
	setLoaded(filename, data)
	return nil
}
