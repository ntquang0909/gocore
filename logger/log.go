package logger

type valueType string

var (
	valueTypeInterface valueType = "interface"
	valueTypeJSON      valueType = "json"
)

type logTask struct {
	logger    *Logger
	logLevel  LogLevel
	time      string
	format    string
	values    []interface{}
	caller    string
	valueType valueType
}
