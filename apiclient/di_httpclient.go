// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// di_httpclient.go --- HTTP client DI.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package apiclient

// * Imports:

import (
	"context"

	contextdi "github.com/Asmodai/gohacks/contextdi"
	errx "github.com/Asmodai/gohacks/errx"
)

// * Constants:

const (
	// Key used to store the instance in the context's user value.
	ContextKeyHTTPClient = "gohacks/apiclient/HTTPClient@v1"
)

// * Variables:

var (

	// Signalled if the instance associated with the context key is not of
	// type HTTPClient.
	ErrValueNotHTTPClient = errx.Base("value is not HTTPClient")
)

// * Code:
// ** Functions:

// Set HTTPClient stores the instance in the context map.
func SetHTTPClient(ctx context.Context, inst HTTPClient) (context.Context, error) {
	val, err := contextdi.PutToContext(ctx, ContextKeyHTTPClient, inst)
	if err != nil {
		return nil, errx.WithStack(err)
	}

	return val, nil
}

// Get the instance from the given context.
//
// Will return ErrValueNotHTTPClient if the value in the context is not of type
// HTTPClient.
func GetHTTPClient(ctx context.Context) (HTTPClient, error) {
	var zero HTTPClient

	val, err := contextdi.GetFromContext(ctx, ContextKeyHTTPClient)
	if err != nil {
		return zero, errx.WithStack(err)
	}

	inst, ok := val.(HTTPClient)
	if !ok {
		return zero, errx.WithStack(ErrValueNotHTTPClient)
	}

	return inst, nil
}

// Attempt to get the instance from the given context.  Panics if the
// operation fails.
func MustGetHTTPClient(ctx context.Context) HTTPClient {
	inst, err := GetHTTPClient(ctx)

	if err != nil {
		panic(errx.WithMessage(err, "HTTPClient missing in context"))
	}

	return inst
}

// TryGetHTTPClient returns the instance and true if present and typed.
func TryGetHTTPClient(ctx context.Context) (HTTPClient, bool) {
	var zero HTTPClient

	val, err := contextdi.GetFromContext(ctx, ContextKeyHTTPClient)
	if err != nil {
		return zero, false
	}

	inst, ok := val.(HTTPClient)
	if !ok {
		return zero, false
	}

	return inst, true
}

// * di_httpclient.go ends here.
