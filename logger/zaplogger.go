/*
 * logger.go --- Logger facility.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gitlab.com/tozd/go/errors"

	"log"
)

type Fields map[string]any

// Logger using Zap as its implementation.
type zapLogger struct {
	logger   *zap.SugaredLogger
	logfile  string
	debug    bool
	facility string
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

// Set debug mode.
//
// Debug mode is a production-friendly runtime mode that will print
// human-readable messages to standard output instead of the defined
// log file.
func (l *zapLogger) SetDebug(flag bool) {
	var cfg = zap.NewProductionConfig()

	// If debug, use Zap's development config.
	//
	// This will result in textual logs rather than JSON.
	if flag {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set up an output path.  If there is none, or we are running
	// in debug mode, then output will be to standard output.
	if l.logfile != "" && !flag {
		cfg.OutputPaths = []string{l.logfile}
	}

	// Skip the first frame of the stack trace so we have the true
	// caller rather than our own functions.
	built, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Panic(err.Error())
	}

	l.logger = built.Sugar()
	l.debug = flag

	//nolint:errcheck
	l.logger.Sync()
}

// Set the log file to use.
func (l *zapLogger) SetLogFile(file string) {
	l.logfile = file

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

	nl := l.WithFields(detail)
	nl.Error(e.Error(), rest...)
}

// Write a debug message to the log.
func (l *zapLogger) Debug(msg string, rest ...any) {
	l.logger.Debugw(msg, rest...)
}

// Write an error message to the log.
func (l *zapLogger) Error(msg string, rest ...any) {
	l.logger.Errorw(msg, rest...)
}

// Write a warning message to the log.
func (l *zapLogger) Warn(msg string, rest ...any) {
	l.logger.Warnw(msg, rest...)
}

// Write an information message to the log.
func (l *zapLogger) Info(msg string, rest ...any) {
	l.logger.Infow(msg, rest...)
}

// Write a fatal message to the log and then exit.
func (l *zapLogger) Fatal(msg string, rest ...any) {
	l.logger.Fatalw(msg, rest...)
}

// Write a panic message to the log and then exit.
func (l *zapLogger) Panic(msg string, rest ...any) {
	l.logger.Panicw(msg, rest...)
}

// Compatibility method.
func (l *zapLogger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

// Compatibility method.
func (l *zapLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}

// Compatibility method.
func (l *zapLogger) Panicf(format string, args ...any) {
	l.logger.Panicf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var flds = make([]any, 0)
	for k, v := range fields {
		flds = append(flds, k)
		flds = append(flds, v)
	}

	return &zapLogger{
		logger:   l.logger.With(flds...),
		logfile:  l.logfile,
		debug:    l.debug,
		facility: l.facility,
	}
}

/* logger.go ends here. */
