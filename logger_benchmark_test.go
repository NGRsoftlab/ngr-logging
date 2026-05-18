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
	logger := NewLogger("multicheck", "mp_scanner", "benchmark", "scanner-node", "127.0.0.1")
	logger.SetOutput(io.Discard)
	return logger
}

func newBenchmarkRawLogger() NgrZeroLogger {
	logger := NewLogger("multicheck", "mp_scanner", "benchmark", "scanner-node", "127.0.0.1")
	logger.SetRawOutput(io.Discard)
	return logger
}

func BenchmarkContextLoggerWithFields(b *testing.B) {
	logger := newBenchmarkLogger()
	err := errors.New("scanner error")
	var child NgrZeroLogger

	b.ReportAllocs()

	for b.Loop() {
		child = logger.With().
			Str("chain_id", "chain-1").
			Int("scanners_count", 3).
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
		logger.Info("scanner chain completed")
	}
}

func BenchmarkInfoPlainStringRawOutput(b *testing.B) {
	logger := newBenchmarkRawLogger()

	b.ReportAllocs()

	for b.Loop() {
		logger.Info("scanner chain completed")
	}
}

func BenchmarkDisabledDebugPlainString(b *testing.B) {
	baseLogger := newBenchmarkLogger()
	logger := baseLogger.Level(zerolog.ErrorLevel)

	b.ReportAllocs()

	for b.Loop() {
		logger.Debug("scanner chain completed")
	}
}
