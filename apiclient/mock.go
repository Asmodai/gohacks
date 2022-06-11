/*
 * mock.go --- Mock API client.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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
	"github.com/Asmodai/gohacks/logger"
)

type MockCallbackFn func(data *Params) ([]byte, int, error)

type MockClient struct {
	logger logger.ILogger
	GetFn  MockCallbackFn
	PostFn MockCallbackFn
}

func NewMockClient(config *Config, lgr logger.ILogger) *MockClient {
	return &MockClient{
		logger: lgr,
	}
}

func (c *MockClient) Get(data *Params) ([]byte, int, error) {
	if c.GetFn == nil {
		return []byte(""), 200, nil
	}

	return c.GetFn(data)
}

func (c *MockClient) Post(data *Params) ([]byte, int, error) {
	if c.PostFn == nil {
		return []byte(""), 200, nil
	}

	return c.PostFn(data)
}

/* mock.go ends here. */
