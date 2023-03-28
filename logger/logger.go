package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/google/uuid"
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

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	logger    *logrus.Logger
	tag       string
	trackerID string
	method    string
	username  string
	route     string
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger *logrus.Logger = logrus.New()
	standardLogger := &StandardLogger{
		logger: baseLogger,
	}
	standardLogger.logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	successArgMessage      = Event{1, "Success argument: %s"}
	warningArgMessage      = Event{2, "Warning argument: %s"}
	invalidArgMessage      = Event{3, "Invalid argument: %s"}
	invalidArgValueMessage = Event{4, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{5, "Missing argument: %s"}
)

// SetContext is a standard SetContext
func (logger *StandardLogger) SetContext(c echo.Context) {
	logger.method = c.Request().Method
	logger.username, _ = c.Get("username").(string)
	logger.route = c.Path()
}

func (logger *StandardLogger) SetTrackerID() {
	logger.trackerID = uuid.NewString()
}

func (logger *StandardLogger) withFields() *logrus.Entry {
	var fields logrus.Fields = logrus.Fields{}

	caller, function := fileNameAndfuncName()
	fields["caller"] = caller
	fields["function"] = function

	if logger.method != "" {
		fields["method"] = logger.method
	}

	if logger.username != "" {
		fields["username"] = logger.username
	}

	if logger.trackerID != "" {
		fields["trackerID"] = logger.trackerID
	}

	if logger.route != "" {
		fields["route"] = logger.route
	}

	return logger.logger.WithFields(
		fields,
	)
}

func (logger *StandardLogger) GetUsername() string {
	return logger.username
}

func fileNameAndfuncName() (string, string) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "", ""
	}

	fileName := fmt.Sprintf("%v:%v", path.Base(file), line)
	funcName := runtime.FuncForPC(pc).Name()
	function := funcName[strings.LastIndex(funcName, ".")+1:]
	return fileName, function
}

// SuccessArg is a standard success message
func (logger *StandardLogger) SuccessArg(argumentName string) {
	logger.withFields().Infof(successArgMessage.message, argumentName)
}

// WarningArg is a standard warning message
func (logger *StandardLogger) WarningArg(argumentName string) {
	logger.withFields().Warnf(warningArgMessage.message, argumentName)
}

// InvalidArg is a standard error message
func (logger *StandardLogger) InvalidArg(argumentName string) {
	logger.withFields().Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func (logger *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	logger.withFields().Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (logger *StandardLogger) MissingArg(argumentName string) {
	logger.withFields().Errorf(missingArgMessage.message, argumentName)
}

// logrus.Info: is a standard info message level
func (logger *StandardLogger) Info(args ...interface{}) {
	logger.withFields().Info(args...)
}

// logrus.Warn: is a standard warn message level
func (logger *StandardLogger) Warn(args ...interface{}) {
	logger.withFields().Warn(args...)
}

// logrus.Error: is a standard error message level
func (logger *StandardLogger) Error(args ...interface{}) {
	logger.withFields().Error(args...)
}

// logrus.Fatal: is a standard fatal message level
func (logger *StandardLogger) Fatal(args ...interface{}) {
	logger.withFields().Fatal(args...)
}

// logrus.Panic: is a standard panic message level
func (logger *StandardLogger) Panic(args ...interface{}) {
	logger.withFields().Panic(args...)
}

// logrus.Infof: is a standard info format message level
func (logger *StandardLogger) Infof(format string, args ...interface{}) {
	logger.withFields().Infof(format, args...)
}

// logrus.Warnf: is a standard warn format message level
func (logger *StandardLogger) Warnf(format string, args ...interface{}) {
	logger.withFields().Warnf(format, args...)
}

// logrus.Errorf: is a standard error format message level
func (logger *StandardLogger) Errorf(format string, args ...interface{}) {
	logger.withFields().Errorf(format, args...)
}

// logrus.Fatalf: is a standard fatal format message level
func (logger *StandardLogger) Fatalf(format string, args ...interface{}) {
	logger.withFields().Fatalf(format, args...)
}

// logrus.Panicf: is a standard panic message level
func (logger *StandardLogger) Panicf(format string, args ...interface{}) {
	logger.withFields().Panicf(format, args...)
}
