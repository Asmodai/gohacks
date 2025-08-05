// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// reflect.go --- Reflection.
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
	"reflect"
	"time"
)

// * Constants:

const (
	KindPrimitive FieldKind = iota
	KindStruct
	KindSlice
	KindMap
	KindUnknown
)

// * Variables:

var (
	//nolint:gochecknoglobals
	kindMap = map[FieldKind]string{
		KindPrimitive: "Primitive",
		KindStruct:    "Structure",
		KindSlice:     "Slice",
		KindMap:       "Map",
		KindUnknown:   "Unknown",
	}

	//nolint:gochecknoglobals
	baseTypeTimeType = reflect.TypeOf(time.Time{})
)

// * Code:

// ** Types:

type FieldKind int

// ** Functions:

func KindToString(kind FieldKind) string {
	if val, ok := kindMap[kind]; ok {
		return val
	}

	return ""
}

func isDerivedFrom(t, base reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t == base || t.ConvertibleTo(base)
}

//nolint:exhaustive
func classifyField(typ reflect.Type) FieldKind {
	switch typ.Kind() {
	case reflect.Struct:
		switch {
		case isDerivedFrom(typ, baseTypeTimeType):
			return KindPrimitive

		default:
			return KindStruct
		}

	case reflect.Slice:
		return KindSlice

	case reflect.Map:
		return KindMap

	case reflect.Ptr:
		return classifyField(typ.Elem())

	default:
		return KindPrimitive
	}
}

func reflectField(index int, field reflect.StructField) *FieldInfo {
	if len(field.PkgPath) > 0 {
		return nil
	}

	var (
		elemType reflect.Type
		elemKind FieldKind
	)

	kind := classifyField(field.Type)
	if kind == KindSlice || kind == KindMap {
		elemType = field.Type.Elem()
		elemKind = classifyField(elemType)
	}

	accessor := func(instance any) any {
		return reflect.ValueOf(instance).
			Elem().
			Field(index).
			Interface()
	}

	return &FieldInfo{
		Name:        field.Name,
		Type:        field.Type,
		TypeName:    field.Type.String(),
		TypeKind:    field.Type.Kind(),
		Accessor:    accessor,
		Tags:        field.Tag,
		Kind:        kind,
		ElementType: elemType,
		ElementKind: elemKind,
	}
}

func BuildDescriptor(typ reflect.Type) *StructDescriptor {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	desc := NewStructDescriptor()
	desc.Type = typ
	desc.TypeName = typ.Name()

	for idx := range typ.NumField() {
		field := typ.Field(idx)

		finfo := reflectField(idx, field)
		if finfo == nil {
			continue
		}

		desc.Fields[field.Name] = finfo
	}

	return desc
}

// * reflect.go ends here.
