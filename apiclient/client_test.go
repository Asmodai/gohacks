/*
 * client_test.go --- API client tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package apiclient

import (
	"github.com/Asmodai/gohacks/logger"

	"bytes"
	"fmt"
	"io/ioutil"
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

	body := ioutil.NopCloser(bytes.NewReader(data))
	resp := &http.Response{
		StatusCode: code,
		Body:       body,
	}

	return resp, err
}

// Create default parameters
func defaultParams() *Params {
	p := NewParams()
	p.Url = "http://127.0.0.1/test"

	return p
}

// Invoke afake HTTP magic GET request.
func invokeGet(params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	conf := NewDefaultConfig()
	client := NewClient(conf, logger.NewDefaultLogger())
	client.Client = &FakeHttp{
		Payload: payloadfn,
	}

	return client.Get(params)
}

// Invoke afake HTTP magic POST request.
func invokePost(params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	conf := NewDefaultConfig()
	client := NewClient(conf, logger.NewDefaultLogger())
	client.Client = &FakeHttp{
		Payload: payloadfn,
	}

	return client.Post(params)
}

// Invoke a fake HTTP magic request with a custom HTTP method verb.
func invokeVerb(verb string, params *Params, payloadfn FakeHttpFn) ([]byte, int, error) {
	conf := &Config{
		RequestsPerSecond: 5,
		Timeout:           5,
	}

	client := NewClient(conf, logger.NewDefaultLogger())
	client.Client = &FakeHttp{
		Payload: payloadfn,
	}

	return client.httpAction(verb, params)
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

			if bytes.Compare(data, payload) != 0 {
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
					return []byte("busted"), 0, fmt.Errorf("Broken")
				},
			)

			if data != nil || code != 0 || err.Error() != "APICLIENT: Broken" {
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

			if err.Error() != "APICLIENT: net/http: invalid method \"Do The Thing\"" {
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
					Url: "http://127.0.0.1/test",
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
					Url:      "http://127.0.0.1/test",
				},
				func(_ *http.Request) ([]byte, int, error) {
					return []byte(""), 200, nil
				},
			)

			if err == nil {
				t.Error("No, no error returned.")
				return
			}

			if err.Error() != "APICLIENT: Cannot use Basic Auth and token at the same time." {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	t.Run(
		"Complains if basic auth is missing a username",
		func(t *testing.T) {
			_, _, err := invokeGet(
				&Params{
					Url:      "http://127.0.0.1/test",
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

			if err.Error() != "APICLIENT: No basic auth username given!" {
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
					Url:      "http://127.0.0.1/test",
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

			if err.Error() != "APICLIENT: No auth token header given!" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)

	// Test if queries are appended properly.
	t.Run(
		"Appends query parameters properly",
		func(t *testing.T) {
			params := &Params{
				Url: "http://127.0.0.1/test",
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

			if err.Error() != "APICLIENT: test=value" {
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

			if bytes.Compare(data, payload) != 0 || code != 200 || err != nil {
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
					Url:      "http://127.0.0.1/test",
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

			if bytes.Compare(data, payload) != 0 || code != 200 || err != nil {
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
					Url:      "http://127.0.0.1/test",
					UseBasic: true,
					Basic: AuthBasic{
						Username: "user",
						Password: "pass",
					},
				},
				func(req *http.Request) ([]byte, int, error) {
					u, p, ok := req.BasicAuth()
					if !ok {
						return []byte(""), 500, fmt.Errorf("Basic auth")
					}

					return []byte(u + ":" + p), 200, nil
				},
			)

			if bytes.Compare(data, payload) != 0 || code != 200 || err != nil {
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

			if err.Error() != "APICLIENT: Received status code 404 for http://127.0.0.1/test" {
				t.Errorf("No, '%v'", err.Error())
			}
		},
	)
}

/* client_test.go ends here. */
