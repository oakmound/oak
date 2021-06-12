package oak

import (
	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/render"
	"golang.org/x/sync/errgroup"
)

func (w *Window) loadAssets(imageDir, audioDir string) {
	if w.config.BatchLoad {
		var eg errgroup.Group
		eg.Go(func() error {
			err := render.BlankBatchLoad(imageDir, w.config.BatchLoadOptions.MaxImageFileSize)
			if err != nil {
				return err
			}
			dlog.Verb("Done Loading Images")
			return nil
		})
		eg.Go(func() error {
			var err error
			if w.config.BatchLoadOptions.BlankOutAudio {
				err = audio.BlankBatchLoad(audioDir)
			} else {
				err = audio.BatchLoad(audioDir)
			}
			if err != nil {
				return err
			}
			dlog.Verb("Done Loading Audio")
			return nil
		})
		dlog.ErrorCheck(eg.Wait())
	}
	w.endLoad()
}

func (w *Window) endLoad() {
	dlog.Verb("Done Loading")
	w.startupLoading = false
}

// SetBinaryPayload changes how oak will load files-- instead of loading from the filesystem,
// they'll be loaded from the provided two functions: one to load bytes from a path,
// and one to list paths underneath a directory.
func SetBinaryPayload(payloadFn func(string) ([]byte, error), dirFn func(string) ([]string, error)) {
	fileutil.BindataDir = dirFn
	fileutil.BindataFn = payloadFn
}
