// Copyright © NGR Softlab 2020-2025

package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger - Default logger
var Logger = NewLogger()

// NewLogger Logrus logger creation.
func NewLogger() *logrus.Logger {
	// TODO: for future - write to local journalctl and to remote syslog (for example throw "log/syslog")
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr),
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
