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
		dlog.Verb("Done Loading Audio")
		return err
	})
	dlog.ErrorCheck(eg.Wait())
}

func (w *Window) endLoad() {
	dlog.Verb("Done Loading")
	w.NextScene()
}

// SetFS updates all calls oak or oak's subpackages will make to read from the given filesystem.
// By default, this is set to os.DirFS(".")
func SetFS(filesystem fs.FS) {
	fileutil.FS = filesystem
}
