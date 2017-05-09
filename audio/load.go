package audio

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/oakerr"
)

func (f *Font) GetWav(filename string) (*Audio, error) {
	if IsLoaded(filename) {
		return &Audio{loadedWavs[filename], f, nil, nil}, nil
	}
	return nil, oakerr.NotLoadedError{}
}

func (f *Font) LoadWav(directory, filename string) (*Audio, error) {
	ad, err := LoadWav(directory, filename)
	if err != nil {
		return nil, err
	}
	return &Audio{ad, f, nil, nil}, nil
}

func GetSounds(fileNames ...string) ([]AudioData, error) {
	var err error
	sounds := make([]AudioData, len(fileNames))
	for i, f := range fileNames {
		sounds[i], err = GetWav(f)
		if err != nil {
			return nil, err
		}
	}
	return sounds, nil
}

func GetWav(filename string) (AudioData, error) {
	if IsLoaded(filename) {
		return loadedWavs[filename], nil
	}
	return nil, oakerr.NotLoadedError{}
}

func LoadWav(directory, filename string) (AudioData, error) {
	dlog.Verb("Loading", directory, filename)
	if _, ok := loadedWavs[filename]; !ok {
		buffer, err := winaudio.LoadWav(filepath.Join(directory, filename))
		if err != nil {
			return buffer, err
		}
		loadedWavs[filename] = buffer
	}
	return loadedWavs[filename], nil
}

func IsLoaded(filename string) bool {
	_, ok := loadedWavs[filename]
	return ok
}

func BatchLoad(baseFolder string) error {

	files, err := ioutil.ReadDir(baseFolder)

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
