// Copyright 2020-2024 NGR Softlab
package logging

import (
	"testing"
)

func TestLogging(t *testing.T) {
	Logger.Trace("trace msg")
	Logger.Debug("debug msg")
	Logger.Info("info msg")
	Logger.Warn("warn msg")
}
