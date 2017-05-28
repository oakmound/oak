package fileutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	// BindataFn is a function to access binary data outside of os.Open
	BindataFn func(string) ([]byte, error)
	wd, _     = os.Getwd()
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

// Open is a wrapper around os.Open that will also check a function to access
// byte data. The intended use is to use the go-bindata library to create an
// Asset function that matches this signature.
func Open(file string) (io.ReadCloser, error) {
	// Check bindata
	if BindataFn != nil {
		// It looks like we need to clean this output sometimes--
		// we get capitalization where we don't want it ocassionally
		rel, err := filepath.Rel(wd, file)
		if err == nil {
			data, err := BindataFn(rel)
			if err == nil {
				dlog.Verb("Found file in binary,", rel)
				// convert data to io.Reader
				return nopCloser{bytes.NewReader(data)}, err
			}
			dlog.Warn("File not found in binary", rel)
		} else {
			dlog.Warn("Error in rel", err)
		}
	}
	return os.Open(file)
}
