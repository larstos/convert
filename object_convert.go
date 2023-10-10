package convert

import (
	"errors"
	"fmt"
	"reflect"
)

const default_tag = "json"

// NewDataStructFilled will use reflect to transfer interface to target type.
func NewDataStructFilled(datatype reflect.Type, value any, tagName ...string) (any, error) {
	datanew := reflect.New(datatype).Elem()
	err := fillDataStructValue(datanew, value, tagName...)
	if err != nil {
		return nil, err
	}
	return datanew.Interface(), nil
}

func GetDataStructFilled(input, value any, tagName ...string) error {
	datanew := reflect.ValueOf(input)
	if datanew.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("%v is not a pointer", input)
	}
	err := fillDataStructValue(datanew.Elem(), value, tagName...)
	if err != nil {
		return err
	}
	return nil
}

// NewDataStructFilledWithMap fill value depends on map record.
// Func will ignore fields :
// 1. tag is "-" or tag is empty
// 2. field not accessible
// 3. field not been set in map
func NewDataStructFilledWithMap(datatype reflect.Type, value map[string]any, careTags ...string) (any, error) {
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
		err := fillDataStructValue(vfield, rawVal, tag)
		if err != nil {
			return nil, fmt.Errorf("[error] error parse %s,err:%v", tagFieldName, err)
		}
	}
	if datatype.Kind() != reflect.Ptr {
		return datanew.Elem().Interface(), nil
	}
	return datanew.Interface(), nil
}

func fillDataStructValue(dataVal reflect.Value, raw any, tagName ...string) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("[painc]%v", p)
		}
	}()
	dataType := dataVal.Type()
	tag := default_tag
	if len(tagName) > 0 {
		tag = tagName[0]
	}
	switch dataType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dataVal.SetInt(MustInt64(raw))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dataVal.SetUint(uint64(MustInt64(raw)))
	case reflect.String:
		dataVal.SetString(MustString(raw))
	case reflect.Bool:
		dataVal.SetBool(MustBool(raw))
	case reflect.Float64, reflect.Float32:
		dataVal.SetFloat(MustFloat64(raw))
	case reflect.Ptr, reflect.Struct:
		val, ok := raw.(map[string]any)
		if !ok {
			return errors.New("data type not valid")
		}
		inner, err := NewDataStructFilledWithMap(dataType, val, tag)
		if err != nil {
			return err
		}
		dataVal.Set(reflect.ValueOf(inner))
	case reflect.Array, reflect.Slice:
		list, ok := raw.([]any)
		if !ok {
			return errors.New("data type not valid")
		}
		tList := reflect.MakeSlice(dataType, 0, len(list))
		for _, i2 := range list {
			tmpVal := reflect.New(dataType.Elem())
			err := fillDataStructValue(tmpVal, i2, tag)
			if err != nil {
				return err
			}
			tList = reflect.Append(tList, tmpVal)
		}
		dataVal.Set(tList)
	case reflect.Map:
		innermap, ok := raw.(map[string]any)
		if !ok {
			return errors.New("data type not valid")
		}
		tmap := reflect.MakeMapWithSize(dataType, len(innermap))
		keytype := dataType.Key()
		for k, i2 := range innermap {
			if keytype.Kind() != reflect.String {
				return errors.New("data type not valid,err: map key should be string type")
			}
			tmpVal := reflect.New(dataType.Elem())
			err := fillDataStructValue(tmpVal, i2, tag)
			if err != nil {
				return fmt.Errorf("data type not valid,value:%v,err:%v", raw, err)
			}
			tmap.SetMapIndex(reflect.ValueOf(k), tmpVal)
		}
		dataVal.Set(tmap)
	case reflect.Interface:
		dataVal.Set(reflect.ValueOf(raw))
	default:
		return errors.New("unsupported type")
	}
	return nil
}
