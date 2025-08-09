// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// fieldinfo.go --- Field information.
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

// * Package:

package validator

// * Imports:

import (
	"fmt"
	"reflect"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Types:

type FieldAccessorFn func(any) any

type FieldInfo struct {
	Type        reflect.Type
	ElementType reflect.Type
	Accessor    FieldAccessorFn
	Name        string
	TypeName    string
	Tags        reflect.StructTag
	TypeKind    reflect.Kind
	Kind        FieldKind
	ElementKind FieldKind
}

// ** Methods:

func (fi *FieldInfo) String() string {
	return fmt.Sprintf("%s:%s", fi.Name, fi.TypeName)
}

func (fi *FieldInfo) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug(fi.Name)

	dbg.Init(params...)
	dbg.Printf("Name:         %s", fi.Name)
	dbg.Printf("Kind:         %s", KindToString(fi.Kind))
	dbg.Printf("Type:         %s", fi.TypeName)
	dbg.Printf("Tags:         (%v)", fi.Tags)
	dbg.Printf("")

	if fi.ElementType != nil {
		dbg.Printf("Element kind: %s", KindToString(fi.ElementKind))
		dbg.Printf("Element type: %s", fi.ElementType.Name())
	}

	dbg.End()

	dbg.Print()

	return dbg
}

// ** Functions:

// * fieldinfo.go ends here.
