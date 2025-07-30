// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// zaplogger.go --- Logger using Zap as the implementation.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:

// * Package:

package logger

// * Imports:

import (
	"fmt"
	"log"
	"sync"

	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// * Constants:

const (
	defaultSamplingInitial    = 100 //nolint:unused
	defaultSamplingThereafter = 100 //nolint:unused
)

// * Code:

// ** Types:

// Log field map type.
type Fields map[string]any

// Logger using Zap as its implementation.
type zapLogger struct {
	mu sync.RWMutex

	logger   *zap.Logger
	logfile  string
	debug    bool
	facility string
}

// ** Methods:

// Set debug mode.
//
// Debug mode is a production-friendly runtime mode that will print
// human-readable messages to standard output instead of the defined log file.
func (l *zapLogger) SetDebug(flag bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var cfg zap.Config

	if flag {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.OutputPaths = []string{"stderr"}
	} else {
		cfg = zap.NewProductionConfig()

		if len(l.logfile) > 0 {
			cfg.OutputPaths = []string{l.logfile}
			cfg.ErrorOutputPaths = []string{l.logfile}
		}
	}

	cfg.Sampling = nil
	/* XXX This should be an optional for stuff that wants sampling.
	// Set up sampling.
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    defaultSamplingInitial,
		Thereafter: defaultSamplingThereafter,
	}
	*/

	// Skip the first frame of the stack trace so we have the true
	// caller rather than our own functions.
	built, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		log.Panic(err.Error())
	}

	l.logger = built
	l.debug = flag

	//nolint:errcheck
	l.logger.Sync()
}

func (l *zapLogger) SetSampling(_, _ int) {
	// Not wired in yet.
}

// Set the log file to use.
func (l *zapLogger) SetLogFile(file string) {
	l.mu.Lock()
	l.logfile = file
	l.mu.Unlock()

	l.SetDebug(l.debug)
}

// Write a Go error to the log.
func (l *zapLogger) GoError(err error, rest ...any) {
	e, ok := err.(errors.E) //nolint:errorlint,varnamelen
	if !ok {
		l.Error(err.Error(), rest...)

		return
	}

	detail := errors.Details(e)
	if detail == nil {
		detail = map[string]any{}
	}

	if cause := errors.Cause(e); cause != nil {
		detail["cause"] = cause
	}

	l.WithFields(detail).Error(e.Error(), rest...)
}

// Write a debug message to the log.
func (l *zapLogger) Debug(msg string, rest ...any) {
	l.logger.Debug(msg, buildFields(rest...)...)
}

// Write an error message to the log.
func (l *zapLogger) Error(msg string, rest ...any) {
	l.logger.Error(msg, buildFields(rest...)...)
}

// Write a warning message to the log.
func (l *zapLogger) Warn(msg string, rest ...any) {
	l.logger.Warn(msg, buildFields(rest...)...)
}

// Write an information message to the log.
func (l *zapLogger) Info(msg string, rest ...any) {
	l.logger.Info(msg, buildFields(rest...)...)
}

// Write a fatal message to the log and then exit.
func (l *zapLogger) Fatal(msg string, rest ...any) {
	l.logger.Fatal(msg, buildFields(rest...)...)
}

// Write a panic message to the log and then exit.
func (l *zapLogger) Panic(msg string, rest ...any) {
	l.logger.Panic(msg, buildFields(rest...)...)
}

// Compatibility method.
func (l *zapLogger) Debugf(format string, args ...any) {
	l.logger.Sugar().Debugf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Errorf(format string, args ...any) {
	l.logger.Sugar().Errorf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Warnf(format string, args ...any) {
	l.logger.Sugar().Warnf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Infof(format string, args ...any) {
	l.logger.Sugar().Infof(format, args...)
}

// Compatibility method.
func (l *zapLogger) Fatalf(format string, args ...any) {
	l.logger.Sugar().Fatalf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Panicf(format string, args ...any) {
	l.logger.Sugar().Panicf(format, args...)
}

// Encapsulate user-specified metadata fields.
func (l *zapLogger) WithFields(fields Fields) Logger {
	zapFields := make([]zap.Field, 0, len(fields))

	for key, val := range fields {
		zapFields = append(zapFields, toZapField(key, val))
	}

	return &zapLogger{
		logger:   l.logger.With(zapFields...),
		logfile:  l.logfile,
		debug:    l.debug,
		facility: l.facility,
	}
}

// * Functions:

func buildFields(kvs ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(kvs)/2) //nolint:mnd

	for idx := 0; idx < len(kvs)-1; idx += 2 {
		key, ok := kvs[idx].(string)
		if !ok {
			key = fmt.Sprintf("invalid_key_%d", idx)
		}

		fields = append(fields, toZapField(key, kvs[idx+1]))
	}

	return fields
}

func toZapField(key string, value any) zap.Field {
	switch val := value.(type) {
	case string:
		return zap.String(key, val)

	case int:
		return zap.Int(key, val)

	case int64:
		return zap.Int64(key, val)

	case float64:
		return zap.Float64(key, val)

	case bool:
		return zap.Bool(key, val)

	case error:
		return zap.NamedError(key, val)

	default:
		return zap.Any(key, value)
	}
}

// Create a new logger.
func NewZapLogger() Logger {
	return NewZapLoggerWithFile("")
}

// Create a new logger with the given log file.
func NewZapLoggerWithFile(logfile string) Logger {
	instance := &zapLogger{
		logger:  nil,
		logfile: logfile,
		debug:   false,
	}

	instance.SetDebug(false)

	return instance
}

// zaplogger.go ends here.
