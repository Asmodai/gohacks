/*
 * module.go --- Golang plugin/module support.
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

package modules

import (
	"context"
	"errors"

	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/semver"
)

type StartFn func(context.Context, logger.ILogger) (bool, error)

type Module struct {
	name    string
	version *semver.SemVer
	startfn StartFn
}

func defaultStartFn(_ context.Context, _ logger.ILogger) (bool, error) {
	return false, errors.New("Module does not have a 'start' function!")
}

func NewModule(name string, version *semver.SemVer) *Module {
	return &Module{
		name:    name,
		version: version,
		startfn: defaultStartFn,
	}
}

func (m *Module) Name() string            { return m.name }
func (m *Module) Version() *semver.SemVer { return m.version }
func (m *Module) StartFn() StartFn        { return m.startfn }

func (m *Module) SetStartFn(fn StartFn) { m.startfn = fn }

func (m *Module) Start(ctx context.Context, lgr logger.ILogger) (bool, error) {
	if m.StartFn() == nil {
		return false, errors.New("Module start function is nil.")
	}

	return m.StartFn()(ctx, lgr)
}

/* module.go ends here. */
