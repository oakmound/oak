package render

import (
	"testing"
)

// We are rather intentionally not testing that the decoders work,
// partially because it's a lot of work and partially because decoding images is not
// this package's job (it just calls decoders). It'd be a lot of pointless mocking.

func TestRegisterDecoder(t *testing.T) {
	err := RegisterDecoder(".png", nil)
	if err == nil {
		t.Fatal("expected registering .png to fail")
	}
	err = RegisterDecoder(".new", nil)
	if err != nil {
		t.Fatalf("expected registering .new to succeed: %v", err)
	}
}

func TestRegisterCfgDecoder(t *testing.T) {
	err := RegisterCfgDecoder(".png", nil)
	if err == nil {
		t.Fatal("expected registering .png to fail")
	}
	err = RegisterCfgDecoder(".new", nil)
	if err != nil {
		t.Fatalf("expected registering .new to succeed: %v", err)
	}
}
