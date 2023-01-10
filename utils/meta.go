package utils

import (
	"reflect"
	"strconv"
)

func Inspect(f interface{}) map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(f)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		f := valueField.Interface()
		val := reflect.ValueOf(f)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[typeField.Name] = strconv.FormatInt(val.Int(), 10)
		case reflect.String:
			m[typeField.Name] = val.String()
		case reflect.Bool:
			m[typeField.Name] = strconv.FormatBool(val.Bool())
		case reflect.Float32, reflect.Float64:
			m[typeField.Name] = strconv.FormatFloat(val.Float(), 'E', -1, 64)
		}
	}
	return m
}
