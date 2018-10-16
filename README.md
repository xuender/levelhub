# README
levelhub is a wrapper of leveldb which supports multi-tenancy

## Installation
```shell
go get -u github.com/xuender/levelhub
```
## Usage:
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