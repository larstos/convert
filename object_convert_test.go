package convert

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDataStructFilled(t *testing.T) {
	Convey("test base func works well", t, func() {
		type Class struct {
			Grade       string `json:"grade"`
			ClassNumber int    `json:"class_number"`
		}
		type Student struct {
			Name  string `json:"name"`
			Class *Class `json:"class"`
		}
		data := map[string]interface{}{
			"name": "xiaoming",
			"class": map[string]interface{}{
				"grade":        "first",
				"class_number": 2,
			},
		}
		ret, err := NewDataStructFilled(reflect.TypeOf(Student{}), data)
		So(err, ShouldBeNil)
		So(ret, ShouldResemble, Student{
			Name: "xiaoming",
			Class: &Class{
				Grade:       "first",
				ClassNumber: 2,
			},
		})
		// 已有对象直接填充
		info := &Student{}
		err = GetDataStructFilled(info, data)
		So(err, ShouldBeNil)
		So(info, ShouldResemble, &Student{
			Name: "xiaoming",
			Class: &Class{
				Grade:       "first",
				ClassNumber: 2,
			},
		})
		// 对特殊对象进行填充
		num := new(int64)
		err = GetDataStructFilled(num, 2)
		So(err, ShouldBeNil)
		So(*num, ShouldEqual, 2)
	})
	Convey("test get data struct field fill type will not change", t, func() {
		Convey("test struct and pointer", func() {
			type unNameStruct struct {
				Name string `json:"name"`
			}
			//struct
			ret, err := NewDataStructFilled(reflect.TypeOf(unNameStruct{}), map[string]interface{}{"name": 1})
			So(err, ShouldBeNil)
			So(ret, ShouldResemble, unNameStruct{
				Name: "1",
			})
			//pointer
			ret, err = NewDataStructFilled(reflect.TypeOf(&unNameStruct{}), map[string]interface{}{"name": 1})
			So(err, ShouldBeNil)
			So(ret, ShouldResemble, &unNameStruct{
				Name: "1",
			})
		})
	})
}
