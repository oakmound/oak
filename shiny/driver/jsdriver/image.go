//go:build js
// +build js

package jsdriver

import (
	"github.com/oakmound/oak/v3/driver/common"
)

type imageImpl struct {
	common.Image
	screen *screenImpl
}
