/*
 * iconfig.go --- Config interface.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
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

package config

import (
	"flag"
)

/*
Config interface.
*/
type IConfig interface {
	IsDebug() bool
	AddValidator(name string, fn interface{})
	String() string
	Validate() []error
	AddBoolFlag(p *bool, name string, value bool, usage string)
	AddFloat64Flag(p *float64, name string, value float64, usage string)
	AddIntFlag(p *int, name string, value int, usage string)
	AddInt64Flag(p *int64, name string, value int64, usage string)
	AddStringFlag(p *string, name string, value string, usage string)
	AddUintFlag(p *uint, name string, value uint, usage string)
	AddUint64Flag(p *uint64, name string, value uint64, usage string)
	LookupFlag(name string) *flag.Flag
	Parse()
}

/* iconfig.go ends here. */
