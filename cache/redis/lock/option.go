package lock

import "time"

// Options is lock options
type Options struct {
	// Expiration is expiration of the lock, default is 5s
	Expiration time.Duration
}

// OptionFunc is function to set Options
type OptionFunc func(*Options)

// DefaultOptions is default redis lock options
var DefaultOptions = Options{
	Expiration: 5 * time.Second,
}

// Expire sets lock's expiration
func Expire(d time.Duration) OptionFunc {
	return func(o *Options) {
		o.Expiration = d
	}
}
