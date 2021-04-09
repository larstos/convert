package convert_pkg

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDataStructFilled(t *testing.T) {
	Convey("test get data struct field fill type will not change", t, func() {
		Convey("test struct and pointer", func() {
			type unNameStruct struct {
				Name string `json:"name"`
			}
			//struct
			ret, err := GetDataStructFilled(reflect.TypeOf(unNameStruct{}), map[string]interface{}{"name": 1})
			So(err, ShouldBeNil)
			So(ret, ShouldResemble, unNameStruct{
				Name: "1",
			})
			//pointer
			ret, err = GetDataStructFilled(reflect.TypeOf(&unNameStruct{}), map[string]interface{}{"name": 1})
			So(err, ShouldBeNil)
			So(ret, ShouldResemble, &unNameStruct{
				Name: "1",
			})
		})

	})
}
