package alg

import (
	"testing"

	"bitbucket.org/oakmoundstudio/oak/tests"
)

func TestIntRanges(t *testing.T) {
	_, err := NewLinearIntRange(1, 0)

	tests.ExpectError(err, t)

	_, err = NewBaseSpreadIntRange(1, -1)

	tests.ExpectError(err, t)
}
