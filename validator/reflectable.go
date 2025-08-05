// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// reflectable.go --- `Reflectable` interface.
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

import "reflect"

// * Code:

type Reflectable interface {
	// Return the reflected type for a given object.
	//
	// An example of how this could work is:
	//
	// ```go
	//     var (
	//         typeForYourStruct reflect.Type
	//         onceForYourStruct sync.Once
	//     )
	//
	//     type YourStruct struct {
	//         // ...
	//     }
	//
	//     func (ys *YourStruct) ReflectType() reflect.Type {
	//         onceForYourStruct.Do(func() {
	//             typeForYourStruct = reflect.TypeOf(ys).Elem()
	//         })
	//
	//         return typeForYourStruct
	//     }
	// ```
	//
	// You might need to do things differently for non-pointer types.
	ReflectType() reflect.Type
}

// * reflectable.go ends here.
