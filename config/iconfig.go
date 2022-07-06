/*
 * iconfig.go --- Config interface.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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
	LogFile() string
}

/* iconfig.go ends here. */
