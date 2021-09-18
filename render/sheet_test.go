package render

import (
	"embed"
	"os"
	"testing"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/fileutil"
)

//go:embed testdata/assets/*
var testfs embed.FS

func TestMain(m *testing.M) {
	fileutil.FS = testfs
	os.Exit(m.Run())
}

func TestSheetSequence(t *testing.T) {
	_, err := NewSheetSequence(nil, 10, 0)
	if err == nil {
		t.Fatalf("new sheet sequence with no sheet should fail")
	}

	sheet, err := LoadSheet("testdata/assets/images/16x16/jeremy.png", intgeom.Point2{16, 16})
	if err != nil {
		t.Fatalf("loading jeremy sheet should not fail")
	}
	_, err = NewSheetSequence(sheet, 10, 0, 1, 0, 2)
	if err != nil {
		t.Fatalf("creating jeremy sheet sequence should not fail")
	}

	_, err = NewSheetSequence(sheet, 10, 100, 1, 0, 2)
	if err == nil {
		t.Fatalf("creating jeremy sheet sequence with invalid frames should fail")
	}

	_, err = NewSheetSequence(sheet, 10, 1, 100)
	if err == nil {
		t.Fatalf("creating jeremy sheet sequence with invalid frames should fail")
	}
}
