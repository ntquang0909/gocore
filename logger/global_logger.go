package logger

import (
	"sync"
)

var once = sync.Once{}

var globalLog *Logger

// Global global
func Global() *Logger {
	return setupLog()
}

func setupLog() *Logger {
	if globalLog == nil {
		once.Do(func() {
			globalLog = New(nil)
		})
	}
	return globalLog
}
