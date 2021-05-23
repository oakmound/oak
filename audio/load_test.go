package audio

import "testing"

func TestBatchLoad(t *testing.T) {
	err := BatchLoad("testdata")
	if err != nil {
		t.Fatalf("expected batchload on valid path to succeed")
	}
	err = BlankBatchLoad("testdata")
	if err != nil {
		t.Fatalf("expected blank batchload on valid path to succeed: %v", err)
	}
	err = BatchLoad("GarbagePath")
	if err == nil {
		t.Fatalf("expected batchload on nonexistant path to fail")
	}
}
