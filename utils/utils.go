package utils

import "reflect"

// IsPtr check pointer
func IsPtr(val interface{}) bool {
	v := reflect.ValueOf(val)
	return v.Kind() == reflect.Ptr
}

// IsArr check array
func IsArr(val interface{}) bool {
	v := reflect.ValueOf(val)
	return v.Kind() == reflect.Array
}
