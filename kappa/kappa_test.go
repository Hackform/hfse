package kappa

import (
	"testing"
)

func TestKappa(t *testing.T) {
	k := New()
	k1 := k.Get()
	k2 := k.Get()
	k3 := k.Get()

	if k1 != 1 {
		t.Error("Expected 1, got ", k1)
	}
	if k2 != 2 {
		t.Error("Expected 2, got ", k2)
	}
	if k3 != 3 {
		t.Error("Expected 3, got ", k3)
	}
}
