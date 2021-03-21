package render

import (
	"testing"

	"github.com/oakmound/oak/v2/fileutil"
)

func Test_loadSprite(t *testing.T) {
	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset

	//loadSprite(dir, )
}
