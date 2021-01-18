package logger

import (
	"log"
	"testing"
	"time"

	"github.com/kjk/dailyrotate"
	"github.com/thaitanloi365/gocore/logger/notifier"

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
		Notifier: &notifier.SlackNotifier{
			WebhookURL: "https://hooks.slack.com/services/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",
			Channel:    "#legend-trucking-staging-bot",
		},
	})

	for i := 0; i < 10; i++ {
		logger.Printf("Test printf %s %d", "1231231", i)
		time.Sleep(time.Second)
	}

	return

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
