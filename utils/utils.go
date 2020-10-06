package utils

import (
	"bytes"
	"reflect"
	"text/template"
)

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

// IsUnix check unix
func IsUnix(value *int64, nullable ...bool) bool {
	var canNull = true
	if len(nullable) > 0 {
		canNull = nullable[0]
	}

	if value == nil {
		return canNull
	}

	return *value > 0
}

// IsLat check unix
func IsLat(value *int64, nullable ...bool) bool {
	var canNull = true
	if len(nullable) > 0 {
		canNull = nullable[0]
	}

	if value == nil {
		return canNull
	}

	var lat = *value
	return lat >= -90 && lat <= 90
}

// IsLng check unix
func IsLng(value *int64, nullable ...bool) bool {
	var canNull = true
	if len(nullable) > 0 {
		canNull = nullable[0]
	}

	if value == nil {
		return canNull
	}

	var lng = *value
	return lng >= -180 && lng <= 180
}

// ParseHTML html
func ParseHTML(htmlFile string, i interface{}) (string, error) {
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
