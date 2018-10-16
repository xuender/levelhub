# README

levelhub is a wrapper of leveldb which supports multi-tenancy

## Installation

```shell
go get -u github.com/xuender/levelhub
```

## Usage

### Default

```go
package main

import "github.com/xuender/levelhub"

func main() {
	hub := levelhub.NewLevelHub("dbpath", nil)
	defer hub.Close()
	hub.Put(1, []byte("key"), []byte("val"), nil)
	hub.Put(2, []byte("key"), []byte("val"), nil)
}
```

### Expire

```go
hub := levelhub.NewLevelHub("dbpath", &levelhub.Options{
    Expire: time.Second * 3,
    Min:    1, // More than Min settings will trigger the expiration
})
defer hub.Close()
hub.Put(1, []byte("key"), []byte("val"), nil)
hub.Put(2, []byte("key"), []byte("val"), nil)
fmt.Println(hub.IsOpen(1), hub.IsOpen(2))
time.Sleep(time.Second * 10)
fmt.Println(hub.IsOpen(1), hub.IsOpen(2))
```