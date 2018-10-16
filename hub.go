package levelhub

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	defaultExpire = time.Minute * 5
	defaultMin    = 5
	defaultMax    = 10
	_Get          = iota
	_Clean
	_Force
)

// LevelHub leveldb hub
type LevelHub struct {
	cache     map[int]*data
	path      string
	o         *Options
	callBackC chan *callBack
}

type data struct {
	db     *leveldb.DB
	access time.Time
}

type callBack struct {
	route     int
	num       int
	db        *leveldb.DB
	err       error
	callBackC chan *callBack
}

func (hub *LevelHub) run() {
	for {
		cb := <-hub.callBackC
		switch cb.route {
		case _Get:
			if d, ok := hub.cache[cb.num]; ok {
				d.access = time.Now()
				cb.callBackC <- &callBack{
					db: d.db,
				}
			} else {
				db, err := leveldb.OpenFile(filepath.Join(hub.path, fmt.Sprintf("%x", cb.num/256), fmt.Sprintf("%x", cb.num%256)), hub.o.DBOptions)
				if err != nil {
					cb.callBackC <- &callBack{
						err: err,
					}
				} else {
					hub.cache[cb.num] = &data{
						access: time.Now(),
						db:     db,
					}
					cb.callBackC <- &callBack{
						db: db,
					}
				}
			}
		case _Clean:
			now := time.Now()
			dels := []int{}
			for num, d := range hub.cache {
				if now.Sub(d.access) > hub.o.Expire {
					dels = append(dels, num)
				}
			}
			for _, num := range dels {
				if d, ok := hub.cache[num]; ok {
					d.db.Close()
				}
				delete(hub.cache, num)
			}
			cb.callBackC <- &callBack{}
		case _Force:
			log.Printf("LevelHub cache size is %d > %d ! \n", len(hub.cache), hub.o.Max)
			for len(hub.cache) > hub.o.Max {
				t := time.Now()
				del := 0
				for num, d := range hub.cache {
					if d.access.Before(t) {
						t = d.access
						del = num
					}
				}
				if del > 0 {
					if d, ok := hub.cache[del]; ok {
						d.db.Close()
					}
					delete(hub.cache, del)
				}
			}
			cb.callBackC <- &callBack{}
		}
	}
}

func (hub *LevelHub) db(num int) (*leveldb.DB, error) {
	cs := hub.send(&callBack{
		route: _Get,
		num:   num,
	})
	return cs.db, cs.err
}

func (hub *LevelHub) send(cb *callBack) *callBack {
	cbC := make(chan *callBack, 1)
	defer close(cbC)
	cb.callBackC = cbC
	hub.callBackC <- cb
	return <-cbC
}

// NewLevelHub create LevelHub
func NewLevelHub(path string, o *Options) *LevelHub {
	if o == nil {
		o = NewOptions(nil)
	}
	if o.Expire < defaultExpire {
		o.Expire = defaultExpire
	}
	if o.Min < defaultMin {
		o.Min = defaultMin
	}
	if o.Max < defaultMax {
		o.Max = defaultMax
	}
	hub := &LevelHub{
		path:      path,
		o:         o,
		cache:     map[int]*data{},
		callBackC: make(chan *callBack, 3),
	}

	go hub.run()
	t := time.Minute
	if o.Expire < t {
		t = o.Expire
	}
	ticker := time.NewTicker(t)
	go func() {
		for range ticker.C {
			size := len(hub.cache)
			if size > hub.o.Min {
				now := time.Now()
				for _, d := range hub.cache {
					if now.Sub(d.access) > hub.o.Expire {
						hub.send(&callBack{
							route: _Clean,
						})
						break
					}
				}
				if size > hub.o.Max {
					hub.send(&callBack{
						route: _Force,
					})
				}
			}
		}
	}()
	return hub
}

// Close closes the nums DB.
func (hub *LevelHub) Close(nums ...int) {
	if len(nums) == 0 {
		for num, d := range hub.cache {
			d.db.Close()
			delete(hub.cache, num)
		}
	} else {
		for _, num := range nums {
			if d, ok := hub.cache[num]; ok {
				d.db.Close()
				delete(hub.cache, num)
			}
		}
	}
}

// Get gets the value for the given key.
func (hub *LevelHub) Get(num int, key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	db, err := hub.db(num)
	if err != nil {
		return
	}
	return db.Get(key, ro)
}

// Has returns true if the DB does contains the given key.
func (hub *LevelHub) Has(num int, key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	db, err := hub.db(num)
	if err != nil {
		return
	}
	return db.Has(key, ro)
}

// NewIterator returns an iterator for the latest snapshot of the underlying DB.
func (hub *LevelHub) NewIterator(num int, slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	db, err := hub.db(num)
	if err != nil {
		return iterator.NewEmptyIterator(err)
	}
	return db.NewIterator(slice, ro)
}

// Put sets the value for the given key.
func (hub *LevelHub) Put(num int, key, value []byte, wo *opt.WriteOptions) error {
	db, err := hub.db(num)
	if err != nil {
		return err
	}
	return db.Put(key, value, wo)
}

// Delete deletes the value for the given key.
func (hub *LevelHub) Delete(num int, key []byte, wo *opt.WriteOptions) error {
	db, err := hub.db(num)
	if err != nil {
		return err
	}
	return db.Delete(key, wo)
}
