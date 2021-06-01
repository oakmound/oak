package oak

import (
	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/render"
	"golang.org/x/sync/errgroup"
)

func (c *Controller) loadAssets(imageDir, audioDir string) {
	if c.config.BatchLoad {
		dlog.Info("Loading Images")
		var eg errgroup.Group
		eg.Go(func() error {
			err := render.BlankBatchLoad(imageDir, c.config.BatchLoadOptions.MaxImageFileSize)
			if err != nil {
				return err
			}
			dlog.Info("Done Loading Images")
			return nil
		})
		eg.Go(func() error {
			var err error
			if c.config.BatchLoadOptions.BlankOutAudio {
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
	c.startupLoading = false
}

// SetBinaryPayload changes how oak will load files-- instead of loading from the filesystem,
// they'll be loaded from the provided two functions: one to load bytes from a path,
// and one to list paths underneath a directory.
func SetBinaryPayload(payloadFn func(string) ([]byte, error), dirFn func(string) ([]string, error)) {
	fileutil.BindataDir = dirFn
	fileutil.BindataFn = payloadFn
}
