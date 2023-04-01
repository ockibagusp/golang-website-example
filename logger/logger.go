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

var debugStr = config.GetAPPConfig().Debug

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	logger    *logrus.Logger
	Package   string
	method    string
	trackerID string
	username  string
	route     string
}

func new() *StandardLogger {
	var baseLogger *logrus.Logger = logrus.New()
	standardLogger := &StandardLogger{
		logger: baseLogger,
	}
	standardLogger.logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	standardLogger.logger.SetLevel(logrus.InfoLevel)
	standardLogger.logger.SetOutput(os.Stdout)
	if debugStr == "true" {
		standardLogger.logger.SetLevel(logrus.DebugLevel)
	}

	return standardLogger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	return new()
}

// NewPackage initializes the standard logger package
func NewPackage(Package string) (StandardLogger *StandardLogger) {
	StandardLogger = new()
	StandardLogger.Package = Package

	return
}

// Declare variables to store log messages as new Events
var (
	successArgMessage      = Event{1, "Success argument: %s"}
	warningArgMessage      = Event{2, "Warning argument: %s"}
	invalidArgMessage      = Event{3, "Invalid argument: %s"}
	invalidArgValueMessage = Event{4, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{5, "Missing argument: %s"}
)

// Start is a standard logger start
func (logger *StandardLogger) Start(c echo.Context) *StandardLogger {
	if c != nil {
		logger.method = c.Request().Method
		logger.username, _ = c.Get("username").(string)
		logger.route = c.Path()
	}

	return logger
}

// StartTrackerID is a standard logger start tracker id
func (logger *StandardLogger) StartTrackerID(c echo.Context) (string, *StandardLogger) {
	logger.Start(c)
	trackerID := uuid.NewString()
	logger.trackerID = trackerID

	return trackerID, logger
}

// SetTrackerID is a standard logger set tracker id
func (logger *StandardLogger) SetTrackerID() (trackerID string) {
	trackerID = uuid.NewString()
	logger.trackerID = trackerID

	return trackerID
}

// End is a standard logger end
func (logger *StandardLogger) End() {
	logger.method = ""
	logger.trackerID = ""
	logger.username = ""
	logger.route = ""
}

func (logger *StandardLogger) withFields() *logrus.Entry {
	var fields logrus.Fields = logrus.Fields{}

	if logger.Package != "" {
		fields["package"] = logger.Package
	}

	if logger.method != "" {
		fields["method"] = logger.method
	}

	if logger.trackerID != "" {
		fields["tracker_id"] = logger.trackerID
	}

	if logger.username != "" {
		fields["username"] = logger.username
	}

	if logger.route != "" {
		fields["route"] = logger.route
	}

	caller, function := fileNameAndfuncName()
	fields["caller"] = caller
	fields["function"] = function

	return logger.logger.WithFields(
		fields,
	)
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

// SuccessArg is a standard logger success message
func (logger *StandardLogger) SuccessArg(argumentName string) {
	logger.withFields().Infof(successArgMessage.message, argumentName)
}

// WarningArg is a standard logger warning message
func (logger *StandardLogger) WarningArg(argumentName string) {
	logger.withFields().Warnf(warningArgMessage.message, argumentName)
}

// InvalidArg is a standard logger error message
func (logger *StandardLogger) InvalidArg(argumentName string) {
	logger.withFields().Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard logger error message
func (logger *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	logger.withFields().Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard logger error message
func (logger *StandardLogger) MissingArg(argumentName string) {
	logger.withFields().Errorf(missingArgMessage.message, argumentName)
}

// logrus
func (logger *StandardLogger) Print(argumentName string) {
	logger.withFields().Print(argumentName)
}

func (logger *StandardLogger) Info(argumentName string) {
	logger.withFields().Info(argumentName)
}

func (logger *StandardLogger) Warn(argumentName string) {
	logger.withFields().Warn(argumentName)
}

func (logger *StandardLogger) Error(argumentName string) {
	logger.withFields().Error(argumentName)
}

func (logger *StandardLogger) Fatal(argumentName string) {
	logger.withFields().Fatal(argumentName)
}

func (logger *StandardLogger) Panic(argumentName string) {
	logger.withFields().Panic(argumentName)
}

// Entry Println family functions
func (logger *StandardLogger) Traceln(format string, argumentName ...interface{}) {
	logger.withFields().Traceln(format, argumentName)
}

func (logger *StandardLogger) Debugln(format string, argumentName ...interface{}) {
	logger.withFields().Debugln(format, argumentName)
}

func (logger *StandardLogger) Infoln(format string, argumentName ...interface{}) {
	logger.withFields().Infoln(format, argumentName)
}

func (logger *StandardLogger) Println(format string, argumentName ...interface{}) {
	logger.withFields().Infoln(format, argumentName)
}

func (logger *StandardLogger) Warnln(format string, argumentName ...interface{}) {
	logger.withFields().Warnln(format, argumentName)
}

func (logger *StandardLogger) Errorln(format string, argumentName ...interface{}) {
	logger.withFields().Errorln(format, argumentName)
}

func (logger *StandardLogger) Fatalln(format string, argumentName ...interface{}) {
	logger.withFields().Fatalln(format, argumentName)
}

func (logger *StandardLogger) Panicln(format string, argumentName ...interface{}) {
	logger.withFields().Panicln(format, argumentName)
}

// Entry Printf family functions
func (logger *StandardLogger) Tracef(format string, argumentName ...interface{}) {
	logger.withFields().Tracef(format, argumentName)
}

func (logger *StandardLogger) Debugf(format string, argumentName ...interface{}) {
	logger.withFields().Debugf(format, argumentName)
}

func (logger *StandardLogger) Infof(format string, argumentName ...interface{}) {
	logger.withFields().Infof(format, argumentName)
}

func (logger *StandardLogger) Printf(format string, argumentName ...interface{}) {
	logger.withFields().Infof(format, argumentName)
}

func (logger *StandardLogger) Warnf(format string, argumentName ...interface{}) {
	logger.withFields().Warnf(format, argumentName)
}

func (logger *StandardLogger) Errorf(format string, argumentName ...interface{}) {
	logger.withFields().Errorf(format, argumentName)
}

func (logger *StandardLogger) Fatalf(format string, argumentName ...interface{}) {
	logger.withFields().Fatalf(format, argumentName)
}

func (logger *StandardLogger) Panicf(format string, argumentName ...interface{}) {
	logger.withFields().Panicf(format, argumentName)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (logger *StandardLogger) WithField(key string, value interface{}) *logrus.Entry {
	return logger.logger.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *StandardLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.logger.WithFields(fields)
}
