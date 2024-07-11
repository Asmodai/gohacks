/*
 * reflect.go --- Reflection hacks.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package config

import (
	"github.com/Asmodai/gohacks/types"

	"fmt"
	"reflect"
)

// Call a given function with arguments and return nil or an error.
//
// This is hairy reflection.
func (obj *config) call(field string, name string, params ...any) error {
	if _, ok := obj.Validators[name]; !ok {
		return types.NewError(
			"CONFIG",
			"Validator '%s' was not found.",
			name,
		)
	}
	fn := reflect.ValueOf(obj.Validators[name])

	// Check function arity.
	if len(params) != fn.Type().NumIn() {
		return types.NewError(
			"CONFIG",
			"Validator '%s' expects %d arguments, but only %d were given.",
			name,
			fn.Type().NumIn(),
			len(params),
		)
	}

	// Build funcall params.
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	// Funcall!
	result := fn.Call(in)
	if result[0].Interface() == nil {
		// Call was successful.
		return nil
	}

	// Fallthrough... we got an error.
	return types.NewError(
		"CONFIG",
		"[%s] Validation failed on '%s': %v",
		name,
		field,
		result[0].Interface().(error).Error(),
	)
}

// Attempt to locate a method and then invoke it.
//
// Will attempt to resolve methods called on pointers.
//
// If the method is found, the result will be returned along with `true`.
// Otherwise nil will be returned along with `false`.
func (obj *config) callMethod(value reflect.Value, method string) (any, bool) {
	var final reflect.Value
	var ptr reflect.Value

	// If we're a pointer then use the value of the pointee.
	if value.Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	}

	// Are we valid?
	if value.IsValid() {
		meth := value.MethodByName(method)

		// Better check if the method is valid too.
		if meth.IsValid() {
			final = meth
		}
	}

	// Are we a valid pointer?
	if ptr.IsValid() {
		meth := ptr.MethodByName(method)

		// Better check if the method is valid too.
		if meth.IsValid() {
			final = meth
		}
	}

	// Finally, double-check the method and then invoke it.
	if final.IsValid() {
		return final.Call([]reflect.Value{})[0].Interface(), true
	}

	// Nope, nothing found.
	return nil, false
}

// Attempt to call an `Init` method on a specific thing.
//
// Returns true if the call was successful, otherwise false.
// Discards any value returned from `Init`.
func (obj *config) checkCanInit(val reflect.Value) bool {
	_, ok := obj.callMethod(val, "Init")

	return ok
}

// Recursively pretty-print some value.
//
// `prefix` contains an arbitrary string that is printed before any element.  It
// is intended that this value be composed of spaces.
//
// `val` is the value (atom, list, struct, whatever) that we intend to print.
//
// `visited` is an accumulator that contains a map of pointers that we have
// visited.  Things that are consitered 'visited' will result in no further
// processing of that thing.
func (obj *config) recursePrint(prefix string, val reflect.Value, visited map[any]bool) string {
	var s string = ""
	toString, toStringFound := obj.callMethod(val, "ToString")

	// Reflect over pointers and interfaces.
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.Kind() == reflect.Ptr {
			// If we're a pointer, then check if we've visited the pointee.
			if visited[val.Interface()] {
				return s
			}

			// Tag it as visited.
			visited[val.Interface()] = true
		}

		// Get the pointee.
		val = val.Elem()
	}

	// Figure out what to do now.
	switch val.Kind() {
	case reflect.Struct: // Structure.
		if toStringFound {
			s += fmt.Sprintf("%s%s", prefix, toString.(string))
		} else {
			t := val.Type()

			for i := 0; i < val.NumField(); i++ {
				if t.Field(i).Tag.Get("config_hide") == "true" {
					// Ignore fields with the `config_hide` tag set to `true`.
					continue
				}

				s += fmt.Sprintf("\n%s%s:", prefix, t.Field(i).Name)

				// Is the field exported?
				if !val.Field(i).CanSet() {
					// No, mark it so and ignore it.
					s += " <unexported>"
					continue
				}

				// Should we obscure the field's value?
				if t.Field(i).Tag.Get("config_obscure") == "true" {
					s += " [**********]"
				} else {
					s += obj.recursePrint(prefix+"    ", val.Field(i), visited)
				}
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			s += obj.recursePrint("\n"+prefix, val.Index(i), visited)
		}

	case reflect.Invalid:
		s += " nil"

	default:
		s += fmt.Sprintf(" [%v]", val.Interface())
	}

	return s
}

// Recursive ugly reflective validation.
//
// Will invoke any validation function that is set via the `config_validator`
// tag.
func (obj *config) recurseValidate(v reflect.Value) []error {
	sref := v
	errs := []error{}

	for i := 0; i < sref.NumField(); i++ {
		field := sref.Field(i)
		ftype := sref.Type().Field(i)
		validate := ftype.Tag.Get("config_validator")

		// Nested structure?
		if field.Kind() == reflect.Struct {
			// Yep, recurse.
			nested := reflect.ValueOf(field.Interface())
			errs = append(errs, obj.recurseValidate(nested)...)
		}

		// Is validation function valid?
		if validate != "" {
			result := obj.call(ftype.Name, validate, field.Interface())
			if result != nil {
				errs = append(errs, []error{result}...)
			}
		}
	}

	return errs
}

/* reflect.go ends here. */
