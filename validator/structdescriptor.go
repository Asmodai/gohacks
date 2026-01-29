// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// structdescriptor.go --- Struct descriptor type.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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

package validator

// * Imports:

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Types:

type StructDescriptor struct {
	Type     reflect.Type
	Fields   map[string]*FieldInfo
	TypeName string
}

// ** Methods:

func (sd *StructDescriptor) Find(what string) (any, bool) {
	for elt := range sd.Fields {
		if strings.EqualFold(what, sd.Fields[elt].Name) {
			return sd.Fields[elt], true
		}
	}

	return nil, false
}

func (sd *StructDescriptor) Get(key string) (any, bool) {
	return sd.Find(key)
}

func (sd *StructDescriptor) Keys() []string {
	result := make([]string, 0, len(sd.Fields))

	for key := range sd.Fields {
		result = append(result, key)
	}

	sort.Strings(result)

	return result
}

func (sd *StructDescriptor) String() string {
	elts := make([]string, 0, len(sd.Fields))

	for key := range sd.Fields {
		elts = append(elts, sd.Fields[key].String())
	}

	return fmt.Sprintf("Fields[%s]", strings.Join(elts, " "))
}

func (sd *StructDescriptor) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug(sd.TypeName)

	dbg.Init(params...)
	dbg.Printf("Type name: %s", sd.TypeName)
	dbg.Printf("Fields")

	for _, key := range sd.Keys() {
		sd.Fields[key].Debug(&dbg)
	}

	dbg.End()

	dbg.Print()

	return dbg
}

// ** Functions:

func NewStructDescriptor() *StructDescriptor {
	return &StructDescriptor{
		Fields: make(map[string]*FieldInfo),
	}
}

// * structdescriptor.go ends here.
