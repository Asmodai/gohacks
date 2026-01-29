// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// client_test.go --- API client tests.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
//

// * Package:

package apiclient

// * Imports:

import (
	"github.com/Asmodai/gohacks/logger"

	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// * Code:

// ** Types:

// Callback type for fake HTTP magic.
type FakeHttpFn func(*http.Request) Response

// Fake HTTP magic structure.
type FakeHttp struct {
	Payload FakeHttpFn
}

// ** Methods:

// Perform a fake HTTP magic request thing, FEEL THE POWER!
func (c *FakeHttp) Do(req *http.Request) (*http.Response, error) {
	result := c.Payload(req)

	body := io.NopCloser(bytes.NewReader(result.Body))
	resp := &http.Response{
		StatusCode: result.StatusCode,
		Body:       body,
	}

	return resp, result.Error
}

// ** Functions:

// Create default parameters
func defaultParams() *Params {
	p := NewParams()
	p.URL = "http://127.0.0.1/test"

	return p
}

// Invoke afake HTTP magic GET request.
func invokeGet(params *Params, payloadfn FakeHttpFn) Response {
	c := NewClient(NewDefaultConfig(), logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.Get(params)
}

// Invoke afake HTTP magic POST request.
func invokePost(params *Params, payloadfn FakeHttpFn) Response {
	c := NewClient(NewDefaultConfig(), logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.Post(params)
}

// Invoke a fake HTTP magic request with a custom HTTP method verb.
func invokeVerb(verb string, params *Params, payloadfn FakeHttpFn) Response {
	conf := &Config{
		RequestsPerSecond: defaultRequestsPerSecond,
		Timeout:           defaultTimeout,
	}

	c := NewClient(conf, logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.(*client).httpAction(context.TODO(), verb, params)
}

// ** Tests:

// Test the 'Get' method.
func TestGet(t *testing.T) {
	// Basic run, no errors.
	t.Run(
		"Works as expected",
		func(t *testing.T) {
			payload := []byte("{}")
			resp := invokeGet(
				defaultParams(),
				func(_ *http.Request) Response {
					return NewResponse(
						200,
						payload,
						http.Header{},
						nil)
				},
			)

			if resp.Error != nil {
				t.Errorf("No, %s", resp.Error.Error())
			}

			if resp.StatusCode != 200 {
				t.Errorf("No, code %v", resp.StatusCode)
			}

			if !bytes.Equal(resp.Body, payload) {
				t.Error("No, payload does not match.")
			}
		},
	)

	// Basic run, server-side error.
	t.Run(
		"Handles server-side errors",
		func(t *testing.T) {
			resp := invokeGet(
				defaultParams(),
				func(_ *http.Request) Response {
					return NewResponse(
						0,
						[]byte("busted"),
						http.Header{},
						fmt.Errorf("broken"))
				},
			)

			if resp.Body != nil || resp.StatusCode != 0 || resp.Error.Error() != "broken" {
				t.Errorf(
					"No, unexpected data='%v', code='%v', err='%v'",
					resp.Body,
					resp.StatusCode,
					resp.Error.Error(),
				)
			}
		},
	)

	// Basic run, invalid HTTP method verb.
	t.Run(
		"Handles invalid methods",
		func(t *testing.T) {
			resp := invokeVerb(
				"Do The Thing",
				defaultParams(),
				func(_ *http.Request) Response {
					return NewResponse(
						200,
						[]byte(""),
						http.Header{},
						nil,
					)
				},
			)

			if resp.Error == nil {
				t.Error("No, no error returned.")
				return
			}

			if resp.Error.Error() != "net/http: invalid method \"Do The Thing\"" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)

	// Test if headers are set when present in the params.
	t.Run(
		"Sets headers if needed",
		func(t *testing.T) {
			accept := "text/gibberish"
			contentType := "text/shenanigans"
			resp := invokeGet(
				&Params{
					URL: "http://127.0.0.1/test",
					Content: ContentType{
						Accept: accept,
						Type:   contentType,
					},
				},
				func(req *http.Request) Response {
					payload := []byte("")

					if val := req.Header.Get("Accept"); val != accept {
						return NewResponseWithCodeFromError(
							500,
							fmt.Errorf("Invalid Accept hdr: '%v'", val))
					}

					if val := req.Header.Get("Content-Type"); val != contentType {
						return NewResponseWithCodeFromError(
							500,
							fmt.Errorf("Invalid Content-Type hdr: '%v'", val))
					}

					return NewResponse(
						200,
						payload,
						http.Header{},
						nil,
					)
				},
			)

			if resp.StatusCode == 500 {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)

	// Test if client complains if multiple auth types are requested.
	t.Run(
		"Complains if both auth types are used",
		func(t *testing.T) {
			resp := invokeGet(
				&Params{
					UseBasic: true,
					UseToken: true,
					URL:      "http://127.0.0.1/test",
				},
				func(_ *http.Request) Response {
					return NewResponse(
						200,
						[]byte(""),
						http.Header{},
						nil)
				},
			)

			if resp.Error == nil {
				t.Error("No, no error returned.")
				return
			}

			if resp.Error.Error() != "cannot use basic auth and token at the same time" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)

	t.Run(
		"Complains if basic auth is missing a username",
		func(t *testing.T) {
			resp := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseBasic: true,
				},
				func(_ *http.Request) Response {
					return Response{StatusCode: 200}
				},
			)

			if resp.Error == nil {
				t.Error("No, no error returned.")
				return
			}

			if resp.Error.Error() != "no basic auth username given" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)

	// Test if client complains if token auth lacks a header.
	t.Run(
		"Complains if token auth is missing headers",
		func(t *testing.T) {
			resp := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseToken: true,
				},
				func(_ *http.Request) Response {
					return Response{StatusCode: 200}
				},
			)

			if resp.Error == nil {
				t.Error("No, no error returned.")
				return
			}

			if resp.Error.Error() != "no auth token header given" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)

	// Test if queries are appended properly.
	t.Run(
		"Appends query parameters properly",
		func(t *testing.T) {
			params := &Params{
				URL: "http://127.0.0.1/test",
			}
			params.AddQueryParam("test", "value")

			resp := invokeGet(
				params,
				func(req *http.Request) Response {
					return NewResponse(
						500,
						[]byte(""),
						http.Header{},
						fmt.Errorf("%v", req.URL.RawQuery))
				},
			)

			if resp.Error == nil {
				t.Error("Unexpected result!")
				return
			}

			if resp.Error.Error() != "test=value" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)
}

// POST tests.
func TestPost(t *testing.T) {
	// Can we POST?
	t.Run(
		"Works as expected",
		func(t *testing.T) {
			payload := []byte("woot")
			resp := invokePost(
				defaultParams(),
				func(_ *http.Request) Response {
					return NewResponse(
						200,
						payload,
						http.Header{},
						nil,
					)
				},
			)

			if !bytes.Equal(resp.Body, payload) || resp.StatusCode != 200 || resp.Error != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'",
					resp.Body,
					resp.StatusCode,
					resp.Error.Error(),
				)
			}
		},
	)
}

// Token auth tests.
func TestTokenAuth(t *testing.T) {
	// Do auth token headers get appended?
	t.Run(
		"Does auth token header get added if required",
		func(t *testing.T) {
			payload := []byte("TOKEN")
			resp := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseToken: true,
					Token: AuthToken{
						Header: "HEADER",
						Data:   "TOKEN",
					},
				},
				func(req *http.Request) Response {
					tok := req.Header.Get("HEADER")

					return NewResponse(
						200,
						[]byte(tok),
						req.Header,
						nil,
					)
				},
			)

			if !bytes.Equal(resp.Body, payload) || resp.StatusCode != 200 || resp.Error != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'",
					resp.Body,
					resp.StatusCode,
					resp.Error.Error(),
				)
			}
		},
	)
}

// Basic auth tests.
func TestBasicAuth(t *testing.T) {
	// Are the basic auth details set in the request?
	t.Run(
		"Is basic auth added when required",
		func(t *testing.T) {
			payload := []byte("user:pass")
			resp := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseBasic: true,
					Basic: AuthBasic{
						Username: "user",
						Password: "pass",
					},
				},
				func(req *http.Request) Response {
					u, p, ok := req.BasicAuth()
					if !ok {
						return NewResponseWithCodeFromError(
							500,
							fmt.Errorf("basic auth"),
						)
					}

					return NewResponse(
						200,
						[]byte(u+":"+p),
						http.Header{},
						nil,
					)
				},
			)

			if !bytes.Equal(resp.Body, payload) || resp.StatusCode != 200 || resp.Error != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'",
					resp.Body,
					resp.StatusCode,
					resp.Error.Error(),
				)
			}
		},
	)
}

func TestStatusCode(t *testing.T) {
	t.Run(
		"Non-200 status codes should generate errors",
		func(t *testing.T) {
			resp := invokeGet(
				defaultParams(),
				func(req *http.Request) Response {
					return NewResponse(
						404,
						[]byte("404 Not Found"),
						http.Header{},
						nil,
					)
				},
			)

			if resp.Error == nil {
				t.Error("No, no error generated.")
				return
			}

			if resp.Error.Error() != "received status code 404 for http://127.0.0.1/test" {
				t.Errorf("No, '%v'", resp.Error.Error())
			}
		},
	)
}

// * client_test.go ends here.
