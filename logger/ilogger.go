/*
 * ilogger.go --- Logger interface.
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

type ILogger interface {
	SetDebug(bool)
	SetLogFile(string)

	Debug(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Fatal(string, ...interface{})

	Debugf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})

	WithFields(Fields) ILogger
}

/* ilogger.go ends here. */
