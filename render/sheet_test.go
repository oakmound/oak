package render

import (
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/v3/fileutil"
)

func TestSheetSequence(t *testing.T) {

	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset

	_, err := NewSheetSequence(nil, 10, 0)
	if err == nil {
		t.Fatalf("new sheet sequence with no sheet should fail")
	}

	sheet, err := LoadSheet(filepath.Join("assets", "images", "16", "jeremy.png"), 16, 16, 0)
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
