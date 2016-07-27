package audio

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"errors"
	"io/ioutil"
	"path/filepath"
)

var (
	loadedWavs = make(map[string]Audio, 0)
)

// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type Audio winaudio.Audio

func InitWinAudio() {
	winaudio.InitWinAudio()
}

func GetWav(fileName string) (Audio, error) {
	if _, ok := loadedWavs[fileName]; !ok {
		return nil, errors.New("File not loaded")
	}
	return loadedWavs[fileName], nil
}

func LoadWav(directory, fileName string) (Audio, error) {
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

	files, _ := ioutil.ReadDir(baseFolder)

	for _, file := range files {
		if !file.IsDir() {
			n := file.Name()
			switch n[len(n)-4:] {
			case ".wav":
				dlog.Verb("loading file ", n)
				LoadWav(baseFolder, n)
			default:
				dlog.Error("Unsupported file ending for batchLoad: ", n)
			}
		}
	}
	return nil
}
