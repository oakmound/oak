package fileutil

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"strings"
)

var (
	// FS is the filesystem that Open, ReadFile and ReadDir will query.
	FS fs.FS = os.DirFS(".")
	// FixWindowsPaths will reset all file paths loaded to replace windows style slashes
	// with unix style slashes. This is important when using io/fs or embed, because the
	// path/filepath package will produce windows style paths on a windows system, but
	// these stdlib packages will reject all windows paths.
	FixWindowsPaths = true
	// OSFallback will fallback to loading via os.Open / io.ReadFile if loading otherwise fails.
	// This is necessary when reading system level fallback fonts. Fixed paths will not be applied
	// to this fallback route.
	OSFallback = true
)

// Open is a wrapper around os.Open that will also check FS to access
// embedded data. The intended use is to use the an embedding library to create an
// Asset function that matches this signature.
func Open(file string) (io.ReadCloser, error) {
	fixedPath := fixWindowsPath(file)
	f, readErr := FS.Open(fixedPath)
	if readErr != nil && OSFallback {
		byt, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(byt)), nil
	}
	return f, readErr
}

// ReadFile replaces ioutil.ReadFile, trying to use FS.
func ReadFile(file string) ([]byte, error) {
	fixedPath := fixWindowsPath(file)
	data, readErr := fs.ReadFile(FS, fixedPath)
	if readErr != nil && OSFallback {
		return os.ReadFile(file)
	}
	return data, readErr
}

// ReadDir replaces ioutil.ReadDir, trying to use FS.
func ReadDir(file string) ([]fs.DirEntry, error) {
	fixedPath := fixWindowsPath(file)
	return fs.ReadDir(FS, fixedPath)
}

func fixWindowsPath(file string) string {
	if !FixWindowsPaths {
		return file
	}
	file = strings.Replace(file, "\\", "/", -1)
	return file
}
