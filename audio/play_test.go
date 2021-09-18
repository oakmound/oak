package audio

import (
	"path/filepath"
	"testing"
	"time"
)

func TestPlayAndLoad(t *testing.T) {
	_, err := Load(filepath.Join("testdata", "test.wav"))
	if err != nil {
		t.Fatalf("failed to load test.wav")
	}
	_, err = Load("badfile.wav")
	if err == nil {
		t.Fatalf("expected loading badfile to fail")
	}
	_, err = Load("play_test.go")
	if err == nil {
		t.Fatalf("expected loading non-wav file to fail")
	}
	err = Play(DefaultFont, "test.wav")
	if err != nil {
		t.Fatalf("failed to play test.wav (1)")
	}
	time.Sleep(1 * time.Second)
	err = DefaultPlay("test.wav")
	if err != nil {
		t.Fatalf("failed to play test.wav (2)")
	}
	time.Sleep(1 * time.Second)
	// Assert something was played twice
	DefaultCache.Clear("test.wav")
	err = Play(DefaultFont, "test.wav")
	if err == nil {
		t.Fatalf("expected playing unloaded test.wav to fail")
	}
}
