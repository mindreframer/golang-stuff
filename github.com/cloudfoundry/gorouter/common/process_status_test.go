package common

import (
	"testing"
)

func BenchmarkUpdate(b *testing.B) {
	stat := NewProcessStatus()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stat.Update()
	}
}
