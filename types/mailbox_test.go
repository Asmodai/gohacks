// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// mailbox_test.go --- Mailbox tests.
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

package types

import (
	"testing"
	"time"
)

func TestMailboxNoContext(t *testing.T) {
	var mbox *Mailbox = nil

	//
	// Test creation.
	t.Run("Can create new mailbox", func(t *testing.T) {
		mbox = NewMailbox()

		if mbox.Full() {
			t.Error("Somehow the mailbox is already full!")
			return
		}
	})

	//
	// Test single write/read.
	t.Run("Can write and read in a single routine", func(t *testing.T) {
		mbox.Put("test")

		val, ok := mbox.Get()
		if !ok {
			t.Error("Get failed")
			return
		}

		if val.(string) != "test" {
			t.Errorf("Unexpected value '%v' returned.", val)
			return
		}
	})

	//
	// Test with goroutines.
	t.Run("Can write and read with concurrency", func(t *testing.T) {
		writeChan := make(chan *Pair, 1)
		readChan := make(chan *Pair, 1)

		go func(mb *Mailbox) {
			mb.Put("concurrency")
			writeChan <- NewPair(true, nil)
		}(mbox)

		go func(mb *Mailbox) {
			val, ok := mb.Get()
			for ok == false {
				val, ok = mb.Get()
				time.Sleep(50 * time.Millisecond)
			}

			readChan <- NewPair(val, ok)
		}(mbox)

		<-writeChan

		select {
		case readRes := <-readChan:
			{
				if readRes.Second.(bool) != true {
					t.Error("Problem with reading value!")
					return
				}

				if readRes.First.(string) != "concurrency" {
					t.Errorf("Unexpected value '%v' returned.", readRes.First)
					return
				}
			}

		case <-time.After(5 * time.Second):
			t.Error("Timeout after 5 seconds!")
			return
		}
	})
}

// mailbox_test.go ends here.
