/*
 * rlhttp.go --- Rate-limited HTTP client.
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

package rlhttp

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RLHTTPClient struct {
	client  *http.Client
	limiter *rate.Limiter
}

func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if err := c.limiter.Wait(context.Background()); err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewClient(rl *rate.Limiter, timeout time.Duration) *RLHTTPClient {
	client := &http.Client{
		Timeout: timeout,
	}

	return &RLHTTPClient{
		client:  client,
		limiter: rl,
	}
}

/* rlhttp.go ends here. */
