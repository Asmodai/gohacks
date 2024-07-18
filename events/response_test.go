// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// response_test.go --- Response event tests.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package events

import (
	"fmt"
	"testing"
	"time"
)

func TestResponseEvent(t *testing.T) {
	command := 42
	data := "Test"
	reply := "Reply"

	msg := NewMessage(command, data)
	time.Sleep(1 * time.Second)
	rsp := NewResponse(msg, reply)

	t.Logf("Message:  %s\n", msg)
	t.Logf("Response: %s\n", rsp)

	t.Run("Accessors", func(t *testing.T) {
		if rsp.Index() != msg.Index() {
			t.Errorf("Index is wrong:  %d != %d", rsp.Index(), msg.Index())
		}

		if rsp.Command() != command {
			t.Error("Command is wrong.")
		}

		if rsp.Response().(string) != reply {
			t.Error("Response is wrong.")
		}

		dur := rsp.When().Sub(rsp.Received())
		str := fmt.Sprintf(
			"Response Event: index:%d duration:%s",
			rsp.Index(),
			dur.String(),
		)

		if rsp.String() != str {
			t.Errorf("String is wrong: '%s' != '%s'", rsp.String(), str)
		}
	})
}

// response_test.go ends here.
