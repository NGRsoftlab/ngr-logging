// Copyright 2020-2024 NGR Softlab
package logging

import (
	"testing"
)

func TestLogging(t *testing.T) {
	Logger.Trace("trace msg1")
	Logger.Debug("debug msg2")
	Logger.Info("info msg3")
	Logger.Warn("warn msg4")
}
