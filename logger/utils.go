package logger

import (
	"encoding/json"
	"fmt"
)

// Writer interface
type Writer interface {
	Printf(string, ...interface{})
	Print(...interface{})
}

// LogLevel log level
type LogLevel int

// All levels
const (
	Debug LogLevel = iota + 1
	Warn
	Info
	Error
)

// Colors
var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Blue    = Color("\033[1;34m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

// Color format text color
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

// ToJSONString convert to json string
func ToJSONString(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Couldn't marshal error %v", err).Error()
	}

	return string(data)
}

// ToPrettyJSONString convert to json string
func ToPrettyJSONString(value interface{}) string {
	data, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return fmt.Errorf("Couldn't marshal error %v", err).Error()
	}

	return string(data)
}
