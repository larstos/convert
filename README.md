# Convert

----
Just a tool package support transfer interface to specific type in go.  
* support transfer interface{} to int64\string\bool\float64
* support to figure out if interface{} likely to use as int64\string\bool\float64
* support fill []interface and map[string]interface{} to struct follow optional tag in struct field
* support fill struct with simple type suit. For example, func will try to fill string value to int64 field.

# How To Use
````
package main

import (
	"github.com/larstos/convert"
)

type Student struct {
	Name  string `json:"name"`
	Class *Class `json:"class"`
}

type Class struct {
	Grade       string `json:"grade"`
	ClassNumber int    `json:"class_number"`
}

func main() {
    //data is the result unmarshalled by common serialization api
    data:=map[string]interface{}{
        "name": "xiaoming",
        "class": map[string]interface{}{
            "grade":"first",
            "class_number":2,
        },
    }
    ret, err := convert_pkg.GetDataStructFilled(reflect.TypeOf(Student{}),data )
    if err != nil {
    	// error handler
    }
    val:=ret.(Student)
    
    ...
}
````