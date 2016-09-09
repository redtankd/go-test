package mytest

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
    // t.Fail()
    // t.Logf("%d", 3)
	t.Skip()
    
}

func BenchmarkAdd(b *testing.B) {
	// If a benchmark needs some expensive setup before running
	b.ResetTimer()
	// b.StartTimer()
	// b.StopTimer()

	for i := 0; i < b.N; i++ {
		Add(5, 6)
	}
}

func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Add(5, 6)
		}
	})
}

func ExampleAdd() {
	fmt.Println(Add(1, 2))
	// Output: 3
}
