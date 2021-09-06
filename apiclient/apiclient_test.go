/*
 * apiclient_test.go --- API client tests.
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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type FakeHttpFn func(*http.Request) (int, []byte, error)

type FakeHttp struct {
	Payload FakeHttpFn
}

func (c *FakeHttp) Do(req *http.Request) (*http.Response, error) {
	code, data, err := c.Payload(req)

	body := ioutil.NopCloser(bytes.NewReader(data))
	resp := &http.Response{
		StatusCode: code,
		Body:       body,
	}

	return resp, err
}

func TestGet(t *testing.T) {
	t.Log("Does 'Get' work as expected?")

	payload := []byte("{}")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(_ *http.Request) (int, []byte, error) {
			return 200, payload, nil
		},
	}

	params := &Params{
		Url: "http://127.0.0.1/test",
	}

	data, code, err := client.Get(params)
	if err != nil {
		t.Errorf("No, %s", err.Error())
		return
	}

	if code == 200 && bytes.Compare(data, payload) == 0 {
		t.Log("Yes.")
	} else {
		t.Errorf("No.")
	}
}

func TestPost(t *testing.T) {
	t.Log("Does 'Post' work as expected?")

	payload := []byte("{}")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(_ *http.Request) (int, []byte, error) {
			return 200, payload, nil
		},
	}

	params := &Params{
		Url: "http://127.0.0.1/test",
	}

	data, code, err := client.Post(params)
	if err != nil {
		t.Errorf("No, %s", err.Error())
		return
	}

	if code == 200 && bytes.Compare(data, payload) == 0 {
		t.Log("Yes.")
	} else {
		t.Errorf("No.")
	}
}

func TestAims(t *testing.T) {
	t.Log("Does 'X-Aimms-Auth-Token' get added if there is a token?")

	payload := []byte("TOKEN")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(req *http.Request) (int, []byte, error) {
			tok := req.Header.Get("X-Aims-Auth-Token")

			return 200, []byte(tok), nil
		},
	}

	params := &Params{
		Url:   "http://127.0.0.1/test",
		Token: "TOKEN",
	}

	data, code, err := client.Get(params)
	if err != nil {
		t.Errorf("No, %s", err.Error())
		return
	}

	if code == 200 && bytes.Compare(data, payload) == 0 {
		t.Log("Yes.")
	} else {
		t.Errorf("No.")
	}
}

func TestAuth(t *testing.T) {
	t.Log("Is basic auth added when username/password exists?")

	payload := []byte("user:pass")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(req *http.Request) (int, []byte, error) {
			u, p, ok := req.BasicAuth()
			if ok != true {
				return 500, []byte(""), fmt.Errorf("Basic auth")
			}

			return 200, []byte(u + ":" + p), nil
		},
	}

	params := &Params{
		Url:      "http://127.0.0.1/test",
		Username: "user",
		Password: "pass",
	}

	data, code, err := client.Get(params)
	if err != nil {
		t.Errorf("No, %s", err.Error())
		return
	}

	if code == 200 && bytes.Compare(data, payload) == 0 {
		t.Log("Yes.")
	} else {
		t.Errorf("No.")
	}
}

func TestStatusCode(t *testing.T) {
	t.Log("Is an error returned when a non-200 status code is received?")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(req *http.Request) (int, []byte, error) {
			return 404, []byte("404 Not Found"), nil
		},
	}

	params := &Params{
		Url: "http://127.0.0.1/test",
	}

	_, _, err := client.Get(params)
	if err != nil {
		if err.Error() == "APICLIENT: Received status code of 404 for http://127.0.0.1/test" {
			t.Log("Yes.")
		} else {
			t.Errorf("Error, but wrong type: %s", err.Error())
		}

		return
	}

	t.Errorf("No.")
}

func TestNoConnect(t *testing.T) {
	t.Log("Is an error returned when a non-200 status code is received?")

	conf := &Config{
		BytesPerSecond: 10000,
		MaxBytes:       10000,
		Timeout:        5,
	}

	client := NewClient(conf)
	client.Client = &FakeHttp{
		Payload: func(req *http.Request) (int, []byte, error) {
			return 0, []byte(""), fmt.Errorf("OMG IT BROKE!")
		},
	}

	params := &Params{
		Url: "http://127.0.0.1/test",
	}

	_, _, err := client.Get(params)
	if err != nil {
		if err.Error() == "APICLIENT: OMG IT BROKE!" {
			t.Log("Yes.")
		} else {
			t.Errorf("Error, but wrong type: %s", err.Error())
		}

		return
	}

	t.Errorf("No.")
}

/* apiclient_test.go ends here. */
