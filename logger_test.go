// Copyright 2020-2026 NGR Softlab

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

func TestConfiguredOutput(t *testing.T) {
	type testParams struct {
		fieldKey   string
		fieldValue string
		message    string
		level      zerolog.Level
		setOutput  func(*NgrZeroLogger, *bytes.Buffer)
	}

	tests := []struct {
		name   string
		params testParams
		want   []string
	}{
		{
			name: "writes console output debug with context field",
			params: testParams{
				fieldKey:   "chain_id",
				fieldValue: "chain-1",
				message:    "scanner chain completed",
				level:      zerolog.DebugLevel,
				setOutput: func(logger *NgrZeroLogger, buf *bytes.Buffer) {
					logger.SetOutput(buf)
				},
			},
			want: []string{
				"scanner chain completed",
				"chain_id=",
				"chain-1",
			},
		},
		{
			name: "writes raw output info with context field",
			params: testParams{
				fieldKey:   "chain_id",
				fieldValue: "chain-1",
				message:    "raw message",
				level:      zerolog.InfoLevel,
				setOutput: func(logger *NgrZeroLogger, buf *bytes.Buffer) {
					logger.SetRawOutput(buf)
				},
			},
			want: []string{
				`"message":"raw message"`,
				`"chain_id":"chain-1"`,
			},
		},
		{
			name: "writes raw output debug with context field",
			params: testParams{
				fieldKey:   "chain_id",
				fieldValue: "chain-1",
				message:    "raw message",
				level:      zerolog.DebugLevel,
				setOutput: func(logger *NgrZeroLogger, buf *bytes.Buffer) {
					logger.SetRawOutput(buf)
				},
			},
			want: []string{
				`"message":"raw message"`,
				`"chain_id":"chain-1"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)

			logger := NewLogger("", "", "", "", "")
			tt.params.setOutput(&logger, buf)

			child := logger.With().
				Str(tt.params.fieldKey, tt.params.fieldValue).
				Logger()

			switch tt.params.level {
			case zerolog.InfoLevel:
				child.Info(tt.params.message)
			case zerolog.DebugLevel:
				child.Debug(tt.params.message)
			default:
				t.Fatalf("unhandled level: %s", tt.params.level)
			}

			output := buf.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("\ngot: %s\nwant: %s", output, want)
				}
			}
		})
	}
}
