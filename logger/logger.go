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

var logger = NewLogger()

func init() {
	// ?
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
	Package   string
	method    string
	trackerID string
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

// NewPackage initializes the standard logger package
func NewPackage(Package string) *StandardLogger {
	standardLogger := NewLogger()
	standardLogger.Package = Package

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

// Start is a standard start
func (logger *StandardLogger) Start(c echo.Context) *StandardLogger {
	if c != nil {
		logger.method = c.Request().Method
		logger.username, _ = c.Get("username").(string)
		logger.route = c.Path()
	}

	return logger
}

// StartTrackerID is a standard Start Tracker ID
func (logger *StandardLogger) StartTrackerID(c echo.Context) (string, *StandardLogger) {
	logger.Start(c)
	trackerID := uuid.NewString()
	logger.trackerID = trackerID

	return trackerID, logger
}

// End is a standard end
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

//

func (logger *StandardLogger) Info(argumentName ...interface{}) {
	logger.withFields().Info(argumentName)
}

func (logger *StandardLogger) Infof(format string, argumentName ...interface{}) {
	logger.withFields().Infof(format, argumentName)
}

func (logger *StandardLogger) Warn(argumentName ...interface{}) {
	logger.withFields().Warn(argumentName)
}

func (logger *StandardLogger) Warnf(format string, argumentName ...interface{}) {
	logger.withFields().Warnf(format, argumentName)
}
