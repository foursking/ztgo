package backoff

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDuration(t *testing.T) {
	Convey("测试 Duration()", t, func() {
		Convey("Factor == 2", func() {
			b := &Backoff{
				Min:    100 * time.Millisecond,
				Max:    10 * time.Second,
				Factor: 2,
			}
			So(b.Duration(), ShouldEqual, 100*time.Millisecond)
			So(b.Duration(), ShouldEqual, 200*time.Millisecond)
			So(b.Duration(), ShouldEqual, 400*time.Millisecond)
			b.Reset()
			So(b.Duration(), ShouldEqual, 100*time.Millisecond)
		})

		Convey("Factor == 1.5", func() {
			b := &Backoff{
				Min:    100 * time.Millisecond,
				Max:    10 * time.Second,
				Factor: 1.5,
			}
			So(b.Duration(), ShouldEqual, 100*time.Millisecond)
			So(b.Duration(), ShouldEqual, 150*time.Millisecond)
			So(b.Duration(), ShouldEqual, 225*time.Millisecond)
			b.Reset()
			So(b.Duration(), ShouldEqual, 100*time.Millisecond)
		})

		Convey("Min > Max", func() {
			b := &Backoff{
				Min:    500 * time.Second,
				Max:    100 * time.Second,
				Factor: 1,
			}
			So(b.Duration(), ShouldEqual, b.Max)
		})
	})
}

func TestForAttempt(t *testing.T) {
	Convey("测试 ForAttempt()", t, func() {
		b := &Backoff{
			Min:    100 * time.Millisecond,
			Max:    10 * time.Second,
			Factor: 2,
		}
		So(b.ForAttempt(0), ShouldEqual, 100*time.Millisecond)
		So(b.ForAttempt(1), ShouldEqual, 200*time.Millisecond)
		So(b.ForAttempt(2), ShouldEqual, 400*time.Millisecond)
		b.Reset()
		So(b.ForAttempt(0), ShouldEqual, 100*time.Millisecond)
	})
}

func TestGetAttempt(t *testing.T) {
	Convey("测试 获取Attempt", t, func() {
		b := &Backoff{
			Min:    100 * time.Millisecond,
			Max:    10 * time.Second,
			Factor: 2,
		}
		So(b.Attempt(), ShouldEqual, float64(0))
		So(b.Duration(), ShouldEqual, 100*time.Millisecond)
		So(b.Attempt(), ShouldEqual, float64(1))
		So(b.Duration(), ShouldEqual, 200*time.Millisecond)
		So(b.Attempt(), ShouldEqual, float64(2))
		So(b.Duration(), ShouldEqual, 400*time.Millisecond)
		b.Reset()
		So(b.Attempt(), ShouldEqual, float64(0))
		So(b.Duration(), ShouldEqual, 100*time.Millisecond)
		So(b.Attempt(), ShouldEqual, float64(1))
	})
}

func TestJitter(t *testing.T) {
	Convey("测试 Jitter", t, func() {
		b := &Backoff{
			Min:    100 * time.Millisecond,
			Max:    10 * time.Second,
			Factor: 2,
			Jitter: true,
		}
		So(b.Duration(), ShouldEqual, 100*time.Millisecond)
		So(b.Duration(), ShouldBeBetween, 100*time.Millisecond, 200*time.Millisecond)
		So(b.Duration(), ShouldBeBetween, 100*time.Millisecond, 400*time.Millisecond)
		b.Reset()
		So(b.Duration(), ShouldEqual, 100*time.Millisecond)
	})
}
