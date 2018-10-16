package levelhub

import (
	"time"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Options hub options
type Options struct {
	DBOptions *opt.Options
	Min       int
	Max       int
	Expire    time.Duration
}

// NewOptions create Options
func NewOptions(o *opt.Options) *Options {
	return &Options{
		DBOptions: o,
		Min:       defaultMin,
		Max:       defaultMax,
		Expire:    defaultExpire,
	}
}
