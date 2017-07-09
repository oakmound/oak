package audio

import (
	"path/filepath"
	"strings"

	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/mp3"
	"github.com/200sc/klangsynthese/wav"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/fileutil"
	"bitbucket.org/oakmoundstudio/oak/oakerr"
)

var (
	wavController = wav.NewController()
	mp3Controller = mp3.NewController()
)

type Data audio.FullAudio

// GetSounds returns a set of Data for a set of input filenames
func GetSounds(fileNames ...string) ([]Data, error) {
	var err error
	sounds := make([]Data, len(fileNames))
	for i, f := range fileNames {
		sounds[i], err = GetWav(f)
		if err != nil {
			return nil, err
		}
	}
	return sounds, nil
}

// GetWav without a font will just return the raw audio data
func GetWav(filename string) (Data, error) {
	if IsLoaded(filename) {
		return loadedWavs[filename], nil
	}
	return nil, oakerr.NotLoadedError{}
}

// LoadWav joins the directory and filename, attempts to find
// the input file, and stores it as filename in the set of
// loadedWav files.
// This can cause a conflict when multiple files have the same
// name but different directories-- the first file loaded wil be the
// one stored in the loeaded map.
func LoadWav(directory, filename string) (Data, error) {
	dlog.Verb("Loading", directory, filename)
	if !IsLoaded(filename) {
		f, err := fileutil.Open(filepath.Join(directory, filename))
		if err != nil {
			return nil, err
		}
		buffer, err := wavController.Load(f)
		if err != nil {
			return nil, err
		}
		loadedWavs[filename] = buffer.(audio.FullAudio)
	}
	return loadedWavs[filename], nil
}

// IsLoaded is shorthand for (if _, ok := loadedWavs[filename]; ok)
func IsLoaded(filename string) bool {
	_, ok := loadedWavs[filename]
	return ok
}

// BatchLoad attempts to load all files within a given directory
// depending on their file ending (currently supporting .wav only)
func BatchLoad(baseFolder string) error {

	files, err := fileutil.ReadDir(baseFolder)

	if err != nil {
		dlog.Error(err)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			n := file.Name()
			dlog.Verb(n)
			switch strings.ToLower(n[len(n)-4:]) {
			case ".wav":
				dlog.Verb("loading file ", n)
				_, err := LoadWav(baseFolder, n)
				if err != nil {
					dlog.Error(err)
					return err
				}
			default:
				dlog.Error("Unsupported file ending for batchLoad: ", n)
			}
		}
	}
	dlog.Verb("Loading complete")
	return nil
}
