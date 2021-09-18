package fileutil

import (
	"embed"
	"io"
	"testing"
)

//go:embed testdata/*
var testfs embed.FS

func TestOpen(t *testing.T) {
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
}

func TestReadFile(t *testing.T) {
	FS = testfs
	_, err := ReadFile("testdata/test.txt")
	if err != nil {
		t.Fatalf("read all failed: %v", err)
	}
}

func TestReadDir(t *testing.T) {
	FS = testfs
	ds, err := ReadDir("testdata")
	if err != nil {
		t.Fatalf("read dir failed: %v", err)
	}
	if len(ds) != 1 {
		t.Fatalf("read dir had %v elements, expected 1", len(ds))
	}
}
