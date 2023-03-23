package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var Logger = New()

func init() {
	// debugStr := config.GetAPPConfig().Debug
	// log := logrus.New()
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02T15:04:05.9999999Z07:00",
	// })
	// log.SetLevel(logrus.InfoLevel)
	// log.SetReportCaller(true)
	// log.SetOutput(os.Stdout)
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	DataKey: "caller",
	// 	CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
	// 		fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
	// 		//return frame.Function, fileName
	// 		return "", fileName
	// 	},
	// })

	// if debugStr == "true" {
	// 	logrus.SetLevel(logrus.DebugLevel)
	// }
}

type logger struct {
}

func New() *logger {
	return &logger{}
}

func LogEntry(c echo.Context) *logrus.Entry {
	if c == nil {
		return logrus.WithFields(logrus.Fields{
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return logrus.WithFields(logrus.Fields{
		"time":   time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

func (logger) fileNameAndfuncName() (string, string) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "", ""
	}

	fileName := fmt.Sprintf("%v:%v", path.Base(file), line)
	funcName := runtime.FuncForPC(pc).Name()
	function := funcName[strings.LastIndex(funcName, ".")+1:]
	return fileName, function
}

func (logger *logger) Error(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Error(args...)
}

func (logger *logger) Fatal(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Fatal(args...)
}

func (logger *logger) Warn(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Warn(args...)
}

func (logger *logger) Warnf(format string, args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Warnf(format, args...)
}

func (logger *logger) Info(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Info(args...)
}
