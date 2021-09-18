package oak

import (
	"io/fs"

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

// SetFS updates all calls oak or oak's subpackages will make to read from the given filesystem.
// By default, this is set to os.DirFS(".")
func SetFS(filesystem fs.FS) {
	fileutil.FS = filesystem
}
