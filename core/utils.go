package core

import (
	"errors"
	"reflect"
	"strings"
)

func StringInArray(key string, array []string) bool {
	if key == "" || len(array) == 0 {
		return false
	}
	for _, v := range array {
		if strings.ToLower(key) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func StructToMap(data interface{}, lowerKey bool, except ...string) (m map[string]interface{}, err error) {
	m = make(map[string]interface{}, 0)
	if data == nil {
		err = errors.New("struct to map error data is nil")
		return
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		err = errors.New("data must be struct")
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if StringInArray(t.Field(i).Name, except) {
			continue
		}
		if lowerKey {
			m[strings.ToLower(t.Field(i).Name)] = f.Interface()
		} else {
			m[t.Field(i).Name] = f.Interface()
		}
	}
	return
}
