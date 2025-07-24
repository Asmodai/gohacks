// -*- Mode: Go -*-
//
// types.go --- Worker types.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
//go:build amd64 || arm64 || riscv64

// * Comments:
//
//

// * Package:

package dynworker

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/logger"
)

// * Code:

// ** User data type:

// *** Type:

// User-supplied data.
type UserData any

// ** Task data type:

// *** Type:

// Task structure.
type Task struct {
	parent context.Context // Parent context.
	logger logger.Logger   // Logger instance.
	data   UserData        // User-supplied data.
}

// *** Functions:

func NewTask(ctx context.Context, lgr logger.Logger, data UserData) *Task {
	return &Task{
		parent: ctx,
		logger: lgr,
		data:   data,
	}
}

// *** Methods:

// Get the parent context for the task.
func (obj *Task) Parent() context.Context {
	return obj.parent
}

// Get the logger instance for the task.
func (obj *Task) Logger() logger.Logger {
	return obj.logger
}

// Get the user-supplied data for the task.
func (obj *Task) Data() UserData {
	return obj.data
}

// Reset the task.
func (obj *Task) reset() {
	// NOTE: setting `obj.Data` to `struct{}{}` might be better.
	obj.parent = nil
	obj.logger = nil
	obj.data = nil
}

// ** Task function type:

// Task callback function type.
type TaskFn func(*Task) error

// * types.go ends here.
