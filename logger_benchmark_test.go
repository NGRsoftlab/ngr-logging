// Copyright © NGR Softlab 2025-2026

package logging

import (
	"errors"
	"io"
	"runtime"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func newBenchmarkLogger() NgrZeroLogger {
	logger := NewLogger("product", "component", "benchmark", "hostname", "127.0.0.1")
	logger.SetOutput(io.Discard)
	return logger
}

func newBenchmarkRawLogger() NgrZeroLogger {
	logger := NewLogger("product", "component", "benchmark", "hostname", "127.0.0.1")
	logger.SetRawOutput(io.Discard)
	return logger
}

func BenchmarkContextLoggerWithFields(b *testing.B) {
	logger := newBenchmarkLogger()
	err := errors.New("scan error")
	var child NgrZeroLogger

	b.ReportAllocs()

	for b.Loop() {
		child = logger.With().
			Str("chain_id", "chain-1").
			Int("count", 3).
			Dur("duration", time.Second).
			Err(err).
			Logger()
	}

	runtime.KeepAlive(child)
}

func BenchmarkInfoPlainString(b *testing.B) {
	logger := newBenchmarkLogger()

	b.ReportAllocs()

	for b.Loop() {
		logger.Info("chain completed")
	}
}

func BenchmarkInfoPlainStringRawOutput(b *testing.B) {
	logger := newBenchmarkRawLogger()

	b.ReportAllocs()

	for b.Loop() {
		logger.Info("chain completed")
	}
}

func BenchmarkDisabledDebugPlainString(b *testing.B) {
	baseLogger := newBenchmarkLogger()
	logger := baseLogger.Level(zerolog.ErrorLevel)

	b.ReportAllocs()

	for b.Loop() {
		logger.Debug("chain completed")
	}
}
