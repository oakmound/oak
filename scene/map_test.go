package scene

import (
	"errors"
	"image"
	"testing"

	"github.com/oakmound/oak/v3/oakerr"
)

func TestMap(t *testing.T) {
	m := NewMap()
	_, ok := m.Get("badScene")
	if ok {
		t.Fatalf("get on an undefined scene should fail")
	}
	if err := m.AddScene("test", Scene{}); err != nil {
		t.Fatalf("scene add should succeed, got %v", err)
	}

	err := m.AddScene("test", Scene{})
	if err == nil {
		t.Fatalf("duplicate scene add should fail")
	}
	exists := &oakerr.ExistingElement{}
	if !errors.As(err, exists) {
		t.Fatalf("expected ExistingElement error type, got %T", err)
	}
	if exists.InputName != "test" {
		t.Fatalf("expected error input 'test', got %q", exists.InputName)
	}
	if exists.InputType != "scene" {
		t.Fatalf("expected error input type 'scene', got %q", exists.InputType)
	}

	_, ok = m.Get("test")
	if !ok {
		t.Fatalf("getting test scene failed")
	}
	m.CurrentScene = "test"
	_, ok = m.GetCurrent()
	if !ok {
		t.Fatalf("getting current test scene failed")
	}
}

func TestFade(t *testing.T) {
	fadeFn := Fade(1, 10)
	if fadeFn(nil, 11) {
		t.Fatalf("fade should not proceed after its frames have elapsed")
	}
	if !fadeFn(image.NewRGBA(image.Rect(0, 0, 50, 50)), 2) {
		t.Fatalf("fade should proceed before its frames have elapsed")
	}
}

func TestZoom(t *testing.T) {
	zoomFn := Zoom(.5, .5, 10, .1)
	if zoomFn(nil, 11) {
		t.Fatalf("zoom should not proceed after its frames have elapsed")
	}
	if !zoomFn(image.NewRGBA(image.Rect(0, 0, 50, 50)), 2) {
		t.Fatalf("zoom should proceed before its frames have elapsed")
	}
}

func TestAddScene(t *testing.T) {
	m := NewMap()
	_, ok := m.Get("badScene")
	if ok {
		t.Fatalf("unexpected success getting undefined scene")
	}

	m.AddScene("test1", Scene{})
	test1, ok := m.Get("test1")
	if !ok {
		t.Fatalf("getting test scene failed")
	}

	if !test1.Loop() {
		t.Fatalf("test loop failed")
	}
	eStr, _ := test1.End()
	if eStr != "test1" {
		t.Fatalf("looping test end did not return test1, got %v", eStr)
	}
	test1.Start(&Context{PreviousScene: "test"})
}
