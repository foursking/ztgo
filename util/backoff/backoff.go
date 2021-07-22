package backoff

import (
	"math"
	"math/rand"
	"time"
)

const (
	defaultFactor = 2
	defaultMin    = 100 * time.Millisecond
	defaultMax    = 10 * time.Second
	maxInt64      = float64(math.MaxInt64 - 512)
)

// Backoff 是一个 time.Duration 计数器
type Backoff struct {
	// attempt 是尝试次数
	attempt int
	// Factor 是指数的底数，Factor的值越大，time.Duration 增幅越大
	Factor float64
	// Jitter 是增加随机数的开关
	Jitter bool
	// Min 最小等待时间，供首次使用
	Min time.Duration
	// Max 允许的最大等待时间
	Max time.Duration
}

// Duration 返回当次尝试需要等待的时间，attempt从0自增
// 如果需要手动指定 attempt，使用 ForAttempt
func (b *Backoff) Duration() time.Duration {
	d := b.ForAttempt(b.attempt)
	b.attempt++
	return d
}

// ForAttempt 根据传入的 attempt，计算并返回当次尝试需要等待的时间
// 这个方法是线程安全的
func (b *Backoff) ForAttempt(attempt int) time.Duration {
	min := b.Min
	if min <= 0 {
		min = defaultMin
	}
	max := b.Max
	if max <= 0 {
		max = defaultMax
	}
	if min >= max {
		return max
	}
	factor := b.Factor
	if factor <= 0 {
		factor = defaultFactor
	}
	minf := float64(min)
	durf := minf * math.Pow(factor, float64(attempt))
	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}
	if durf > maxInt64 {
		return max
	}
	dur := time.Duration(durf)
	if dur < min {
		return min
	}
	if dur > max {
		return max
	}
	return dur
}

// Reset 重置尝试次数
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt 返回当前尝试次数
func (b *Backoff) Attempt() int {
	return b.attempt
}
