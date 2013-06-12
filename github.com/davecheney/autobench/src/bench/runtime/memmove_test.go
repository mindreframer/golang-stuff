package runtime_test

import "testing"

func bmMemmove(b *testing.B, size int) {
	src, dst := make([]byte, size), make([]byte, size)
	b.SetBytes(int64(size))
	b.ResetTimer()
	for i := 0 ; i < b.N ; i++ {
		copy(dst, src)
	}
}

func BenchmarkMemmove32(b *testing.B) { bmMemmove(b, 32) }
func BenchmarkMemmove4K(b *testing.B) { bmMemmove(b, 4<<10) }
func BenchmarkMemmove64K(b *testing.B) { bmMemmove(b, 64<<10) }
func BenchmarkMemmove4M(b *testing.B) { bmMemmove(b, 4<<20) }
func BenchmarkMemmove64M(b *testing.B) { bmMemmove(b, 64<<20) }
