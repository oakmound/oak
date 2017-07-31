package oak

import (
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render"
)

var (
	startupLoadCh = make(chan bool)
	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable
)

func loadAssets(imageDir, audioDir string) {
	if conf.BatchLoad {
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
	}
	endLoad()
}

func endLoad() {
	dlog.Info("Done Loading")
	startupLoadCh <- true
	dlog.Info("Startup load signal sent")
}

// SetBinaryPayload just sets some public fields on packages that require access to binary functions
// as alternatives to os file functions. This is no longer necessary, as a single package uses these
// now.
func SetBinaryPayload(payloadFn func(string) ([]byte, error), dirFn func(string) ([]string, error)) {
	fileutil.BindataDir = dirFn
	fileutil.BindataFn = payloadFn
}
