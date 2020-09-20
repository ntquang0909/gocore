package logger

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Logger instance
type Logger struct {
	context    context.Context
	cancelFunc context.CancelFunc
	config     *Config

	mutex sync.RWMutex

	queue chan logTask

	writer     Writer
	fileWriter Writer

	debugStr      string
	debugColorStr string

	infoStr      string
	infoColorStr string

	warnStr      string
	warnColorStr string

	errStr      string
	errColorStr string
}

// Config log config
type Config struct {
	BufferedSize int
	Colorful     bool
	TimeLocation *time.Location
	DateFormat   string
	Prefix       string

	Writer                Writer
	WriteFileExceptLevels []LogLevel
}

// New new writter
func New(config *Config) *Logger {
	var bufferedSize = 10
	var dateFormat = "2006-01-02 15:04:05 Z07:00"
	timeLocation, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	var defaultConfig = config
	if defaultConfig == nil {
		defaultConfig = &Config{
			Prefix:       "",
			BufferedSize: bufferedSize,
			DateFormat:   dateFormat,
			TimeLocation: timeLocation,
			Colorful:     true,
		}
	}

	if defaultConfig.BufferedSize == 0 {
		defaultConfig.BufferedSize = bufferedSize
	}

	if defaultConfig.DateFormat == "" {
		defaultConfig.DateFormat = dateFormat
	}

	if defaultConfig.TimeLocation == nil {
		defaultConfig.TimeLocation = timeLocation
	}

	var writer = log.New(os.Stdout, "\r\n", 0)
	var fileWriter Writer = log.New(ioutil.Discard, "", 0)

	if defaultConfig.Writer != nil {
		fileWriter = defaultConfig.Writer
	}

	var (
		debugStr      = "%s DEBUG %s "
		infoStr       = "%s INFO %s "
		warnStr       = "%s WARN %s "
		errStr        = "%s ERROR %s "
		debugColorStr = "%s " + Green("DEBUG %s\n")
		infoColorStr  = "%s " + Blue("INFO %s\n")
		warnColorStr  = "%s " + Yellow("WARN %s\n")
		errColorStr   = "%s " + Red("ERROR %s\n")
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	var logger = &Logger{
		config:        defaultConfig,
		writer:        writer,
		fileWriter:    fileWriter,
		mutex:         sync.RWMutex{},
		context:       ctx,
		cancelFunc:    cancelFunc,
		queue:         make(chan logTask, defaultConfig.BufferedSize),
		debugStr:      debugStr,
		debugColorStr: debugColorStr,
		infoStr:       infoStr,
		infoColorStr:  infoColorStr,
		warnStr:       warnStr,
		warnColorStr:  warnColorStr,
		errStr:        errStr,
		errColorStr:   errColorStr,
	}

	logger.run()

	return logger
}

// Debug debug
func (l *Logger) Debug(values ...interface{}) {
	l.queue <- l.buildlog(Debug, l.fileWithLineNum(), valueTypeInterface, "", values...)
}

// Debugf debug with format
func (l *Logger) Debugf(format string, values ...interface{}) {
	l.queue <- l.buildlog(Debug, l.fileWithLineNum(), valueTypeInterface, format, values...)
}

// DebugJSON print pretty json
func (l *Logger) DebugJSON(values ...interface{}) {
	l.queue <- l.buildlog(Debug, l.fileWithLineNum(), valueTypeJSON, "", values...)
}

// Info info
func (l *Logger) Info(values ...interface{}) {
	l.queue <- l.buildlog(Info, l.fileWithLineNum(), valueTypeInterface, "", values...)
}

// Infof info with format
func (l *Logger) Infof(format string, values ...interface{}) {
	l.queue <- l.buildlog(Info, l.fileWithLineNum(), valueTypeInterface, format, values...)
}

// InfoJSON print pretty json
func (l *Logger) InfoJSON(values ...interface{}) {
	l.queue <- l.buildlog(Info, l.fileWithLineNum(), valueTypeJSON, "", values...)
}

// Warn warn
func (l *Logger) Warn(values ...interface{}) {
	l.queue <- l.buildlog(Warn, l.fileWithLineNum(), valueTypeInterface, "", values...)
}

// Warnf info with format
func (l *Logger) Warnf(format string, values ...interface{}) {
	l.queue <- l.buildlog(Warn, l.fileWithLineNum(), valueTypeInterface, format, values...)
}

// WarnJSON print pretty json
func (l *Logger) WarnJSON(values ...interface{}) {
	l.queue <- l.buildlog(Warn, l.fileWithLineNum(), valueTypeJSON, "", values...)
}

// Error error
func (l *Logger) Error(values ...interface{}) {
	l.queue <- l.buildlog(Error, l.fileWithLineNum(), valueTypeInterface, "", values...)
}

// Errorf error with format
func (l *Logger) Errorf(format string, values ...interface{}) {
	l.queue <- l.buildlog(Error, l.fileWithLineNum(), valueTypeInterface, format, values...)
}

// ErrorJSON print pretty json
func (l *Logger) ErrorJSON(values ...interface{}) {
	l.queue <- l.buildlog(Error, l.fileWithLineNum(), valueTypeJSON, "", values...)
}

func (l *Logger) run() {
	go l.cleanup()

	go func(ctx context.Context, queue chan logTask) {
		for {
			select {
			case <-ctx.Done():
				return

			case data := <-queue:
				var format = l.infoStr
				var formatColor = l.infoColorStr
				var extraFormat = data.format
				var extraPrettyFormat = data.format
				switch data.logLevel {
				case Debug:
					format = l.debugStr
					formatColor = l.debugColorStr
				case Error:
					format = l.errStr
					formatColor = l.errColorStr
				case Warn:
					format = l.warnStr
					formatColor = l.warnColorStr
				}

				var separator = " "
				switch data.valueType {
				case valueTypeJSON:
					separator = "\n"
				}

				if extraPrettyFormat == "" {
					for i := 0; i < len(data.values); i++ {
						extraPrettyFormat = "%v" + separator + extraPrettyFormat
					}
				}
				if extraFormat == "" {
					for i := 0; i < len(data.values); i++ {
						extraFormat = "%v" + " " + extraFormat
					}
				}

				switch data.valueType {
				case valueTypeJSON:
					var prettyValues = []interface{}{}
					var values = []interface{}{}
					for _, value := range data.values {
						values = append(values, ToJSONString(value))
						prettyValues = append(prettyValues, ToPrettyJSONString(value))
					}
					l.writer.Printf(formatColor+extraPrettyFormat, append([]interface{}{data.time, data.caller}, prettyValues...)...)
					if l.ignoreWriteFile(data.logLevel) == false {
						l.fileWriter.Printf(format+extraFormat, append([]interface{}{data.time, data.caller}, values...)...)
					}
				default:
					l.writer.Printf(formatColor+extraPrettyFormat, append([]interface{}{data.time, data.caller}, data.values...)...)
					if l.ignoreWriteFile(data.logLevel) == false {
						l.fileWriter.Printf(format+extraFormat, append([]interface{}{data.time, data.caller}, data.values...)...)
					}

				}

				break
			}
		}
	}(l.context, l.queue)
}

func (l *Logger) cleanup() {
	<-l.context.Done()

	// Lock the destinations
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Cleanup the destinations
	close(l.queue)

}

func (l *Logger) buildlog(logtype LogLevel, caller string, valueType valueType, format string, values ...interface{}) (newlog logTask) {
	newlog = logTask{
		logger:    l,
		logLevel:  logtype,
		time:      time.Now().Format(l.config.DateFormat),
		format:    format,
		values:    values,
		caller:    caller,
		valueType: valueType,
	}

	return newlog
}

func (l *Logger) fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		if ok {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

func (l *Logger) ignoreWriteFile(level LogLevel) bool {
	for _, lv := range l.config.WriteFileExceptLevels {
		if lv == level {
			return true
		}
	}

	return false

}
