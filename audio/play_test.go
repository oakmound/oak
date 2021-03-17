package audio

import (
	"testing"
	"time"
)

func TestPlayAndLoad(t *testing.T) {
	_, err := Load("testdata", "test.wav")
	if err != nil {
		t.Fatalf("failed to load test.wav")
	}
	_, err = Load(".", "badfile.wav")
	if err == nil {
		t.Fatalf("expected loading badfile to fail")
	}
	_, err = Load(".", "play_test.go")
	if err == nil {
		t.Fatalf("expected loading non-wav file to fail")
	}
	err = Play(DefFont, "test.wav")
	if err != nil {
		t.Fatalf("failed to play test.wav (1)")
	}
	time.Sleep(1 * time.Second)
	err = DefPlay("test.wav")
	if err != nil {
		t.Fatalf("failed to play test.wav (2)")
	}
	time.Sleep(1 * time.Second)
	// Assert something was played twice
	_, err = GetSounds("test.wav")
	if err != nil {
		t.Fatalf("failed to get test.wav")
	}
	_, err = GetSounds("badfile.wav")
	if err == nil {
		t.Fatalf("expected getting badfile to fail")
	}
	Unload("test.wav")
	err = Play(DefFont, "test.wav")
	if err == nil {
		t.Fatalf("expected playing unloaded test.wav to fail")
	}
}
