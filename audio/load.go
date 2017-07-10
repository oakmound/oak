package audio

import (
	"errors"
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
		sounds[i], err = Get(f)
		if err != nil {
			return nil, err
		}
	}
	return sounds, nil
}

// Get without a font will just return the raw audio data
func Get(filename string) (Data, error) {
	if IsLoaded(filename) {
		return loaded[filename], nil
	}
	return nil, oakerr.NotLoadedError{}
}

// Load joins the directory and filename, attempts to find
// the input file, and stores it as filename in the set of
// loaded files.
// This can cause a conflict when multiple files have the same
// name but different directories-- the first file loaded wil be the
// one stored in the loeaded map.
func Load(directory, filename string) (Data, error) {
	dlog.Verb("Loading", directory, filename)
	if !IsLoaded(filename) {
		f, err := fileutil.Open(filepath.Join(directory, filename))
		if err != nil {
			return nil, err
		}
		var buffer audio.Audio
		end := strings.ToLower(filename[len(filename)-4:])
		switch end {
		case ".wav":
			buffer, err = wavController.Load(f)
		case ".mp3":
			buffer, err = mp3Controller.Load(f)
		default:
			return nil, errors.New("Unsupported file ending " + end)
		}
		if err != nil {
			return nil, err
		}
		loaded[filename] = buffer.(audio.FullAudio)
	}
	return loaded[filename], nil
}

func Unload(filename string) {
	delete(loaded, filename)
}

// IsLoaded is shorthand for (if _, ok := loaded[filename]; ok)
func IsLoaded(filename string) bool {
	_, ok := loaded[filename]
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
			case ".wav", ".mp3":
				dlog.Verb("loading file ", n)
				_, err := Load(baseFolder, n)
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
