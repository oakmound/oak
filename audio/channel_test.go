package audio

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/alg/range/intrange"
)

func TestChannels(t *testing.T) {
	_, err := DefaultChannel(intrange.NewConstant(5))
	if err == nil {
		t.Fatalf("expected error calling DefChannel without file names")
	}
	_, err = Load("testdata", "test.wav")
	if err != nil {
		t.Fatalf("expected no error loading test file")
	}
	ch, err := DefaultChannel(intrange.NewLinear(1, 100), "test.wav")
	if err != nil {
		t.Fatalf("expected no error creating channel with test file")
	}
	if ch == nil {
		t.Fatalf("expected channel to be not-nil post create")
	}
	go func() {
		tm := time.Now().Add(2 * time.Second)
		// This only matters when running a suite of tests
		for time.Now().Before(tm) {
			ch <- Signal{0}
		}
	}()
	time.Sleep(2 * time.Second)
}
