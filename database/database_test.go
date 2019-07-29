package database

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDatabase(t *testing.T) {
	Convey("初始化DB", t, func() {
		_, err := Database("company")
		Convey("初始化DB", func() {
			So(err, ShouldEqual, nil)
		})
	})
}
