package pg

import "testing"

func BenchmarkErrorCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Error("Hello")
	}
}

func BenchmarkError(b *testing.B) {
	err := Error("Hello")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
