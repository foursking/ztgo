// Package rand 是一个并发安全的、多核CPU高效的 rand 包
package rand

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Rand gets random number
type Rand struct {
	seed int64
	pool *sync.Pool
}

// New creates rand
func New() *Rand {
	r := new(Rand)
	r.pool = &sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(r.newSeed()))
		},
	}
	return r
}

func (r *Rand) newSeed() int64 {
	var seed int64
	for {
		seed = time.Now().UnixNano()
		cur := atomic.LoadInt64(&r.seed)
		if cur != seed {
			if atomic.CompareAndSwapInt64(&r.seed, cur, seed) {
				break
			}
		}
		time.Sleep(time.Nanosecond)
	}
	return seed
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0
func (r *Rand) Intn(n int) int {
	rd := r.pool.Get().(*rand.Rand)
	i := rd.Intn(n)
	r.pool.Put(rd)
	return i
}
