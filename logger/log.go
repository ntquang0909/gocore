package logger

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

const RefErrorIDKey = "__ref_error_id"
const UserIDKey = "__userID"

type valueType string

var (
	valueTypeInterface valueType = "interface"
	valueTypeJSON      valueType = "json"
	valueTypeCustom    valueType = "custom"
)

type requestInfo struct {
	reqID      string
	status     int
	method     string
	uri        string
	userID     string
	refErrorID string
}
type logTask struct {
	logger      *Logger
	logLevel    LogLevel
	time        string
	format      string
	values      []interface{}
	caller      string
	valueType   valueType
	requestInfo *requestInfo
}

func (task *logTask) withRequestInfo(requestInfo *requestInfo) *logTask {
	task.requestInfo = requestInfo
	return task
}

func (task *logTask) formatRequestInfo() string {
	if task.requestInfo == nil {
		return ""
	}
	var info = task.requestInfo

	var extras = []string{}
	if info.userID != "" {
		extras = append(extras, info.userID)
	}

	if info.refErrorID != "" {
		extras = append(extras, info.refErrorID)
	}

	if len(extras) > 0 {
		return fmt.Sprintf("%s [%v] %d %s %s", info.reqID, strings.Join(extras, "::"), info.status, info.method, info.uri)
	}

	return fmt.Sprintf("%s %d %s %s", info.reqID, info.status, info.method, info.uri)
}

func (task *logTask) withEchoContext(c echo.Context) *logTask {
	var reqID = ""
	var status = -1
	var userID = ""
	var refErrorID = ""
	if id, ok := c.Get(UserIDKey).(string); ok {
		userID = id
	}

	if id, ok := c.Get(RefErrorIDKey).(string); ok {
		refErrorID = id
	}
	var req = c.Request()
	var res = c.Response()

	if res != nil {
		reqID = res.Header().Get(echo.HeaderXRequestID)
		status = res.Status
	}

	var reqInfo = &requestInfo{
		method:     req.Method,
		reqID:      reqID,
		status:     status,
		refErrorID: refErrorID,
		userID:     userID,
		uri:        req.RequestURI,
	}
	return task.withRequestInfo(reqInfo)
}
