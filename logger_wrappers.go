// Copyright © NGR Softlab 2025

package logging

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
)

// Level creates a child logger with the minimum accepted level set to level.
func (l *NgrZeroLogger) Level(lvl zerolog.Level) NgrZeroLogger {
	child := *l
	zl := child.logger.Level(lvl)
	child.logger = &zl
	return child
}

func (l NgrZeroLogger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

func (l NgrZeroLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l NgrZeroLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l NgrZeroLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l NgrZeroLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l NgrZeroLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

func (l NgrZeroLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

func (l NgrZeroLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}

func (l NgrZeroLogger) Trace(args ...interface{}) {
	l.logger.Trace().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Info(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Warning(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Error(args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l NgrZeroLogger) Panic(args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(args...))
}

type Context struct {
	zc     zerolog.Context
	logger *NgrZeroLogger
}

func (l *NgrZeroLogger) With() Context {
	return Context{zc: l.logger.With(),
		logger: l,
	}
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() NgrZeroLogger {
	nl := *c.logger
	zl := c.zc.Logger()
	nl.logger = &zl
	return nl
}

// Fields is a helper function to use a map or slice to set fields using type assertion.
// Only map[string]interface{} and []interface{} are accepted. []interface{} must
// alternate string keys and arbitrary values, and extraneous ones are ignored.
func (c Context) Fields(fields interface{}) Context {
	return Context{zc: c.zc.Fields(fields),
		logger: c.logger,
	}
}

// Str adds the field key with val as a string to the logger context.
func (c Context) Str(key, val string) Context {
	return Context{zc: c.zc.Str(key, val),
		logger: c.logger,
	}
}

// Strs adds the field key with val as a string to the logger context.
func (c Context) Strs(key string, vals []string) Context {
	return Context{zc: c.zc.Strs(key, vals),
		logger: c.logger,
	}
}

// Stringer adds the field key with val.String() (or null if val is nil) to the logger context.
func (c Context) Stringer(key string, val fmt.Stringer) Context {
	return Context{zc: c.zc.Stringer(key, val),
		logger: c.logger,
	}
}

// Errs adds the field key with errs as an array of serialized errors to the
// logger context.
func (c Context) Errs(key string, errs []error) Context {
	return Context{zc: c.zc.Errs(key, errs),
		logger: c.logger,
	}
}

// Err adds the field "error" with serialized err to the logger context.
func (c Context) Err(err error) Context {
	return Context{zc: c.zc.Err(err),
		logger: c.logger,
	}
}

// Ctx adds the context.Context to the logger context. The context.Context is
// not rendered in the error message, but is made available for hooks to use.
// A typical use case is to extract tracing information from the
// context.Context.
func (c Context) Ctx(ctx context.Context) Context {
	return Context{zc: c.zc.Ctx(ctx),
		logger: c.logger,
	}
}

// Bool adds the field key with val as a bool to the logger context.
func (c Context) Bool(key string, b bool) Context {
	return Context{zc: c.zc.Bool(key, b),
		logger: c.logger,
	}
}

// Bools adds the field key with val as a []bool to the logger context.
func (c Context) Bools(key string, b []bool) Context {
	return Context{zc: c.zc.Bools(key, b),
		logger: c.logger,
	}
}

// Int adds the field key with i as a int to the logger context.
func (c Context) Int(key string, i int) Context {
	return Context{zc: c.zc.Int(key, i),
		logger: c.logger,
	}
}

// Ints adds the field key with i as a []int to the logger context.
func (c Context) Ints(key string, i []int) Context {
	return Context{zc: c.zc.Ints(key, i),
		logger: c.logger,
	}
}

// Int8 adds the field key with i as a int8 to the logger context.
func (c Context) Int8(key string, i int8) Context {
	return Context{zc: c.zc.Int8(key, i),
		logger: c.logger,
	}
}

// Ints8 adds the field key with i as a []int8 to the logger context.
func (c Context) Ints8(key string, i []int8) Context {
	return Context{zc: c.zc.Ints8(key, i),
		logger: c.logger,
	}
}

// Int16 adds the field key with i as a int16 to the logger context.
func (c Context) Int16(key string, i int16) Context {
	return Context{zc: c.zc.Int16(key, i),
		logger: c.logger,
	}
}

// Ints16 adds the field key with i as a []int16 to the logger context.
func (c Context) Ints16(key string, i []int16) Context {
	return Context{zc: c.zc.Ints16(key, i),
		logger: c.logger,
	}
}

// Int32 adds the field key with i as a int32 to the logger context.
func (c Context) Int32(key string, i int32) Context {
	return Context{zc: c.zc.Int32(key, i),
		logger: c.logger,
	}
}

// Ints32 adds the field key with i as a []int32 to the logger context.
func (c Context) Ints32(key string, i []int32) Context {
	return Context{zc: c.zc.Ints32(key, i),
		logger: c.logger,
	}
}

// Int64 adds the field key with i as a int64 to the logger context.
func (c Context) Int64(key string, i int64) Context {
	return Context{zc: c.zc.Int64(key, i),
		logger: c.logger,
	}
}

// Ints64 adds the field key with i as a []int64 to the logger context.
func (c Context) Ints64(key string, i []int64) Context {
	return Context{zc: c.zc.Ints64(key, i),
		logger: c.logger,
	}
}

// Uint adds the field key with i as a uint to the logger context.
func (c Context) Uint(key string, i uint) Context {
	return Context{zc: c.zc.Uint(key, i),
		logger: c.logger,
	}
}

// Uints adds the field key with i as a []uint to the logger context.
func (c Context) Uints(key string, i []uint) Context {
	return Context{zc: c.zc.Uints(key, i),
		logger: c.logger,
	}
}

// Uint8 adds the field key with i as a uint8 to the logger context.
func (c Context) Uint8(key string, i uint8) Context {
	return Context{zc: c.zc.Uint8(key, i),
		logger: c.logger,
	}
}

// Uints8 adds the field key with i as a []uint8 to the logger context.
func (c Context) Uints8(key string, i []uint8) Context {
	return Context{zc: c.zc.Uints8(key, i),
		logger: c.logger,
	}
}

// Uint16 adds the field key with i as a uint16 to the logger context.
func (c Context) Uint16(key string, i uint16) Context {
	return Context{zc: c.zc.Uint16(key, i),
		logger: c.logger,
	}
}

// Uints16 adds the field key with i as a []uint16 to the logger context.
func (c Context) Uints16(key string, i []uint16) Context {
	return Context{zc: c.zc.Uints16(key, i),
		logger: c.logger,
	}
}

// Uint32 adds the field key with i as a uint32 to the logger context.
func (c Context) Uint32(key string, i uint32) Context {
	return Context{zc: c.zc.Uint32(key, i),
		logger: c.logger,
	}
}

// Uints32 adds the field key with i as a []uint32 to the logger context.
func (c Context) Uints32(key string, i []uint32) Context {
	return Context{zc: c.zc.Uints32(key, i),
		logger: c.logger,
	}
}

// Uint64 adds the field key with i as a uint64 to the logger context.
func (c Context) Uint64(key string, i uint64) Context {
	return Context{zc: c.zc.Uint64(key, i),
		logger: c.logger,
	}
}

// Uints64 adds the field key with i as a []uint64 to the logger context.
func (c Context) Uints64(key string, i []uint64) Context {
	return Context{zc: c.zc.Uints64(key, i),
		logger: c.logger,
	}
}

// Float32 adds the field key with f as a float32 to the logger context.
func (c Context) Float32(key string, f float32) Context {
	return Context{zc: c.zc.Float32(key, f),
		logger: c.logger,
	}
}

// Floats32 adds the field key with f as a []float32 to the logger context.
func (c Context) Floats32(key string, f []float32) Context {
	return Context{zc: c.zc.Floats32(key, f),
		logger: c.logger,
	}
}

// Float64 adds the field key with f as a float64 to the logger context.
func (c Context) Float64(key string, f float64) Context {
	return Context{zc: c.zc.Float64(key, f),
		logger: c.logger,
	}
}

// Floats64 adds the field key with f as a []float64 to the logger context.
func (c Context) Floats64(key string, f []float64) Context {
	return Context{zc: c.zc.Floats64(key, f),
		logger: c.logger,
	}
}

// Time adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (c Context) Time(key string, t time.Time) Context {
	return Context{zc: c.zc.Time(key, t),
		logger: c.logger,
	}
}

// Times adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (c Context) Times(key string, t []time.Time) Context {
	return Context{zc: c.zc.Times(key, t),
		logger: c.logger,
	}
}

// Interface adds the field key with obj marshaled using reflection.
func (c Context) Interface(key string, i interface{}) Context {
	return Context{zc: c.zc.Interface(key, i),
		logger: c.logger,
	}
}

// Type adds the field key with val's type using reflection.
func (c Context) Type(key string, val interface{}) Context {
	return Context{zc: c.zc.Interface(key, val),
		logger: c.logger,
	}
}

// Any is a wrapper around Context.Interface.
func (c Context) Any(key string, i interface{}) Context {
	return c.Interface(key, i)
}

// IPAddr adds IPv4 or IPv6 Address to the context
func (c Context) IPAddr(key string, ip net.IP) Context {
	return Context{zc: c.zc.Interface(key, ip),
		logger: c.logger,
	}
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the context
func (c Context) IPPrefix(key string, pfx net.IPNet) Context {
	return Context{zc: c.zc.Interface(key, pfx),
		logger: c.logger,
	}
}

// MACAddr adds MAC address to the context
func (c Context) MACAddr(key string, ha net.HardwareAddr) Context {
	return Context{zc: c.zc.Interface(key, ha),
		logger: c.logger,
	}
}

func (c Context) Tracef(format string, args ...interface{}) {
	c.Logger().logger.Trace().Msgf(format, args...)
}

func (c Context) Debugf(format string, args ...interface{}) {
	c.Logger().logger.Debug().Msgf(format, args...)
}

func (c Context) Infof(format string, args ...interface{}) {
	c.Logger().logger.Info().Msgf(format, args...)
}

func (c Context) Warnf(format string, args ...interface{}) {
	c.Logger().logger.Warn().Msgf(format, args...)
}

func (c Context) Warningf(format string, args ...interface{}) {
	c.Logger().logger.Warn().Msgf(format, args...)
}

func (c Context) Errorf(format string, args ...interface{}) {
	c.Logger().logger.Error().Msgf(format, args...)
}

func (c Context) Fatalf(format string, args ...interface{}) {
	c.Logger().logger.Fatal().Msgf(format, args...)
}

func (c Context) Panicf(format string, args ...interface{}) {
	c.Logger().logger.Panic().Msgf(format, args...)
}

func (c Context) Trace(args ...interface{}) {
	c.Logger().logger.Trace().Msg(fmt.Sprint(args...))
}

func (c Context) Debug(args ...interface{}) {
	c.Logger().logger.Debug().Msg(fmt.Sprint(args...))
}

func (c Context) Info(args ...interface{}) {
	c.Logger().logger.Info().Msg(fmt.Sprint(args...))
}

func (c Context) Warn(args ...interface{}) {
	c.Logger().logger.Warn().Msg(fmt.Sprint(args...))
}

func (c Context) Warning(args ...interface{}) {
	c.Logger().logger.Warn().Msg(fmt.Sprint(args...))
}

func (c Context) Error(args ...interface{}) {
	c.Logger().logger.Error().Msg(fmt.Sprint(args...))
}

func (c Context) Fatal(args ...interface{}) {
	c.Logger().logger.Fatal().Msg(fmt.Sprint(args...))
}

func (c Context) Panic(args ...interface{}) {
	c.Logger().logger.Panic().Msg(fmt.Sprint(args...))
}
