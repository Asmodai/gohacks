/*
 * mailbox_test.go --- Mailbox tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

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

/* mailbox_test.go ends here. */
