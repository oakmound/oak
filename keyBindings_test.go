package oak

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadKeyBindings(t *testing.T) {

	type testCase struct {
		input         string
		shouldSucceed bool
	}

	cases := []testCase{
		{
			"MoveUp = \"W\"\nMoveDown = \"S\"\nAttack = \"Spacebar\"\n", true,
		},
		{
			"MoveUp\"W\"\nMoveDown = \"S\"\nAttack = \"Spacebar\"\n", false,
		},
		{
			"MoveUp = \"W\"\n\nMoveDown = \"S\"\nAttack = \"Spacebar\"\n", true,
		},
	}

	for _, cas := range cases {
		r := bytes.NewReader([]byte(cas.input))
		_, err := LoadKeyBindings(r)
		if cas.shouldSucceed {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
