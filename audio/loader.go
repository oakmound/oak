package plastic

import (
	"os"
	"path/filepath"
	//	"golang.org/x/mobile/exp/audio"
)

var (
	// Form ...main/core.go/../assets/audio,
	// the audio directory.
	wd, _ = os.Getwd()
	dir   = filepath.Join(
		filepath.Dir(wd),
		"assets",
		"audio")
)

//func loadSFX(fileName string) {
//
//}
