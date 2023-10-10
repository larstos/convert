package convert

import (
	"errors"
	"fmt"
	"reflect"
)

const default_tag = "json"

// GetDataStructFilled will use reflect to transfer interface to target type.
func GetDataStructFilled(datatype reflect.Type, value any, tagName ...string) (any, error) {
	ret, err := fillDataStructValue(datatype, value, tagName...)
	if err != nil {
		return nil, err
	}
	return ret.Interface(), err
}

// GetDataStructFilledWithMap fill value depends on map record.
// Func will ignore fields :
// 1. tag is "-" or tag is empty
// 2. field not accessible
// 3. field not been set in map
func GetDataStructFilledWithMap(datatype reflect.Type, value map[string]any, careTags ...string) (any, error) {
	tag := default_tag
	if len(careTags) > 0 {
		tag = careTags[0]
	}
	operateDataType := datatype
	if datatype.Kind() == reflect.Ptr {
		operateDataType = datatype.Elem()
	}
	datanew := reflect.New(operateDataType)
	datanewval := datanew.Elem()
	for i := 0; i < operateDataType.NumField(); i++ {
		var ok bool
		tfield := operateDataType.Field(i)
		if len(tfield.PkgPath) > 0 {
			continue
		}
		tagFieldName := tfield.Tag.Get(tag)
		if tagFieldName == "" || tagFieldName == "-" {
			continue
		}
		rawVal, ok := value[tagFieldName]
		if !ok {
			continue
		}
		vfield := datanewval.Field(i)
		if !vfield.CanSet() {
			continue
		}
		value, err := fillDataStructValue(tfield.Type, rawVal, tag)
		if err != nil {
			return nil, fmt.Errorf("[error] error parse %s,err:%v", tagFieldName, err)
		}
		vfield.Set(value)
	}
	if datatype.Kind() != reflect.Ptr {
		return datanew.Elem().Interface(), nil
	}
	return datanew.Interface(), nil
}

func fillDataStructValue(datatype reflect.Type, value any, tagName ...string) (val reflect.Value, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("[painc]%v", p)
		}
	}()
	ret := reflect.New(datatype).Elem()
	tag := default_tag
	if len(tagName) > 0 {
		tag = tagName[0]
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
		val, ok := value.(map[string]any)
		if !ok {
			return reflect.Zero(datatype), errors.New("data type not valid")
		}
		inner, err := GetDataStructFilledWithMap(datatype, val, tag)
		if err != nil {
			return reflect.Zero(datatype), err
		}
		return reflect.ValueOf(inner), nil
	case reflect.Array, reflect.Slice:
		list, ok := value.([]any)
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
		innermap, ok := value.(map[string]any)
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
	return ret, nil
}
