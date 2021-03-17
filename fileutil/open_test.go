package fileutil

import (
	"errors"
	"testing"
)

func TestNopCloser(t *testing.T) {
	if (nopCloser{}).Close() != nil {
		t.Fatalf("no op closer failed to close")
	}
}

func TestOpen(t *testing.T) {
	_, err := Open("notafile")
	if err == nil {
		t.Fatalf("expected open to fail")
	}
	BindataFn = func(s string) ([]byte, error) {
		if s == "exists" {
			return []byte{0}, nil
		}
		return []byte{}, errors.New("Doesn't Exist")
	}
	_, err = Open("exists")
	if err != nil {
		t.Fatalf("expected open (exists) to pass")
	}
	_, err = Open("doesntexist")
	if err == nil {
		t.Fatalf("expected open (doesntexist) to fail")
	}
	BindataFn = nil
}

func TestReadFile(t *testing.T) {
	_, err := ReadFile("notafile")
	if err == nil {
		t.Fatalf("expected read to fail")
	}
	BindataFn = func(s string) ([]byte, error) {
		if s == "exists" {
			return []byte{0}, nil
		}
		return []byte{}, errors.New("Doesn't Exist")
	}
	_, err = ReadFile("exists")
	if err != nil {
		t.Fatalf("expected read (exists) to pass")
	}
	_, err = ReadFile("doesntexist")
	if err == nil {
		t.Fatalf("expected read (doesntexist) to fail")
	}
	BindataFn = nil
}

func TestReadDir(t *testing.T) {
	_, err := ReadDir("notafile")
	if err == nil {
		t.Fatalf("expected read dir to fail")
	}
	BindataDir = func(s string) ([]string, error) {
		if s == "exists" {
			return []string{""}, nil
		}
		return []string{}, errors.New("Doesn't Exist")
	}
	_, err = ReadDir("exists")
	if err != nil {
		t.Fatalf("expected read dir (exists) to pass")
	}
	_, err = ReadDir("doesntexist")
	if err == nil {
		t.Fatalf("expected read dir (doesntexist) to fail")
	}
	BindataDir = nil
}
