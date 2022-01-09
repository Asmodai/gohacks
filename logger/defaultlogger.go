/*
 * logger.go --- Default logger.
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
	"log"
)

type DefaultLogger struct {
}

func (l *DefaultLogger) SetDebug(junk bool) {
}

func (l *DefaultLogger) SetLogFile(junk string) {
}

func (l *DefaultLogger) Debug(msg string, rest ...interface{}) {
	log.Printf("DEBUG: %s  %v", msg, rest)
}

func (l *DefaultLogger) Warn(msg string, rest ...interface{}) {
	log.Printf("WARN: %s  %v", msg, rest)
}

func (l *DefaultLogger) Info(msg string, rest ...interface{}) {
	log.Printf("INFO: %s  %v", msg, rest)
}

func (l *DefaultLogger) Fatal(msg string, rest ...interface{}) {
	log.Fatalf("FATAL: %s  %v", msg, rest)
}

/* logger.go ends here. */
