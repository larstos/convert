package convert

import (
	"fmt"
	"reflect"
)

func GetTagValues(object any, tagName string) ([]string, error) {
	if tagName == "" {
		return []string{}, nil
	}
	rType := reflect.TypeOf(object)
	if rType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("object must be a pointer")
	}
	rType = rType.Elem()
	if rType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("object must be a struct")
	}
	fNames := make([]string, 0, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		name := rType.Field(i).Tag.Get(tagName)
		if name == "" || name == "-" {
			continue
		}
		fNames = append(fNames, name)
	}
	return fNames, nil
}

func GetDataByTagValue(object any, tagName, tagValue string) (any, bool) {
	if tagName == "" {
		return nil, false
	}
	rType := reflect.TypeOf(object)
	if rType.Kind() != reflect.Ptr {
		return nil, false
	}
	rType = rType.Elem()
	if rType.Kind() != reflect.Struct {
		return nil, false
	}
	for i := 0; i < rType.NumField(); i++ {
		name := rType.Field(i).Tag.Get(tagName)
		if name != "" && name == tagName {
			val := reflect.ValueOf(object).Elem().Field(i)
			if !val.CanInterface() {
				return nil, false
			}
			return val.Interface(), true
		}
	}
	return nil, false
}
