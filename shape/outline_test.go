package shape

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquareOutine(t *testing.T) {
	out, err := Rectangle.Outline(2, 2)
	assert.Nil(t, err)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 4)

	out, err = Rectangle.Outline(3, 3)
	assert.Nil(t, err)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 8)

	out, err = Rectangle.Outline(4, 4)
	assert.Nil(t, err)
	fmt.Println(out, err)
	assert.Equal(t, len(out), 12)

	in := JustIn(func(x, y int, sizes ...int) bool {
		return x > 5
	})
	out, err = in.Outline(10, 3)
	assert.Nil(t, err)
	assert.Equal(t, len(out), 10)
}
