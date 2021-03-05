// Copyright 2020 NGR Softlab
//
// Logging - pack with common logger.
//
package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// Default logger
var Logger = NewLogger()

// Logrus logger creation.
func NewLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr), // for future - write to local journalctl and to remote syslog (for example throw "log/syslog")
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		},
		ReportCaller: true,
	}
	return logger
}

// Zap logger

// Smth else
