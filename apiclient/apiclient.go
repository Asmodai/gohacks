/*
 * apiclient.go --- Web API client.
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
	"fmt"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTP client interface.
//
// This is to allow mocking of `net/http`'s `http.Client` in unit tests.
// You could also, if drunk enough, provide your own HTTP client, as long as
// it conforms to the interface.  But you wouldn't want to do that, would you.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// API client request parameters.
type Params struct {
	Url      string // API endpoint URL.
	Token    string // API AIMS authentication token.
	Username string // Authentication username for basic auth.
	Password string // Authentication password for basic auth.
}

// Create a new API parameters object.
func NewParams() *Params {
	return &Params{}
}

/*

API Client

Like a finely-crafted sword, it can be wielded with skill if instructions
are followed.

1) Create your config:

	conf := &apiclient.Config{
		BytesPerSecond: 1024, // 1KiB/s
		MaxBytes:       8192, // 8KiB
		Timeout:        5,    // 5 seconds
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
	Client  HTTPClient
	Limiter *rate.Limiter
}

// API client configuration.
//
// `BytesPerSecond` is the transfer speed rate limit.
// `MaxBytes` is the maximum number of bytes limit.
// `Timeout` is obvious.
type Config struct {
	BytesPerSecond int `json:"bytes_per_second"`
	MaxBytes       int `json:"max_bytes"`
	Timeout        int `json:"timeout"`
}

// Return a new default API client configuration.
func NewDefaultConfig() *Config {
	return &Config{
		BytesPerSecond: 8192,
		MaxBytes:       999999,
		Timeout:        5,
	}
}

// Create a new API client with the given configuration.
func NewClient(config *Config) *Client {
	limiter := rate.NewLimiter((rate.Limit)(config.BytesPerSecond*1024), config.MaxBytes*1024)
	client := &http.Client{Timeout: (time.Duration)(config.Timeout) * time.Second}

	return &Client{
		Limiter: limiter,
		Client:  client,
	}
}

// The actual meat of the API client.
func (c *Client) httpAction(verb string, data *Params) ([]byte, int, error) {
	rv := c.Limiter.Reserve()
	if !rv.OK() {
		return nil, 0, fmt.Errorf("APICLIENT: Exceeded burst limit")
	}

	delay := rv.Delay()
	time.Sleep(delay)

	req, err := http.NewRequest(verb, data.Url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("APICLIENT: %s", err.Error())
	}

	req.Header.Add("Accept", "application/json")
	if data.Token != "" {
		req.Header.Add("X-Aims-Auth-Token", data.Token)
	}

	if data.Username != "" {
		req.SetBasicAuth(data.Username, data.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("APICLIENT: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf(
			"APICLIENT: Received status code of %d for %s",
			resp.StatusCode,
			req.URL,
		)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("APICLIENT: %s", err.Error())
	}

	return bytes, 200, nil
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

/* apiclient.go ends here. */
