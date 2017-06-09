package oak

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/audio"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/fileutil"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	startupLoadCh = make(chan bool)
	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable
)

func loadAssets(imageDir, audioDir string) {
	dlog.Info("Loading Images")
	err := render.BatchLoad(imageDir)
	if err != nil {
		dlog.Error(err)
		endLoad()
		return
	}
	dlog.Info("Done Loading Images")
	dlog.Info("Loading Audio")
	err = audio.BatchLoad(audioDir)
	if err != nil {
		dlog.Error(err)
		endLoad()
		return
	}
	dlog.Info("Done Loading Audio")
	endLoad()
}

func endLoad() {
	dlog.Info("Done Loading")
	startupLoadCh <- true
	dlog.Info("Startup load signal sent")
}

func SetBinaryPayload(payloadFn func(string) ([]byte, error), dirFn func(string) ([]string, error)) {
	winaudio.BindataFn = payloadFn
	fileutil.BindataDir = dirFn
	fileutil.BindataFn = payloadFn
}
