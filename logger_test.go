// Copyright 2020-2024 NGR Softlab

package logging

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var _ Logger = NewLogger("", "", "", "", "")

func TestLogging(t *testing.T) {
	Log.Trace("trace msg1")
	Log.Debug("debug msg2")
	Log.Info("info msg3")
	Log.With().Dur("duration", 100*time.Nanosecond).Error("duration test")
	Log.Warn("warn msg4")
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		expectedLevel zerolog.Level
		logMessage    string
	}{
		{
			name:          "default logger settings",
			expectedLevel: zerolog.DebugLevel,
			logMessage:    "Test debug message",
		},
		{
			name:          "check trace level logging",
			expectedLevel: zerolog.TraceLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check info level logging",
			expectedLevel: zerolog.InfoLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check warn level logging",
			expectedLevel: zerolog.WarnLevel,
			logMessage:    "Test info message",
		},
		{
			name:          "check error level logging",
			expectedLevel: zerolog.ErrorLevel,
			logMessage:    "Test error message",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				buf := new(bytes.Buffer)

				logger := NewLogger("", "", "", "", "")
				logger.SetOutput(buf)

				switch tt.expectedLevel {
				case zerolog.DebugLevel:
					logger.Debug(tt.logMessage)
				case zerolog.TraceLevel:
					logger.Trace(tt.logMessage)
				case zerolog.InfoLevel:
					logger.Info(tt.logMessage)
				case zerolog.WarnLevel:
					logger.Warn(tt.logMessage)
				case zerolog.ErrorLevel:
					logger.Error(tt.logMessage)
				default:
					panic("unhandled default case")
				}

				output := buf.String()

				if !strings.Contains(output, tt.logMessage) {
					t.Errorf("\ngot: %s\nwant: %s", output, tt.logMessage)
				}

				if !strings.Contains(output, zerolog.FormattedLevels[tt.expectedLevel]) {
					t.Errorf("\ngot: %s\nwant: %s", output, tt.logMessage)
				}
			},
		)
	}
}
