package levelhub

import (
	"testing"
)

func BenchmarkPut(b *testing.B) {
	b.StopTimer()
	hub := NewLevelHub("tmp", nil)
	defer hub.Close()
	key := []byte("key")
	val := []byte("val")
	hub.Put(1, key, val, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hub.Put(1, key, val, nil)
	}
}

func BenchmarkClean(b *testing.B) {
	b.StopTimer()
	hub := NewLevelHub("tmp", nil)
	defer hub.Close()
	key := []byte("key")
	val := []byte("val")
	hub.Put(1, key, val, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hub.Clean()
	}
}
