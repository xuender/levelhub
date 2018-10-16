# README
levelhub is a wrapper of leveldb which supports multi-tenancy

## Installation
```shell
go get -u github.com/xuender/levelhub
```
## Usage:
```go
import "github.com/xuender/levelhub"
...
hub := levelhub.NewLevelHub("path", nil)
hub.Put(1, []byte("key"), []byte("val"))
hub.Put(2, []byte("key"), []byte("val"))
```