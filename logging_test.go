// Copyright 2020-2024 NGR Softlab
package logging

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogging(t *testing.T) {
	Logger.Trace("trace msg1")
	Logger.Debug("debug msg2")
	Logger.Info("info msg3")
	Logger.Warn("warn msg4")
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		expectedLevel logrus.Level
		logMessage    string
	}{
		{
			name:          "default logger settings",
			expectedLevel: logrus.DebugLevel,
			logMessage:    "Test debug message",
		},
		{
			name:          "check trace level logging",
			expectedLevel: logrus.TraceLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check info level logging",
			expectedLevel: logrus.InfoLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check warn level logging",
			expectedLevel: logrus.WarnLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check error level logging",
			expectedLevel: logrus.ErrorLevel,
			logMessage:    "Test error message",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				buf := new(bytes.Buffer)

				logger := NewLogger()
				logger.Out = buf
				logger.Level = tt.expectedLevel
				logger.ReportCaller = false

				switch tt.expectedLevel {
				case logrus.DebugLevel:
					logger.Debug(tt.logMessage)
				case logrus.TraceLevel:
					logger.Trace(tt.logMessage)
				case logrus.InfoLevel:
					logger.Info(tt.logMessage)
				case logrus.WarnLevel:
					logger.Warn(tt.logMessage)
				case logrus.ErrorLevel:
					logger.Error(tt.logMessage)
				default:
					panic("unhandled default case")
				}

				output := buf.String()

				if !strings.Contains(output, tt.logMessage) {
					t.Errorf("\ngot: %s\nwant: %s", output, tt.logMessage)
				}

				if !strings.Contains(output, strings.ToUpper(tt.expectedLevel.String())) {
					t.Errorf("\ngot: %s\nwant: %s", output, tt.logMessage)
				}
			},
		)
	}
}
