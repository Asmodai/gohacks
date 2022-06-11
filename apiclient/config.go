/*
 * config.go --- API client configuration.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package apiclient

// API client configuration.
//
// `RequstsPerSecond` is the number of requests per second rate limiting.
// `Timeout` is obvious.
type Config struct {
	RequestsPerSecond int `json:"requests_per_second"`
	Timeout           int `json:"timeout"`
}

// Create a new API client configuration.
func NewConfig(ReqsPerSec, Timeout int) *Config {
	return &Config{
		RequestsPerSecond: ReqsPerSec,
		Timeout:           Timeout,
	}
}

// Return a new default API client configuration.
func NewDefaultConfig() *Config {
	return NewConfig(5, 5)
}

/* config.go ends here. */
