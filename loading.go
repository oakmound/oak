package oak

import (
	"github.com/oakmound/oak/v2/audio"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/fileutil"
	"github.com/oakmound/oak/v2/render"
	"golang.org/x/sync/errgroup"
)

var (
	startupLoadCh = make(chan bool)
	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable
)

func loadAssets(imageDir, audioDir string) {
	if conf.BatchLoad {
		dlog.Info("Loading Images")
		var eg errgroup.Group
		eg.Go(func() error {
			err := render.BatchLoad(imageDir)
			if err != nil {
				return err
			}
			dlog.Info("Done Loading Images")
			return nil
		})
		eg.Go(func() error {
			err := audio.BatchLoad(audioDir)
			if err != nil {
				return err
			}
			dlog.Info("Done Loading Audio")
			return nil
		})
		dlog.ErrorCheck(eg.Wait())
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
