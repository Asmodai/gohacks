/*
 * logger.go --- Logger facility.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"log"
	"sync"
)

var (
	once     sync.Once
	instance *Logger
)

type Fields map[string]interface{}

/*

Logging structure.

To use,

1) Create a logger:

    lgr := logger.NewLogger("/path/to/log")

2) Do things with it:

    lgr.Warn("Not enough coffee!")
    lgr.Info("Water is heating up.")
    // and so on.

If an empty string is passed to `NewLogger`, then the log facility will
display messages on standard output.

*/
type Logger struct {
	logger   *zap.SugaredLogger
	logfile  string
	debug    bool
	facility string
}

// Initialise our singleton instance.
func initInstance(logfile string) {
	once.Do(func() {
		instance = &Logger{
			logger:  nil,
			logfile: logfile,
			debug:   false,
		}

		instance.SetDebug(false)
	})
}

// Create a new logger.
func NewLogger() *Logger {
	return NewLoggerWithFile("")
}

// Create a new logger with the given log file.
func NewLoggerWithFile(logfile string) *Logger {
	if instance == nil {
		initInstance(logfile)
	}

	return instance
}

// Set debug mode.
//
// Debug mode is a production-friendly runtime mode that will print
// human-readable messages to standard output instead of the defined
// log file.
func (l *Logger) SetDebug(flag bool) {
	var cfg zap.Config = zap.NewProductionConfig()

	// If debug, use Zap's development config.
	//
	// This will result in textual logs rather than JSON.
	if flag {
		cfg = zap.NewDevelopmentConfig()
	}

	// Set up an output path.  If there is none, or we are running
	// in debug mode, then output will be to standard output.
	if l.logfile != "" && !flag {
		cfg.OutputPaths = []string{l.logfile}
	}

	// Yeah, pweety colours for debug mode!
	if flag {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Skip the first frame of the stack trace so we have the true
	// caller rather than our own functions.
	built, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Panic(err.Error())
	}

	l.logger = built.Sugar()
	l.debug = flag

	l.logger.Sync()
}

// Set the log file to use.
func (l *Logger) SetLogFile(file string) {
	l.logfile = file

	l.SetDebug(l.debug)
}

// Write a debug message to the log.
func (l *Logger) Debug(msg string, rest ...interface{}) {
	l.logger.Debugw(msg, rest...)
}

// Write a warning message to the log.
func (l *Logger) Warn(msg string, rest ...interface{}) {
	l.logger.Warnw(msg, rest...)
}

// Write an information message to the log.
func (l *Logger) Info(msg string, rest ...interface{}) {
	l.logger.Infow(msg, rest...)
}

// Write a fatal message to the log and then exit.
func (l *Logger) Fatal(msg string, rest ...interface{}) {
	l.logger.Fatalw(msg, rest...)
}

// Compatibility method.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Compatibility method.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Compatibility method.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Compatibility method.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Compatibility method.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Compatibility method.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *Logger) WithFields(fields Fields) ILogger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.logger.With(f...)
	return &Logger{
		logger:   newLogger,
		logfile:  l.logfile,
		debug:    l.debug,
		facility: l.facility,
	}
}

/* logger.go ends here. */
