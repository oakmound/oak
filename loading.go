package oak

import (
	"github.com/oakmound/oak/v2/audio"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/fileutil"
	"github.com/oakmound/oak/v2/render"
	"golang.org/x/sync/errgroup"
)

func (c *Controller) loadAssets(imageDir, audioDir string) {
	if conf.BatchLoad {
		dlog.Info("Loading Images")
		var eg errgroup.Group
		eg.Go(func() error {
			err := render.BlankBatchLoad(imageDir, conf.BatchLoadOptions.MaxImageFileSize)
			if err != nil {
				return err
			}
			dlog.Info("Done Loading Images")
			return nil
		})
		eg.Go(func() error {
			var err error
			if conf.BatchLoadOptions.BlankOutAudio {
				err = audio.BlankBatchLoad(audioDir)
			} else {
				err = audio.BatchLoad(audioDir)
			}
			if err != nil {
				return err
			}
			dlog.Info("Done Loading Audio")
			return nil
		})
		dlog.ErrorCheck(eg.Wait())
	}
	c.endLoad()
}

func (c *Controller) endLoad() {
	dlog.Info("Done Loading")
	close(c.startupLoadCh)
	dlog.Info("Startup load closed")
}

// SetBinaryPayload just sets some public fields on packages that require access to binary functions
// as alternatives to os file functions. This is no longer necessary, as a single package uses these
// now.
func SetBinaryPayload(payloadFn func(string) ([]byte, error), dirFn func(string) ([]string, error)) {
	fileutil.BindataDir = dirFn
	fileutil.BindataFn = payloadFn
}
