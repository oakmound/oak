package scene

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBooleanLoop(t *testing.T) {
	v := true
	l := BooleanLoop(&v)
	require.True(t, l())
	v = false
	require.False(t, l())
	require.True(t, l())
}

func randString() string {
	length := rand.Intn(100)
	data := make([]byte, length)
	for i := range data {
		data[i] = byte(rand.Intn(255))
	}
	return string(data)
}

func TestGoTo(t *testing.T) {
	tests := 10
	for i := 0; i < tests; i++ {
		s := randString()
		gt := GoTo(s)
		s2, result := gt()
		require.Equal(t, s, s2)
		require.Nil(t, result)
	}
}

func TestGoToPtr(t *testing.T) {
	tests := 10
	s := new(string)
	gt := GoToPtr(s)
	for i := 0; i < tests; i++ {
		*s = randString()
		s2, result := gt()
		require.Equal(t, *s, s2)
		require.Nil(t, result)
	}
}

func TestGoToPtrNil(t *testing.T) {
	s, _ := GoToPtr(nil)()
	require.Equal(t, "", s)
}
