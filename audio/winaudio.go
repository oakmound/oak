package audio

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	loadedWavs = make(map[string]Audio, 0)
)

// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type Audio winaudio.Audio

func InitWinAudio() {
	err := winaudio.InitWinAudio()
	if err != nil {
		panic(err)
	}
}

func GetSounds(fileNames ...string) ([]Audio, error) {
	var err error
	sounds := make([]Audio, len(fileNames))
	for i, f := range fileNames {
		sounds[i], err = GetWav(f)
		if err != nil {
			return nil, err
		}
	}
	return sounds, nil
}

func GetWav(fileName string) (Audio, error) {
	if _, ok := loadedWavs[fileName]; !ok {
		return nil, errors.New("File not loaded")
	}
	return loadedWavs[fileName], nil
}

func PlayWav(fileName string) error {
	a, err := GetWav(fileName)
	if err == nil {
		err = a.Play()
	} else {
		dlog.Error(err)
	}
	return err
}

func LoadWav(directory, fileName string) (Audio, error) {
	dlog.Verb("Loading", directory, fileName)
	if _, ok := loadedWavs[fileName]; !ok {
		buffer, err := winaudio.LoadWav(filepath.Join(directory, fileName))
		if err != nil {
			return buffer, err
		}
		loadedWavs[fileName] = buffer
	}
	return loadedWavs[fileName], nil
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
