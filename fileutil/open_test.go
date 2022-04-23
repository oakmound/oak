package fileutil

import (
	"embed"
	"errors"
	"io"
	"os"
	"testing"
)

//go:embed testdata/*
var testfs embed.FS

func TestOpen(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		FS = testfs
		f, err := Open("testdata/test.txt")
		if err != nil {
			t.Fatalf("open failed: %v", err)
		}
		_, err = io.ReadAll(f)
		if err != nil {
			t.Fatalf("read all failed: %v", err)
		}
		err = f.Close()
		if err != nil {
			t.Fatalf("close failed: %v", err)
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		FS = testfs
		_, err := Open("testdata/notfound.txt")
		perr := &os.PathError{}
		if !errors.As(err, &perr) {
			t.Fatalf("expected path error: %v", err)
		}
	})
	t.Run("OSFallback", func(t *testing.T) {
		FS = testfs
		f, err := os.CreateTemp(".", "test")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		f.Close()
		f2, err := Open(f.Name())
		if err != nil {
			t.Fatalf("open failed: %v", err)
		}
		f2.Close()
	})
}

func TestReadFile(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		FS = testfs
		_, err := ReadFile("testdata/test.txt")
		if err != nil {
			t.Fatalf("read all failed: %v", err)
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		FS = testfs
		_, err := ReadFile("testdata/notfound.txt")
		perr := &os.PathError{}
		if !errors.As(err, &perr) {
			t.Fatalf("expected path error: %v", err)
		}
	})
}

func TestReadDir(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		FS = testfs
		ds, err := ReadDir("testdata")
		if err != nil {
			t.Fatalf("read dir failed: %v", err)
		}
		if len(ds) != 1 {
			t.Fatalf("read dir had %v elements, expected 1", len(ds))
		}
	})
	t.Run("NoWindowsPaths", func(t *testing.T) {
		FixWindowsPaths = false
		defer func() {
			FixWindowsPaths = true
		}()
		FS = testfs
		ds, err := ReadDir("testdata")
		if err != nil {
			t.Fatalf("read dir failed: %v", err)
		}
		if len(ds) != 1 {
			t.Fatalf("read dir had %v elements, expected 1", len(ds))
		}
	})
}
