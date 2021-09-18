package render

import (
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/fileutil"
	"github.com/oakmound/oak/v3/oakerr"
)

// BatchLoad loads subdirectories from the given base folder and imports all files,
// using alias rules to automatically determine the size of sprites and sheets in
// subfolders.
func BatchLoad(baseFolder string) error {
	return BlankBatchLoad(baseFolder, 0)
}

// BlankBatchLoad acts like BatchLoad, but will not load and instead return a blank image
// of the appropriate dimensions for anything above maxFileSize.
func BlankBatchLoad(baseFolder string, maxFileSize int64) error {
	var wg sync.WaitGroup
	err := fs.WalkDir(fileutil.FS, baseFolder, func(file string, d fs.DirEntry, err error) error {
		if d == nil {
			// We've been given a bad base directory
			return oakerr.InvalidInput{InputName: "baseFolder"}
		}
		if d.IsDir() {
			return nil
		}
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			ext := filepath.Ext(file)
			if _, ok := fileDecoders[ext]; !ok {
				// Ignore files we know we can't parse
				return
			}
			_, err := DefaultCache.loadSprite(file, maxFileSize)
			if err != nil {
				dlog.Error(err)
				return
			}
			if cell, ok := shouldLoadSheet(file); ok {
				_, err = DefaultCache.LoadSheet(file, cell)
				if err != nil {
					dlog.Error(err)
					return
				}
			}
		}(file)
		return nil
	})
	wg.Wait()
	return err
}

var (
	sheetFileRegex      = regexp.MustCompile(`^[^\d]*(\d+)x(\d+)\..*$`)
	sheetDirectoryRegex = regexp.MustCompile(`^[^\d]*(\d+)x(\d+)$`)
)

func shouldLoadSheet(file string) (intgeom.Point2, bool) {
	// when should we determine a file should be loaded as a sheet?
	// 1. If the file itself ends in that syntax: image_%dx%d.png
	// 2. If the file's final directory ends in special syntax: (%dx%d)
	// ... preferring the former

	matches := sheetFileRegex.FindStringSubmatch(file)
	if len(matches) != 0 {
		// if len matches != 3, or if these fail to parse, our regex is wrong
		w, _ := strconv.Atoi(matches[1])
		h, _ := strconv.Atoi(matches[2])
		return intgeom.Point2{w, h}, true
	}
	dirMatches := sheetDirectoryRegex.FindStringSubmatch(path.Dir(file))
	if len(dirMatches) != 0 {
		w, _ := strconv.Atoi(dirMatches[1])
		h, _ := strconv.Atoi(dirMatches[2])
		return intgeom.Point2{w, h}, true
	}
	return intgeom.Point2{}, false
}
