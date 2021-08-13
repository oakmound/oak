package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

var (
	imgPath1 = "16x16/jeremy.png"
)

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
	_, err = DefaultCache.loadSprite("dir/dummy.jpg", 0)
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

func TestSetAssetPath(t *testing.T) {
	_, err := LoadSheet("testdata/assets/images/"+imgPath1, intgeom.Point2{16, 16})
	if err != nil {
		t.Fatalf("load sheet 1 failed: %v", err)
	}
	DefaultCache.ClearAll()
	wd, _ := os.Getwd()
	_, err = LoadSheet(filepath.Join(wd, imgPath1), intgeom.Point2{16, 16})
	if err == nil {
		t.Fatalf("load sheet should have failed")
	}
	DefaultCache.ClearAll()
}

func TestBadSheetParams(t *testing.T) {
	_, err := LoadSheet(filepath.Join("assets", "images", imgPath1), intgeom.Point2{0, 16})
	if err == nil {
		t.Fatalf("load sheet should have failed")
	}
	_, err = LoadSheet(filepath.Join("assets", "images", imgPath1), intgeom.Point2{16, 0})
	if err == nil {
		t.Fatalf("load sheet should have failed")
	}
}

func TestSheetStorage(t *testing.T) {
	if _, err := GetSheet("jeremy.png"); err == nil {
		t.Fatalf("sheets should not be loaded at startup")
	}
	_, err := LoadSheet("testdata/assets/images/"+imgPath1, intgeom.Point2{16, 16})
	if err != nil {
		t.Fatalf("load sheet failed: %v", err)
	}
	if _, err := GetSheet("jeremy.png"); err != nil {
		t.Fatalf("sheet did not load: %v", err)
	}
	DefaultCache.ClearAll()
}
