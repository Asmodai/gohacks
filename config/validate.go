// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// validate.go --- Validators.
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

package config

import (
	"reflect"
	"slices"
)

// Add a validator function for a tag value.
//
// This is so one can add validators after instance creation.
func (c *config) AddValidator(name string, fn any) {
	c.Validators[name] = fn
}

// Validate configuration.
//
// Should validation fail, then a list of errors is returned.
// Should validation pass, an empty list is returned.
func (c *config) Validate() []error {
	sref := reflect.ValueOf(c.App).Elem()
	sval := reflect.ValueOf(c.App)

	// Deal with validators.
	res := c.recurseValidate(sref)

	// Now deal with `Validate` functions.
	if mthd := sval.MethodByName("Validate"); mthd.IsValid() {
		rval := mthd.Call(nil)
		if rerr := rval[0].Interface(); rerr != nil {
			if errs, ok := rerr.([]error); ok {
				res = slices.Concat(res, errs)
			}
		}
	}

	return res
}

// validate.go ends here.
