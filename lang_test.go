package oak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLang(t *testing.T) {
	SetLang("Gibberish")
	assert.Equal(t, Lang, ENGLISH)
	SetLang("German")
	assert.Equal(t, Lang, GERMAN)
	SetLang("English")
	assert.Equal(t, Lang, ENGLISH)
}
