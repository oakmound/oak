package oak

import (
	"bytes"
	"strconv"
	"testing"
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

	for i, cas := range cases {
		cas := cas
		t.Run("case"+strconv.Itoa(i), func(t *testing.T) {
			r := bytes.NewReader([]byte(cas.input))
			_, err := LoadKeyBindings(r)
			if cas.shouldSucceed {
				if err != nil {
					t.Fatalf("case failed: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("case should have failed")
				}
			}
		})
	}
}
