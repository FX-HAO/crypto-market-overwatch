package log

import (
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()

	// Debug logs a message using DEBUG as log level.
	Debug  = logger.Debug
	Debugf = logger.Debugf
	// Info logs a message using INFO as log level.
	Info  = logger.Info
	Infof = logger.Infof
	// Warn logs a message using WARN as log level.
	Warn  = logger.Warn
	Warnf = logger.Warnf
	// Error logs a message using ERROR as log level.
	Error  = logger.Error
	Errorf = logger.Errorf
	// Fatal logs a message using FATAL as log level and followed by a call to os.Exit(1).
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
	// Panic logs a message using ERROR as log level and followed by a call to panic().
	Panic  = logger.Panic
	Panicf = logger.Panicf
)

var logFile = "crypto-market-overwatch.log"

func init() {
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.Formatter = formatter

	logger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.DebugLevel: logFile,
		logrus.InfoLevel:  logFile,
		logrus.WarnLevel:  logFile,
		logrus.ErrorLevel: logFile,
		logrus.FatalLevel: logFile,
		logrus.PanicLevel: logFile,
	}, formatter))
}

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

// WithFields adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

// New creates a new logger.
func New() *logrus.Logger {
	return logrus.New()
}

// GetLogger returns a global instance
func GetLogger() *logrus.Logger {
	return logger
}
