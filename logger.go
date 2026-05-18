// Copyright © NGR Softlab 2026

package logging

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
	"go.uber.org/atomic"
)

const consoleKey = "console"

type NgrZeroLogger struct {
	logger *zerolog.Logger

	outputs map[string]io.Writer

	Product   string
	Component string
	Version   string
	Hostname  string
	Address   string

	ShowProduct   atomic.Bool
	ShowComponent atomic.Bool
	ShowVersion   atomic.Bool
	ShowHostname  atomic.Bool
	ShowAddress   atomic.Bool

	skipFrameCount int // frames to additionally skip in getCaller. Needed to report caller if logger is wrapped
}

// Log - Default logger
var Log = NewLogger("", "", "", "", "")

var partsOrder = []string{"time", "level", "component", "version", "hostname", "address", "caller", "message"}
var fieldsExclude = []string{"component", "version", "hostname", "address"}

func NewLogger(product, component, version, hostname, address string) NgrZeroLogger {
	ngrLog := NgrZeroLogger{
		Product:   product,
		Component: component,
		Version:   version,
		Hostname:  hostname,
		Address:   address,
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	consoleOut := zerolog.ConsoleWriter{
		Out:                   os.Stdout,
		TimeFormat:            time.DateTime,
		PartsOrder:            partsOrder,
		FieldsExclude:         fieldsExclude,
		FormatPartValueByName: formatPartByName,
		FormatCaller:          formatCaller(ngrLog.skipFrameCount),
	}

	log := zerolog.New(os.Stdout).Output(consoleOut)

	ngrLog.outputs = map[string]io.Writer{consoleKey: consoleOut}

	// we skip additional caller frame here (default skipFrameCount is 2) because we wrap zerolog's methods (add 1 caller to the stack)
	ctx := log.With().Timestamp().CallerWithSkipFrameCount(3)
	// ctx = ctx.Str("product", ngrLog.Product)
	ctx = ctx.Str("component", ngrLog.Component)
	ctx = ctx.Str("version", ngrLog.Version)
	ctx = ctx.Str("hostname", ngrLog.Hostname)
	ctx = ctx.Str("address", ngrLog.Address)

	log = ctx.Logger()
	// log = log.With().Stack().Logger()

	ngrLog.logger = &log

	return ngrLog
}

func (l *NgrZeroLogger) SetOutput(w io.Writer) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:                   w,
		TimeFormat:            time.DateTime,
		PartsOrder:            partsOrder,
		FieldsExclude:         fieldsExclude,
		FormatPartValueByName: formatPartByName,
		FormatCaller:          formatCaller(l.skipFrameCount),
	}

	l.outputs = make(map[string]io.Writer)
	l.outputs[consoleKey] = consoleWriter

	l.logger = new(l.logger.Output(consoleWriter))
}

// SetRawOutput - switch logger output to the writer without console formatting.
func (l *NgrZeroLogger) SetRawOutput(w io.Writer) {
	l.outputs = make(map[string]io.Writer)
	l.outputs[consoleKey] = w

	l.logger = new(l.logger.Output(w))
}

func (l *NgrZeroLogger) AddOutput(key string, w io.Writer) {
	nw := diode.NewWriter(w, 1000, 0, func(missed int) {
		fmt.Printf("%s dropped %d messages\n", key, missed)
	})
	l.outputs[key] = nw
	l.updateOutputs()
}

// updateOutputs - refresh zerolog output after the configured writers change.
func (l *NgrZeroLogger) updateOutputs() {
	l.logger = new(l.logger.Output(l.output()))
}

// output - build one writer from all configured outputs.
func (l *NgrZeroLogger) output() io.Writer {
	switch len(l.outputs) {
	case 0:
		return io.Discard
	case 1:
		for _, w := range l.outputs {
			return w
		}
	}

	writers := make([]io.Writer, 0, len(l.outputs))
	for _, w := range l.outputs {
		writers = append(writers, w)
	}

	return io.MultiWriter(writers...)
}

// RemoveOutput removes output by key
func (l *NgrZeroLogger) RemoveOutput(key string) {
	if _, exist := l.outputs[key]; !exist {
		return
	}

	delete(l.outputs, key)

	l.updateOutputs()
}

func formatPartByName(i interface{}, s string) string {
	ret := fmt.Sprintf("%s", i)
	if len(ret) == 0 {
		return ""
	}

	switch s {
	case "hostname":
		ret = "(" + ret
	case "address":
		ret = "[" + ret + "])"
	case "component":
	case "version":
		ret = "v" + ret
	}

	return ret
}
