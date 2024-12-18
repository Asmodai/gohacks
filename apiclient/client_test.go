// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// client_test.go --- API client tests.
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

package apiclient

import (
	"github.com/Asmodai/gohacks/logger"

	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// Callback type for fake HTTP magic.
type FakeHttpFn func(*http.Request) ([]byte, int, error)

// Fake HTTP magic structure.
type FakeHttp struct {
	Payload FakeHttpFn
}

// Perform a fake HTTP magic request thing, FEEL THE POWER!
func (c *FakeHttp) Do(req *http.Request) (*http.Response, error) {
	data, code, err := c.Payload(req)

	body := io.NopCloser(bytes.NewReader(data))
	resp := &http.Response{
		StatusCode: code,
		Body:       body,
	}

	return resp, err
}

// Create default parameters
func defaultParams() *Params {
	p := NewParams()
	p.URL = "http://127.0.0.1/test"

	return p
}

// Invoke afake HTTP magic GET request.
func invokeGet(params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	c := NewClient(NewDefaultConfig(), logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.Get(params)
}

// Invoke afake HTTP magic POST request.
func invokePost(params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	c := NewClient(NewDefaultConfig(), logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.Post(params)
}

// Invoke a fake HTTP magic request with a custom HTTP method verb.
func invokeVerb(verb string, params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	conf := &Config{
		RequestsPerSecond: 5,
		Timeout:           5,
	}

	c := NewClient(conf, logger.NewDefaultLogger())
	c.(*client).Client = &FakeHttp{
		Payload: payloadfn,
	}

	return c.(*client).httpAction(context.TODO(), verb, params)
}

// Test the 'Get' method.
func TestGet(t *testing.T) {
	// Basic run, no errors.
	t.Run(
		"Works as expected",
		func(t *testing.T) {
			payload := []byte("{}")
			data, code, err := invokeGet(
				defaultParams(),
				func(_ *http.Request) ([]byte, int, error) {
					return payload, 200, nil
				},
			)

			if err != nil {
				t.Errorf("No, %s", err.Error())
			}

			if code != 200 {
				t.Errorf("No, code %v", code)
			}

			if !bytes.Equal(data, payload) {
				t.Error("No, payload does not match.")
			}
		},
	)

	// Basic run, server-side error.
	t.Run(
		"Handles server-side errors",
		func(t *testing.T) {
			data, code, err := invokeGet(
				defaultParams(),
				func(_ *http.Request) ([]byte, int, error) {
					return []byte("busted"), 0, fmt.Errorf("broken")
				},
			)

			if data != nil || code != 0 || err.Error() != "broken" {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'", data, code, err.Error())
			}
		},
	)

	// Basic run, invalid HTTP method verb.
	t.Run(
		"Handles invalid methods",
		func(t *testing.T) {
			_, _, err := invokeVerb(
				"Do The Thing",
				defaultParams(),
				func(_ *http.Request) ([]byte, int, error) {
					return []byte(""), 200, nil
				},
			)

			if err == nil {
				t.Error("No, no error returned.")
				return
			}

			if err.Error() != "net/http: invalid method \"Do The Thing\"" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	// Test if headers are set when present in the params.
	t.Run(
		"Sets headers if needed",
		func(t *testing.T) {
			accept := "text/gibberish"
			contentType := "text/shenanigans"
			_, code, err := invokeGet(
				&Params{
					URL: "http://127.0.0.1/test",
					Content: ContentType{
						Accept: accept,
						Type:   contentType,
					},
				},
				func(req *http.Request) ([]byte, int, error) {
					payload := []byte("")

					if val := req.Header.Get("Accept"); val != accept {
						return payload, 500, fmt.Errorf("Invalid Accept hdr: '%v'", val)
					}

					if val := req.Header.Get("Content-Type"); val != contentType {
						return payload, 500, fmt.Errorf("Invalid Content-Type hdr: '%v'", val)
					}

					return payload, 200, nil
				},
			)

			if code == 500 {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	// Test if client complains if multiple auth types are requested.
	t.Run(
		"Complains if both auth types are used",
		func(t *testing.T) {
			_, _, err := invokeGet(
				&Params{
					UseBasic: true,
					UseToken: true,
					URL:      "http://127.0.0.1/test",
				},
				func(_ *http.Request) ([]byte, int, error) {
					return []byte(""), 200, nil
				},
			)

			if err == nil {
				t.Error("No, no error returned.")
				return
			}

			if err.Error() != "cannot use basic auth and token at the same time" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	t.Run(
		"Complains if basic auth is missing a username",
		func(t *testing.T) {
			_, _, err := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseBasic: true,
				},
				func(_ *http.Request) ([]byte, int, error) {
					return []byte(""), 200, nil
				},
			)

			if err == nil {
				t.Error("No, no error returned.")
				return
			}

			if err.Error() != "no basic auth username given" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	// Test if client complains if token auth lacks a header.
	t.Run(
		"Complains if token auth is missing headers",
		func(t *testing.T) {
			_, _, err := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseToken: true,
				},
				func(_ *http.Request) ([]byte, int, error) {
					return []byte(""), 200, nil
				},
			)

			if err == nil {
				t.Error("No, no error returned.")
				return
			}

			if err.Error() != "no auth token header given" {
				t.Errorf("No, '%v'", err.Error())
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

			_, _, err := invokeGet(
				params,
				func(req *http.Request) ([]byte, int, error) {
					return []byte(""), 500, fmt.Errorf("%v", req.URL.RawQuery)
				},
			)

			if err == nil {
				t.Error("Unexpected result!")
				return
			}

			if err.Error() != "test=value" {
				t.Errorf("No, '%v'", err.Error())
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
			data, code, err := invokePost(
				defaultParams(),
				func(_ *http.Request) ([]byte, int, error) {
					return payload, 200, nil
				},
			)

			if !bytes.Equal(data, payload) || code != 200 || err != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'", data, code, err.Error())
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
			data, code, err := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseToken: true,
					Token: AuthToken{
						Header: "HEADER",
						Data:   "TOKEN",
					},
				},
				func(req *http.Request) ([]byte, int, error) {
					tok := req.Header.Get("HEADER")

					return []byte(tok), 200, nil
				},
			)

			if !bytes.Equal(data, payload) || code != 200 || err != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'", data, code, err.Error())
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
			data, code, err := invokeGet(
				&Params{
					URL:      "http://127.0.0.1/test",
					UseBasic: true,
					Basic: AuthBasic{
						Username: "user",
						Password: "pass",
					},
				},
				func(req *http.Request) ([]byte, int, error) {
					u, p, ok := req.BasicAuth()
					if !ok {
						return []byte(""), 500, fmt.Errorf("basic auth")
					}

					return []byte(u + ":" + p), 200, nil
				},
			)

			if !bytes.Equal(data, payload) || code != 200 || err != nil {
				t.Errorf("No, unexpected data='%v', code='%v', err='%v'", data, code, err.Error())
			}
		},
	)
}

func TestStatusCode(t *testing.T) {
	t.Run(
		"Non-200 status codes should generate errors",
		func(t *testing.T) {
			_, _, err := invokeGet(
				defaultParams(),
				func(req *http.Request) ([]byte, int, error) {
					return []byte("404 Not Found"), 404, nil
				},
			)

			if err == nil {
				t.Error("No, no error generated.")
				return
			}

			if err.Error() != "received status code 404 for http://127.0.0.1/test" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)
}

// client_test.go ends here.
