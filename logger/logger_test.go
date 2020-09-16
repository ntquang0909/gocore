package logger

import (
	"log"
	"testing"
	"time"

	"github.com/kjk/dailyrotate"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestRotateLogger(t *testing.T) {
	rotateLog, err := dailyrotate.NewFile("2006-01-02.log", func(path string, didRotate bool) {})
	if err != nil {
		panic(err)
	}

	var logger = New(&Config{
		BufferedSize: 100,
		Writer:       log.New(rotateLog, "\r\n", 0),
	})
	var data = []interface{}{
		"asdf", "ss", "sss",
	}

	logger.Debugf("%s\n[info] "+"asdf", append([]interface{}{"aaaaa"}, data...)...)
	for i := 0; i < 10; i++ {
		logger.Debugf("count %d \n", i)
		logger.Debug("count sssss", i, "asdfasdf")
		time.Sleep(time.Second)
	}
}

func TestLumperjackLogger(t *testing.T) {
	var writer = &lumberjack.Logger{
		Filename:   "foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	var logger = New(&Config{
		BufferedSize: 100,
		Writer:       log.New(writer, "\r\n", 0),
	})
	var data = []interface{}{
		"asdf", "ss", "sss",
	}

	logger.Debugf("%s\n[info] "+"asdf", append([]interface{}{"aaaaa"}, data...)...)
	for i := 0; i < 10; i++ {
		logger.Debugf("count %d \n", i)
		logger.Debug("count sssss", i, "asdfasdf")
		time.Sleep(time.Second)
	}
}
