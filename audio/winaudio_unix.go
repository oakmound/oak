//+build unix

package audio

// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type Audio interface{}

func InitWinAudio() {
}

func GetSounds(fileNames ...string) ([]Audio, error) {
	return nil, nil
}

func GetWav(fileName string) (Audio, error) {
	return nil, nil
}

func PlayWav(fileName string) error {
	return nil
}

func LoadWav(directory, fileName string) (Audio, error) {
	return nil, nil
}

func BatchLoad(baseFolder string) error {
	return nil
}
