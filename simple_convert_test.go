package convert_pkg

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMustString(t *testing.T) {
	Convey("Given some case for string transfer", t, func() {
		type unNameStruct struct {
			name string
		}
		defaultval := "default"
		Convey("support transfer type", func() {
			So(MustString(1, defaultval), ShouldEqual, "1")
			So(MustString(true, defaultval), ShouldEqual, "true")
			So(MustString(float64(1), defaultval), ShouldEqual, "1")
			So(MustString("1", defaultval), ShouldEqual, "1")
		})
		Convey("unsupport transfer type", func() {
			So(MustString(nil, defaultval), ShouldEqual, defaultval)
			So(MustString(unNameStruct{}, defaultval), ShouldEqual, defaultval)
			So(MustString(&unNameStruct{}, defaultval), ShouldEqual, defaultval)
			So(MustString([]interface{}{}, defaultval), ShouldEqual, defaultval)
			So(MustString(map[string]interface{}{}, defaultval), ShouldEqual, defaultval)
		})
	})
}
