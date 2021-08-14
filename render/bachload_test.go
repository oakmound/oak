package render

import (
	"errors"
	"testing"

	"github.com/oakmound/oak/v3/oakerr"
)

func TestBlankBatchLoad_BadBaseFolder(t *testing.T) {
	err := BlankBatchLoad("notfound", 0)
	expected := oakerr.InvalidInput{}
	if !errors.As(err, &expected) {
		t.Fatalf("error was not expected invalid input: %v", err)
	}
}

func TestBatchLoad(t *testing.T) {
	if BatchLoad("testdata/assets/images") != nil {
		t.Fatalf("batch load failed")
	}
	sh, err := GetSheet("jeremy.png")
	if err != nil {
		t.Fatalf("get sheet failed: %v", err)
	}
	if len(sh.ToSprites()) != 8 {
		t.Fatalf("sheet did not contain 8 sprites")
	}
	_, err = LoadSprite("dir/dummy.jpg")
	if err == nil {
		t.Fatalf("load sprite should have failed")
	}
	sp, err := GetSprite("dummy.gif")
	if sp != nil {
		t.Fatalf("get sprite should be nil")
	}
	if err == nil {
		t.Fatalf("get sprite should have failed")
	}
	sp, err = GetSprite("jeremy.png")
	if sp == nil {
		t.Fatalf("get sprite failed")
	}
	if err != nil {
		t.Fatalf("get sprite failed: %v", err)
	}
	DefaultCache.ClearAll()
}
