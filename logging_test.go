package logging

import (
	"testing"
)

/////////////////////////////////////////////////
func TestLogging(t *testing.T) {
	Logger.Trace("trace msg")
	Logger.Debug("debug msg")
	Logger.Info("info msg")
	Logger.Warn("warn msg")
	Logger.Error("error msg")
	Logger.Fatal("fatal msg")
	Logger.Panic("panic msg")

}
