/*
 * logger.go --- Logger facility.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"log"
)

type Logger struct {
	logger   *zap.SugaredLogger
	logfile  string
	debug    bool
	facility string
}

func NewLogger(logfile string) *Logger {
	logger := &Logger{
		logger:  nil,
		logfile: logfile,
		debug:   false,
	}

	logger.SetDebug(false)

	return logger
}

func (l *Logger) SetDebug(flag bool) {
	var cfg zap.Config = zap.NewProductionConfig()

	if flag {
		cfg = zap.NewDevelopmentConfig()
	}

	if l.logfile != "" && !flag {
		cfg.OutputPaths = []string{l.logfile}
	}

	if flag {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	built, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Panic(err.Error())
	}

	l.logger = built.Sugar()
	l.debug = flag

	l.logger.Sync()
}

func (l *Logger) SetLogFile(file string) {
	l.logfile = file

	l.SetDebug(l.debug)
}

func (l *Logger) Debug(msg string, rest ...interface{}) {
	l.logger.Debugw(msg, rest...)
}

func (l *Logger) Warn(msg string, rest ...interface{}) {
	l.logger.Warnw(msg, rest...)
}

func (l *Logger) Info(msg string, rest ...interface{}) {
	l.logger.Infow(msg, rest...)
}

func (l *Logger) Fatal(msg string, rest ...interface{}) {
	l.logger.Fatalw(msg, rest...)
}

/* logger.go ends here. */
