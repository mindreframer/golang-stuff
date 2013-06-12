package skyd

import (
	"testing"
)

// Ensure it condenses a 64-bit number down to only it's even bits.
func TestCondenseUint64Even(t *testing.T) {
	if CondenseUint64Even(0x0) != 0x0 {
		t.Fatalf("CondenseUint64: expected %x, got %x", 0x0, CondenseUint64Even(0x0))
	}
	if CondenseUint64Even(0x58AB) != 0xC1 {
		t.Fatalf("CondenseUint64: expected %x, got %x", 0xC1, CondenseUint64Even(0x58AB))
	}
}

// Ensure it condenses a 64-bit number down to only it's odd bits.
func TestCondenseUint64Odd(t *testing.T) {
	if CondenseUint64Odd(0x0) != 0x0 {
		t.Fatalf("CondenseUint64: expected %x, got %x", 0x0, CondenseUint64Odd(0x0))
	}
	if CondenseUint64Odd(0x58AB) != 0x2F {
		t.Fatalf("CondenseUint64: expected %x, got %x", 0x2F, CondenseUint64Odd(0x58AB))
	}
}

func BenchmarkCondenseUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CondenseUint64Odd(0x58AB)
	}
}
