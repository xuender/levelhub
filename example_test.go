package levelhub

import (
	"fmt"
	"time"
)

func ExampleNewLevelHub() {
	hub := NewLevelHub("tmp", nil)
	defer hub.Close()
	hub.Put(1, []byte("key"), []byte("A"), nil)
	hub.Put(2, []byte("key"), []byte("B"), nil)
	a, _ := hub.Get(1, []byte("key"), nil)
	fmt.Println(string(a))
	b, _ := hub.Get(2, []byte("key"), nil)
	fmt.Println(string(b))

	// Output:
	// A
	// B
}

func ExampleOptions() {
	hub := NewLevelHub("tmp", &Options{
		Expire: time.Second * 3,
		Min:    1, // More than Min settings will trigger the expiration
	})
	defer hub.Close()
	hub.Put(1, []byte("key"), []byte("val"), nil)
	hub.Put(2, []byte("key"), []byte("val"), nil)
	fmt.Println(hub.IsOpen(1), hub.IsOpen(2))
	time.Sleep(time.Second * 10)
	fmt.Println(hub.IsOpen(1), hub.IsOpen(2))

	// Output:
	// true true
	// false false
}
