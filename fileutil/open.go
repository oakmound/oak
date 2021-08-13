package fileutil

import (
	"io"
	"io/fs"
	"os"
)

var (
	FS fs.FS = os.DirFS(".")
)

// Open is a wrapper around os.Open that will also check BindataFn to access
// embedded data. The intended use is to use the an embedding library to create an
// Asset function that matches this signature.
func Open(file string) (io.ReadCloser, error) {
	return FS.Open(file)
}

// ReadFile replaces ioutil.ReadFile, trying to use the BindataFn if it exists.
func ReadFile(file string) ([]byte, error) {
	return fs.ReadFile(FS, file)
}

// ReadDir replaces ioutil.ReadDir, trying to use the BinaryDir if it exists.
func ReadDir(file string) ([]fs.DirEntry, error) {
	return fs.ReadDir(FS, file)
}
