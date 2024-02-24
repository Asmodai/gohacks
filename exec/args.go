/*
 * args.go --- Program arguments.
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

/*
 * Future plans:
 *
 * A means to build dynamic arguments via reflection, so that one
 * does not have a base object with a 'port' thing that will be
 * useless with 99% of other use cases beyond the one I wrote this
 * entire shebang for.
 */

package exec

import (
	"fmt"
)

const (
	STRING_LIT_PORT string = "-port"
	STRING_FMT_PORT string = "%d"
)

type Args struct {
	args []string
	base int
}

func NewArgs(port int, args []string) Args {
	return Args{
		args: args,
		base: port,
	}
}

func (a Args) Get(port int) []string {
	return append(
		a.args,
		[]string{
			STRING_LIT_PORT,
			fmt.Sprintf(STRING_FMT_PORT, a.base+port),
		}...,
	)
}

/* args.go ends here. */
