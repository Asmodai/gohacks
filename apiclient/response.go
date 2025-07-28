// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// response.go --- API client response structure.
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

package apiclient

// * Imports:

import (
	"net/http"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type Response struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
	Limits     RateLimitInfo
	Error      error
}

// ** Methods:

// ** Functions:

func NewResponse(code int, body []byte, headers http.Header, err error) Response {
	return Response{
		Body:       body,
		StatusCode: code,
		Headers:    headers,
		Limits:     NewRateLimitInfo(code, headers),
		Error:      err,
	}
}

func NewResponseFromError(err error) Response {
	return NewResponse(0, nil, http.Header{}, err)
}

func NewResponseWithCodeFromError(code int, err error) Response {
	return NewResponse(code, nil, http.Header{}, err)
}

// * response.go ends here.
