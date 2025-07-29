// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// debug.go --- Debug flag in DI context.
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

// * Comments:

//
//
//

// * Package:

package contextdi

// * Imports:

import (
	"context"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ContextKeyDebugMode string = "_DI_FLG_DEBUG"
)

// * Variables:

// * Code:

// ** Functions

// Set the debug mode flag in the DI context to the given value.
func SetDebugMode(ctx context.Context, debugMode bool) (context.Context, error) {
	val, err := PutToContext(ctx, ContextKeyDebugMode, debugMode)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}

// Get the debug mode flag from the DI context.
func GetDebugMode(ctx context.Context) (bool, error) {
	val, err := GetFromContext(ctx, ContextKeyDebugMode)
	if err != nil {
		return false, errors.WithStack(err)
	}

	debugMode, ok := val.(bool)
	if !ok {
		return false, nil
	}

	return debugMode, nil
}

// * debug.go ends here.
