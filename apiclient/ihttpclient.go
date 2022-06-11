/*
 * ihttpclient.go --- HTTP client interface.
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

import (
	"net/http"
)

// HTTP client interface.
//
// This is to allow mocking of `net/http`'s `http.Client` in unit tests.
// You could also, if drunk enough, provide your own HTTP client, as long as
// it conforms to the interface.  But you wouldn't want to do that, would you.
type IHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

/* ihttpclient.go ends here. */
