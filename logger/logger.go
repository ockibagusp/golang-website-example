package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/config"
	"github.com/sirupsen/logrus"
)

func init() {
	debugStr := config.GetAPPConfig().Debug
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	if debugStr == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

type Logger struct {
	trackerID string
}

func New() *Logger {
	return &Logger{}
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

func (Logger) fileNameAndfuncName() (string, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "", ""
	}

	fileName := fmt.Sprintf("%v:%v", path.Base(file), line)
	funcName := runtime.FuncForPC(pc).Name()
	function := funcName[strings.LastIndex(funcName, ".")+1:]
	return fileName, function
}

func (logger *Logger) SetTrackerID() {
	logger.trackerID = "123"
}

func (logger *Logger) Error(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Error(args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Errorf(format, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Fatal(args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Fatalf(format, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Warn(args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Warnf(format, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Info(args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	caller, function := logger.fileNameAndfuncName()
	logrus.WithFields(
		logrus.Fields{
			"caller":   caller,
			"function": function,
		},
	).Infof(format, args...)
}
