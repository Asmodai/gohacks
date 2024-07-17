/* mock:yes */
/*
 * client.go --- HTTP API client.
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

package apiclient

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/rlhttp"

	"gitlab.com/tozd/go/errors"
	"golang.org/x/time/rate"

	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

var (
	ErrInvalidAuthMethod = errors.Base("invalid authentication method")
	ErrMissingArgument   = errors.Base("missing argument")
	ErrNotOk             = errors.Base("not ok")
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
		URL: "http://www.example.com/underpants",
	}

4) Profit

	data, code, err := api.Get(params)
	// check `err` and `code` here.
	// `data` will need to be converted from `[]byte`.
*/
type Client interface {
	Get(*Params) ([]byte, int, error)
	Post(*Params) ([]byte, int, error)
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	Client  HTTPClient
	Limiter *rate.Limiter
	Trace   *httptrace.ClientTrace
	logger  logger.Logger
}

// Create a new API client with the given configuration.
func NewClient(config *Config, logger logger.Logger) Client {
	trace := &httptrace.ClientTrace{}

	limiter := rate.NewLimiter(
		rate.Every(1*time.Second),
		config.RequestsPerSecond,
	)

	//nolint:durationcheck
	rlclient := rlhttp.NewClient(
		limiter,
		(time.Duration)(config.Timeout)*time.Second,
	)

	return &client{
		Limiter: limiter,
		Client:  rlclient,
		Trace:   trace,
		logger:  logger,
	}
}

// Perform a HTTP GET using the given API parameters.
//
// Returns the response body as an array of bytes, the HTTP status code, and
// an error if one is triggered.
//
// You will need to remember to check the error *and* the status code.
func (c *client) Get(data *Params) ([]byte, int, error) {
	return c.httpAction("GET", data)
}

// Perform a HTTP POST using the given API parameters.
//
// Returns the response body as an array of bytes, the HTTP status code, and
// an error if one is triggered.
//
// You will need to remember to check the error *and* the status code.
func (c *client) Post(data *Params) ([]byte, int, error) {
	return c.httpAction("POST", data)
}

// The actual meat of the API client.
// TODO:Rewrite using contexts.
//
//nolint:cyclop,funlen,noctx
func (c *client) httpAction(verb string, data *Params) ([]byte, int, error) {
	req, err := http.NewRequest(verb, data.URL, nil)
	if err != nil {
		return nil, 0, errors.WithStack(err)
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
		return nil, 0, errors.Wrap(
			ErrInvalidAuthMethod,
			"cannot use basic auth and token at the same time",
		)
	}

	// Set up basic authentication if required.
	if data.UseBasic {
		if data.Basic.Username == "" {
			return nil, 0, errors.Wrap(
				ErrMissingArgument,
				"no basic auth username given",
			)
		}

		req.SetBasicAuth(data.Basic.Username, data.Basic.Password)
	}

	// Set up authentication token if required.
	if data.UseToken {
		if data.Token.Header == "" {
			return nil, 0, errors.Wrap(
				ErrMissingArgument,
				"no auth token header given",
			)
		}

		req.Header.Add(data.Token.Header, data.Token.Data)
	}

	// Append URI parameters.
	if len(data.Queries) != 0 {
		qry := req.URL.Query()

		for idx := range data.Queries {
			qry.Add(data.Queries[idx].Name, data.Queries[idx].Content)
		}

		req.URL.RawQuery = qry.Encode()
	}

	// Perform the request.
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	// Did we get a non-200 status code?
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()

		return nil, resp.StatusCode, errors.Wrap(
			ErrNotOk,
			fmt.Sprintf(
				"received status code %d for %s",
				resp.StatusCode,
				req.URL,
			),
		)
	}

	// Decode the body.
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()

		return nil, resp.StatusCode, errors.WithStack(err)
	}

	resp.Body.Close()

	return bytes, resp.StatusCode, nil
}

/* client.go ends here. */
