package render

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/oakerr"
)

var (
	imgPath1 = "16x16/jeremy.png"
)

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

func TestMakeSheet_BadDimensions(t *testing.T) {
	_, err := MakeSheet(nil, intgeom.Point2{0, 5})
	expected := &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "cellSize.X" {
		t.Fatalf("expected invalid width, got %v", expected.InputName)
	}

	_, err = MakeSheet(nil, intgeom.Point2{5, -1})
	expected = &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "cellSize.Y" {
		t.Fatalf("expected invalid height, got %v", expected.InputName)
	}
}

func TestMakeSheet_BadMod(t *testing.T) {
	rgba := &image.RGBA{
		Rect: image.Rect(0, 0, 10, 10),
	}
	_, err := MakeSheet(rgba, intgeom.Point2{4, 5})
	expected := &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "cellSize.X" {
		t.Fatalf("expected invalid width, got %v", expected.InputName)
	}

	_, err = MakeSheet(rgba, intgeom.Point2{5, 4})
	expected = &oakerr.InvalidInput{}
	if !errors.As(err, expected) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if expected.InputName != "cellSize.Y" {
		t.Fatalf("expected invalid height, got %v", expected.InputName)
	}
}
