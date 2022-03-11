/*
 * client.go --- HTTP API client.
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
	"github.com/Asmodai/gohacks/rlhttp"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"

	"golang.org/x/time/rate"
)

/*

API Client

Like a finely-crafted sword, it can be wielded with skill if instructions
are followed.

1) Create your config:

	conf := &apiclient.Config{
    RequestsPerSecond: 5,    // 5 requests per second.
		Timeout:           5,    // 5 seconds.
	}

2) Create your client

	api := apiclient.NewClient(conf)

3) ???

	params := &Params{
		Url: "http://www.example.com/underpants",
	}

4) Profit

	data, code, err := api.Get(params)
	// check `err` and `code` here.
	// `data` will need to be converted from `[]byte`.

*/
type Client struct {
	Client  IHTTPClient
	Limiter *rate.Limiter
	Trace   *httptrace.ClientTrace
}

// Create a new API client with the given configuration.
func NewClient(config *Config) *Client {
	trace := &httptrace.ClientTrace{}

	limiter := rate.NewLimiter(
		rate.Every(1*time.Second),
		config.RequestsPerSecond,
	)

	client := rlhttp.NewClient(
		limiter,
		(time.Duration)(config.Timeout)*time.Second,
	)

	return &Client{
		Limiter: limiter,
		Client:  client,
		Trace:   trace,
	}
}

// Perform a HTTP GET using the given API parameters.
//
// Returns the response body as an array of bytes, the HTTP status code, and
// an error if one is triggered.
//
// You will need to remember to check the error *and* the status code.
func (c *Client) Get(data *Params) ([]byte, int, error) {
	return c.httpAction("GET", data)
}

// Perform a HTTP POST using the given API parameters.
//
// Returns the response body as an array of bytes, the HTTP status code, and
// an error if one is triggered.
//
// You will need to remember to check the error *and* the status code.
func (c *Client) Post(data *Params) ([]byte, int, error) {
	return c.httpAction("POST", data)
}

// The actual meat of the API client.
func (c *Client) httpAction(verb string, data *Params) ([]byte, int, error) {
	req, err := http.NewRequest(verb, data.Url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("APICLIENT: %s", err.Error())
	}

	// Set `Accept` header if required.
	if data.Content.Accept != "" {
		req.Header.Add("Accept", data.Content.Accept)
	}

	// Set `Content-Type` header if required.
	if data.Content.Type != "" {
		req.Header.Add("Content-Type", data.Content.Type)
	}

	// Error if we're told to use both basic auth and an auth token.
	if data.UseBasic && data.UseToken {
		return nil, 0, fmt.Errorf(
			"APICLIENT: Cannot use Basic Auth and token at the same time.",
		)
	}

	// Set up basic authentication if required.
	if data.UseBasic {
		if data.Basic.Username == "" {
			return nil, 0, fmt.Errorf(
				"APICLIENT: No basic auth username given!",
			)
		}

		req.SetBasicAuth(data.Basic.Username, data.Basic.Password)
	}

	// Set up authentication token if required.
	if data.UseToken {
		if data.Token.Header == "" {
			return nil, 0, fmt.Errorf(
				"APICLIENT: No auth token header given!",
			)
		}

		req.Header.Add(data.Token.Header, data.Token.Data)
	}

	// Append URI parameters.
	if len(data.Queries) != 0 {
		q := req.URL.Query()

		for idx, _ := range data.Queries {
			q.Add(data.Queries[idx].Name, data.Queries[idx].Content)
		}

		req.URL.RawQuery = q.Encode()
	}

	// Perform the request.
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("APICLIENT: %s", err.Error())
	}

	// Did we get a non-200 status code?
	if resp.StatusCode != 200 {
		resp.Body.Close()

		return nil, resp.StatusCode, fmt.Errorf(
			"APICLIENT: Received status code %d for %s",
			resp.StatusCode,
			req.URL,
		)
	}

	// Decode the body.
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()

		return nil, resp.StatusCode, fmt.Errorf("APICLIENT: %s", err.Error())
	}

	resp.Body.Close()

	return bytes, resp.StatusCode, nil
}

/* client.go ends here. */
