// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// rfc3339_test.go --- RFC3339 tests.
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

package rfc3339

import (
	"testing"
	"time"
)

var (
	TextObj1 string       = "2020-01-02T01:02:03Z"
	TextObj2 string       = "2020-01-02T01:02:03+00:00"
	TextObj3 string       = "2020-01-02T01:02:03"
	SqlObj   string       = "2020-01-02 01:02:03"
	UnixTime int64        = 1577926923
	JsonObj  *JSONRFC3339 = nil
)

func TestMarshal(t *testing.T) {
	JsonObj = &JSONRFC3339{}

	t.Log("Can we unmarshal from a JSON string?")
	err := JsonObj.UnmarshalJSON([]byte(TextObj1))
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err.Error())
		return
	}
	t.Log("Yes.")

	t.Log("Can we marshal to a JSON string?")
	str, err := JsonObj.MarshalJSON()
	if err != nil {
		t.Errorf("Marshal failed: %s", err.Error())
		return
	}
	t.Log("Yes.")

	t.Log("Are results identical?")
	if string(str) != "\""+TextObj1+"\"" {
		t.Errorf("No, '%s' != '%s'", string(str), TextObj1)
		return
	}
	t.Log("Yes.")
}

func TestInvalidMarshal(t *testing.T) {
	JsonObj = &JSONRFC3339{}

	t.Log("Do we get an error from invalid input?")
	err := JsonObj.UnmarshalJSON([]byte("This is not json"))
	if err == nil {
		t.Error("Invalid data did not generate an error.")
		return
	}

	t.Log("Yes.")
}

func TestTimeFuncs(t *testing.T) {
	t.Log("Do time functions work as expected?")

	JsonObj = &JSONRFC3339{}

	err := JsonObj.UnmarshalJSON([]byte(TextObj1))
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err.Error())
		return
	}

	utc := JsonObj.UTC()
	unix := JsonObj.Unix()
	golang := JsonObj.Time()
	sql := JsonObj.MySQL()

	if !golang.Equal(utc) {
		t.Error("Golang time.Time does not equal UTC time.Time.")
	}

	if unix != UnixTime {
		t.Error("Unix time does not equal computed timestamp.")
	}

	if sql != SqlObj {
		t.Error("MySQL timestamp is not the same.")
	}
}

func TestFutureTime(t *testing.T) {
	t.Log("Does the parser handle future times?")

	now := time.Now().UTC().Add(24 * time.Hour)

	res, err := Parse(now.Format("2006-01-02T15:04:05"))
	if err != nil {
		t.Errorf("No.  %s", err.Error())
		return
	}

	then := res.Format("2006-01-02T15:04:05Z")
	if then == time.Now().UTC().Format("2006-01-02T15:04:05Z") {
		t.Log("Yes.")
		return
	}

	t.Errorf(
		"No, %v != %v",
		then,
		time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	)
}

// rfc3339_test.go ends here.
