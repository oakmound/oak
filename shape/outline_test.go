package shape

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquareOutine(t *testing.T) {
	out, err := Rectangle.Outline(2, 2)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 4)

	out, err = Rectangle.Outline(3, 3)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 8)

	out, err = Rectangle.Outline(4, 4)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 12)
}
