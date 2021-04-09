# ConvertPkg

----
Just a tool package support transfer interface to specific type in go.  
* support transfer interface{} to int64\string\bool\float64
* support to figure out if interface{} likely to use as int64\string\bool\float64
* support fill []interface and map[string]interface{} to struct follow optional tag in struct field
* support fill struct with simple type suit. For example, func will try to fill string value to int64 field.
