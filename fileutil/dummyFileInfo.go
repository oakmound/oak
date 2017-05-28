package fileutil

import (
	"os"
	"time"
)

type dummyfileinfo struct {
	file  string
	isdir bool
}

func (dfi dummyfileinfo) Name() string {
	return dfi.file
}

func (dummyfileinfo) Size() int64 {
	return 0
}

func (dummyfileinfo) Mode() os.FileMode {
	return os.ModeTemporary
}

func (dummyfileinfo) ModTime() time.Time {
	return time.Time{}
}

func (dfi dummyfileinfo) IsDir() bool {
	return dfi.isdir
}

func (dummyfileinfo) Sys() interface{} {
	return nil
}
