package fileutil

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

var (
	FS fs.FS = os.DirFS(".")
	// FixWindowsPaths will reset all file paths loaded to replace windows style slashes
	// with unix style slashes. This is important when using io/fs or embed, because the
	// path/filepath package will produce windows style paths on a windows system, but
	// these stdlib packages will reject all windows paths.
	FixWindowsPaths = true
)

// Open is a wrapper around os.Open that will also check BindataFn to access
// embedded data. The intended use is to use the an embedding library to create an
// Asset function that matches this signature.
func Open(file string) (io.ReadCloser, error) {
	file = fixWindowsPath(file)
	return FS.Open(file)
}

// ReadFile replaces ioutil.ReadFile, trying to use the BindataFn if it exists.
func ReadFile(file string) ([]byte, error) {
	file = fixWindowsPath(file)
	return fs.ReadFile(FS, file)
}

// ReadDir replaces ioutil.ReadDir, trying to use the BinaryDir if it exists.
func ReadDir(file string) ([]fs.DirEntry, error) {
	file = fixWindowsPath(file)
	return fs.ReadDir(FS, file)
}

func fixWindowsPath(file string) string {
	if !FixWindowsPaths {
		return file
	}
	file = strings.Replace(file, "\\", "/", -1)
	return file
}
