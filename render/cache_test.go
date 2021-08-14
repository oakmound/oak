package render

import "testing"

func TestCache_Clear(t *testing.T) {
	err := BatchLoad("testdata/assets/images")
	if err != nil {
		t.Fatalf("batch load failed: %v", err)
	}
	file := "jeremy.png"
	_, err = GetSprite(file)
	if err != nil {
		t.Fatalf("get jeremy should have succeeded: %v", err)
	}
	DefaultCache.Clear(file)
	_, err = GetSprite(file)
	if err == nil {
		t.Fatal("get jeremy should have failed post-Clear")
	}
}
