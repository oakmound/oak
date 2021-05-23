package fileutil

import (
	"os"
	"testing"
	"time"
)

func TestDummyFileInfo(t *testing.T) {
	dfi := dummyfileinfo{"file", false}
	if dfi.Name() != "file" {
		t.Fatalf("name mismatch")
	}
	if dfi.Size() != 0 {
		t.Fatalf("size mismatch")
	}
	if dfi.Mode() != os.ModeTemporary {
		t.Fatalf("mode mismatch")
	}
	if dfi.ModTime() != (time.Time{}) {
		t.Fatalf("modTime mismatch")
	}
	if dfi.IsDir() != false {
		t.Fatalf("isDir mismatch")
	}
	if dfi.Sys() != nil {
		t.Fatalf("sys mismatch")
	}
}
