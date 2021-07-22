package rand

import (
	"math/rand"
	"testing"
	"time"
)

// go test -bench=. -test.cpu=1,4,8

const maxInt = 1000

var r *Rand

// 初始化
func init() {
	r = New()
	rand.Seed(time.Now().UnixNano())
}

// 测试可扩展的随机数
func BenchmarkRand_Intn(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Intn(maxInt)
		}
	})
}

// 测试默认的随机数
func BenchmarkBuildInRand_Intn(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Intn(maxInt)
		}
	})
}
