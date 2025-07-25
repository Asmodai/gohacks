// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// debug_test.go --- Debug tool tests.
//
// Copyright (c) 2022-2025 Paul Ward <paul@lisphacker.uk>
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

package debug

// * Imports:

import (
	"testing"

	"github.com/google/uuid"
)

// * Code:

// ** Types:

type SomeNestedStruct struct {
	Nested string
}

func (obj SomeNestedStruct) Debug(params ...any) *Debug {
	dbg := NewDebug("Nested Structure")

	dbg.Init(params...)
	dbg.Printf("Nested: %s", obj.Nested)
	dbg.End()

	return dbg
}

type SomeStruct struct {
	One    int
	Two    string
	Three  uuid.UUID
	Nested SomeNestedStruct
}

func (obj SomeStruct) Debug(params ...any) *Debug {
	dbg := NewDebug("Some Structure")

	dbg.Init(params...)
	dbg.Printf("One:   %d", obj.One)
	dbg.Printf("Two:   %s", obj.Two)
	dbg.Printf("Three: %s", obj.Three)

	obj.Nested.Debug(&dbg)

	dbg.End()

	return dbg
}

// ** Functions:

func DebugGet(thing any, params ...any) *Debug {
	impl, ok := thing.(Debugable)
	if !ok {
		return nil
	}

	return impl.Debug(params...)
}

// ** Tests:

func TestFunctional(t *testing.T) {
	inst := &SomeStruct{
		One:   42,
		Two:   "Forty-two",
		Three: uuid.New(),
		Nested: SomeNestedStruct{
			Nested: "Yep",
		},
	}

	dbg := DebugGet(inst, "Some Structure")
	if dbg == nil {
		t.Fatal("Could not get debugging info")
	}

	t.Logf("Result:\n%s\n", dbg.String())
}

// * debug_test.go ends here.
