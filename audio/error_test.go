package audio

import (
	"testing"

	"github.com/oakmound/oak/v2/oakerr"
)

func TestErrorChannel(t *testing.T) {
	err := oakerr.ExistingElement{}
	err2 := <-errChannel(err)
	if err != err2 {
		t.Fatalf("err channel did not propagate error")
	}
}
