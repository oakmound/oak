package scene

import (
	"math/rand"
	"testing"
)

func randString() string {
	length := rand.Intn(100)
	data := make([]byte, length)
	for i := range data {
		data[i] = byte(rand.Intn(255))
	}
	return string(data)
}

func TestGoTo(t *testing.T) {
	tests := 10
	for i := 0; i < tests; i++ {
		s := randString()
		gt := GoTo(s)
		s2, result := gt()
		if s != s2 {
			t.Fatalf("expected goto to return %v, got %v", s, s2)
		}
		if result != nil {
			t.Fatalf("expected goto to return nil result, got %v", result)
		}
	}
}

func TestGoToPtr(t *testing.T) {
	tests := 10
	s := new(string)
	gt := GoToPtr(s)
	for i := 0; i < tests; i++ {
		*s = randString()
		s2, result := gt()
		if *s != s2 {
			t.Fatalf("expected gotoptr to return %v, got %v", *s, s2)
		}
		if result != nil {
			t.Fatalf("expected gotoptr to return nil result, got %v", result)
		}
	}
}

func TestGoToPtrNil(t *testing.T) {
	s, _ := GoToPtr(nil)()
	if s != "" {
		t.Fatalf("expected nil gotoptr to return empty string, got %v", s)
	}
}
