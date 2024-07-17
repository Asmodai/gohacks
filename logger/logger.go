/* mock:yes */
/*
 * ilogger.go --- Logger interface.
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

/*
Logging structure.

To use,

1) Create a logger:

```go

	lgr := logger.NewLogger()

```

2) Do things with it:

```go

	lgr.Warn("Not enough coffee!")
	lgr.Info("Water is heating up.")
	// and so on.

```

If an empty string is passed to `NewLogger`, then the log facility will
display messages on standard output.
*/
type Logger interface {
	SetDebug(bool)
	SetLogFile(string)

	GoError(error, ...any)
	Debug(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
	Info(string, ...any)
	Fatal(string, ...any)
	Panic(string, ...any)

	Debugf(string, ...any)
	Errorf(string, ...any)
	Warnf(string, ...any)
	Infof(string, ...any)
	Fatalf(string, ...any)
	Panicf(string, ...any)

	WithFields(Fields) Logger
}

/* ilogger.go ends here. */
