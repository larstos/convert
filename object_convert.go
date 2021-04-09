package convert_pkg

import (
	"errors"
	"fmt"
	"reflect"
)

const default_tag = "json"

func GetDataStructFilledWithMap(datatype reflect.Type, value map[string]interface{}, careTags ...string) (interface{}, error) {
	tag := default_tag
	if len(careTags) > 0 {
		tag = careTags[0]
	}
	var datanew reflect.Value
	operateDataType := datatype
	if datatype.Kind() == reflect.Ptr {
		operateDataType = datatype.Elem()
	}
	datanew = reflect.New(operateDataType)
	numField := operateDataType.NumField()
	datanewval := datanew.Elem()
	for i := 0; i < numField; i++ {
		var ok bool
		tfield := operateDataType.Field(i)
		if len(tfield.PkgPath) > 0 {
			return nil, fmt.Errorf("data type=%v,field %s is not public", datatype, tfield.Name)
		}
		//tag:"yaml",若忽略该项则`yaml:"-"`或不填
		yamlFieldName := tfield.Tag.Get(tag)
		if len(yamlFieldName) == 0 || yamlFieldName == "-" {
			continue
		} else {
			_, ok = value[yamlFieldName]
			//兼容select部分字段
			if !ok {
				continue
			}
		}
		vfield := datanewval.Field(i)
		if vfield.CanSet() {
			if ok {
				switch tfield.Type.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					vfield.SetInt(MustInt64(value[yamlFieldName]))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					vfield.SetUint(uint64(MustInt64(value[yamlFieldName])))
				case reflect.String:
					vfield.SetString(MustString(value[yamlFieldName]))
				case reflect.Bool:
					vfield.SetBool(MustBool(value[yamlFieldName]))
				case reflect.Float64, reflect.Float32:
					vfield.SetFloat(MustFloat64(value[yamlFieldName]))
				case reflect.Array, reflect.Slice:
					list, ok := value[yamlFieldName].([]interface{})
					if !ok {
						return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
					}
					val := reflect.Zero(tfield.Type)
					elem := tfield.Type.Elem()
					for _, i2 := range list {
						interval, err := GetDataStructFilled(elem, i2, tag)
						if err != nil {
							return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
						}
						val = reflect.Append(val, interval)
					}
					vfield.Set(val)
				case reflect.Ptr, reflect.Struct:
					val, ok := value[yamlFieldName].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
					}
					inner, err := GetDataStructFilled(tfield.Type, val, tag)
					if err != nil {
						return nil, err
					}
					vfield.Set(inner)
				case reflect.Map:
					innermap, ok := value[yamlFieldName].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
					}
					keytype := tfield.Type.Key()
					valtype := tfield.Type.Elem()
					for k, i2 := range innermap {
						innerkey, err := GetDataStructFilled(keytype, k, tag)
						if err != nil {
							return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
						}
						innervalue, err := GetDataStructFilled(valtype, i2, tag)
						if err != nil {
							return nil, fmt.Errorf("data type not valid,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
						}
						vfield.SetMapIndex(innerkey, innervalue)
					}
				default:
					return nil, fmt.Errorf("find unsupport result,data type=%v,field name=%v,type=%v,result=%s", datatype, tfield.Name, tfield.Type.String(), value[yamlFieldName])
				}
			}
		} else {
			return nil, fmt.Errorf("data type=%v, %v can not be set", datatype, yamlFieldName)
		}
	}
	if datatype.Kind() != reflect.Ptr {
		return datanew.Elem().Interface(), nil
	}
	return datanew.Interface(), nil
}

func GetDataStructFilled(datatype reflect.Type, value interface{}, careTags ...string) (reflect.Value, error) {
	tag := default_tag
	if len(careTags) > 0 {
		tag = careTags[0]
	}
	ret := reflect.New(datatype).Elem()
	switch datatype.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret.SetInt(MustInt64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret.SetUint(uint64(MustInt64(value)))
	case reflect.String:
		ret.SetString(MustString(value))
	case reflect.Bool:
		ret.SetBool(MustBool(value))
	case reflect.Float64, reflect.Float32:
		ret.SetFloat(MustFloat64(value))
	case reflect.Ptr, reflect.Struct:
		val, ok := value.(map[string]interface{})
		if !ok {
			return reflect.Zero(datatype), errors.New("data type not valid")
		}
		inner, err := GetDataStructFilledWithMap(datatype, val, tag)
		if err != nil {
			return reflect.Zero(datatype), err
		}
		return reflect.ValueOf(inner), nil
	case reflect.Array, reflect.Slice:
		list, ok := value.([]interface{})
		if !ok {
			return reflect.Zero(datatype), errors.New("data type not valid")
		}
		elem := datatype.Elem()
		for _, i2 := range list {
			interval, err := GetDataStructFilled(elem, i2, tag)
			if err != nil {
				return reflect.Zero(datatype), errors.New("data type not valid")
			}
			ret = reflect.Append(ret, interval)
		}
	default:
		return reflect.Zero(datatype), errors.New("unsupported type")
	}
	return ret, nil
}
