package convert_pkg

import (
	"errors"
	"fmt"
	"reflect"
)

const default_tag = "json"

//fillDataStructValue will use reflect to transfer interface to target type.
func GetDataStructFilled(datatype reflect.Type, value interface{}, careTags ...string) (interface{}, error) {
	ret, err := fillDataStructValue(datatype, value, careTags...)
	if err != nil {
		return nil, err
	}
	return ret.Interface(), err
}

//GetDataStructFilledWithMap fill value depends on map record.
//Func will ignore fields :
// 1. tag is "-" or tag is empty
// 2. field not accessible
// 3. field not been set in map
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
		tagFieldName := tfield.Tag.Get(tag)
		if len(tagFieldName) == 0 || tagFieldName == "-" {
			continue
		} else {
			_, ok = value[tagFieldName]
			if !ok {
				continue
			}
		}
		vfield := datanewval.Field(i)
		if vfield.CanSet() {
			value, err := fillDataStructValue(tfield.Type, value[tagFieldName], tag)
			if err != nil {
				return nil, fmt.Errorf("[error] error parse %s,err:%v", tagFieldName, err)
			}
			vfield.Set(value)
		} else {
			return nil, fmt.Errorf("data type=%v, %v can not be set", datatype, tagFieldName)
		}
	}
	if datatype.Kind() != reflect.Ptr {
		return datanew.Elem().Interface(), nil
	}
	return datanew.Interface(), nil
}

func fillDataStructValue(datatype reflect.Type, value interface{}, careTags ...string) (reflect.Value, error) {
	var panicerr error
	defer func() {
		if err := recover(); err != nil {
			panicerr = fmt.Errorf("%v", err)
		}
	}()
	ret := reflect.New(datatype).Elem()
	tag := default_tag
	if len(careTags) > 0 {
		tag = careTags[0]
	}
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
		ret = reflect.MakeSlice(datatype, 0, len(list))
		for _, i2 := range list {
			interval, err := fillDataStructValue(elem, i2, tag)
			if err != nil {
				return reflect.Zero(datatype), errors.New("data type not valid")
			}
			ret = reflect.Append(ret, interval)
		}
	case reflect.Map:
		innermap, ok := value.(map[string]interface{})
		if !ok {
			return reflect.Zero(datatype), errors.New("data type not valid")
		}
		keytype := datatype.Key()
		valtype := datatype.Elem()
		for k, i2 := range innermap {
			if keytype.Kind() != reflect.String {
				return reflect.Zero(datatype), errors.New("data type not valid,err: map key should be string type")
			}
			innerkey := reflect.ValueOf(k)
			//fill value
			innervalue, err := fillDataStructValue(valtype, i2, tag)
			if err != nil {
				return reflect.Zero(datatype), fmt.Errorf("data type not valid,value:%v,err:%v", value, err)
			}
			ret.SetMapIndex(innerkey, innervalue)
		}
	case reflect.Interface:
		ret = reflect.ValueOf(value)
	default:
		return ret, errors.New("unsupported type")
	}
	return ret, panicerr
}
