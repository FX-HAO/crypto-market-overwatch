package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

var (
	logger = logrus.New()

	// DEBUG ...
	Debug  = logger.Debug
	Debugf = logger.Debugf
	// INFO ...
	Info  = logger.Info
	Infof = logger.Infof
	// WARNING ...
	Warn  = logger.Warn
	Warnf = logger.Warnf
	// ERROR ...
	Error  = logger.Error
	Errorf = logger.Errorf
	// FATAL ...
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
	// Panic ...
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

type Fields logrus.Fields

func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

func New() *logrus.Logger {
	return logrus.New()
}

func GetLogger() *logrus.Logger {
	return logger
}
