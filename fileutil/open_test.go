package fileutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopCloser(t *testing.T) {
	assert.Nil(t, nopCloser{}.Close())
}

func TestOpen(t *testing.T) {
	_, err := Open("notafile")
	assert.NotNil(t, err)
	BindataFn = func(s string) ([]byte, error) {
		if s == "exists" {
			return []byte{0}, nil
		}
		return []byte{}, errors.New("Doesn't Exist")
	}
	_, err = Open("exists")
	assert.Nil(t, err)
	_, err = Open("doesntexist")
	assert.NotNil(t, err)
	BindataFn = nil
}

func TestReadFile(t *testing.T) {
	_, err := ReadFile("notafile")
	assert.NotNil(t, err)
	BindataFn = func(s string) ([]byte, error) {
		if s == "exists" {
			return []byte{0}, nil
		}
		return []byte{}, errors.New("Doesn't Exist")
	}
	_, err = ReadFile("exists")
	assert.Nil(t, err)
	_, err = ReadFile("doesntexist")
	assert.NotNil(t, err)
	BindataFn = nil
}

func TestReadDir(t *testing.T) {
	_, err := ReadDir("notafile")
	assert.NotNil(t, err)
	BindataDir = func(s string) ([]string, error) {
		if s == "exists" {
			return []string{""}, nil
		}
		return []string{}, errors.New("Doesn't Exist")
	}
	_, err = ReadDir("exists")
	assert.Nil(t, err)
	_, err = ReadDir("doesntexist")
	assert.NotNil(t, err)
	BindataDir = nil
}
