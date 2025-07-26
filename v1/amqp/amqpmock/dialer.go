// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// dialer.go --- AMQP mock dialer.
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

// * Package:

package amqpmock

// * Imports:

import (
	"crypto/tls"

	"github.com/Asmodai/gohacks/v1/amqp/amqpshim"
)

// * Constants:

// * Variables:

// * Code:

// ** Type:

type MockDialer struct {
	URL        string
	Connection amqpshim.Connection
	Error      error
}

// ** Methods:

func (obj *MockDialer) Dial(url string) (amqpshim.Connection, error) {
	obj.URL = url

	return obj.Connection, obj.Error
}

func (obj *MockDialer) DialTLS(url string, _ *tls.Config) (amqpshim.Connection, error) {
	return obj.Dial(url)
}

// * dialer.go ends here.
