package addressing

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddress(t *testing.T) {
	convey.Convey("test address", t, func() {
		convey.Convey("l5 address", func() {
			addr, err := Address("l5://1149313:131072")
			convey.So(err, convey.ShouldBeNil)
			addr.ReportL5(&err)
			t.Logf("l5 address result(%s)", addr)
		})
	})
}
