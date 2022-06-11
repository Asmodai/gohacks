/*
 * params.go --- API client parameters.
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

type AuthBasic struct {
	Username string
	Password string
}

type AuthToken struct {
	Header string
	Data   string
}

type ContentType struct {
	Accept string
	Type   string
}

// API client request parameters
type Params struct {
	Url string // API URL.

	UseBasic bool
	UseToken bool

	Content ContentType
	Token   AuthToken
	Basic   AuthBasic

	Queries []*QueryParam
}

// Add a new query parameter.
func (p *Params) AddQueryParam(name, content string) *QueryParam {
	q := NewQueryParam(name, content)

	p.Queries = append(p.Queries, q)

	return q
}

// Clear all query parameters.
func (p *Params) ClearQueryParams() {
	p.Queries = []*QueryParam{}
}

// Enable/disable basic authentication.
func (p *Params) SetUseBasic(val bool) {
	p.UseBasic = val
}

// Enable/disable authentication token.
func (p *Params) SetUseToken(val bool) {
	p.UseToken = val
}

// Create a new API parameters object,
func NewParams() *Params {
	return &Params{
		Url:      "",
		UseBasic: false,
		UseToken: false,
		Content: ContentType{
			Accept: "application/json",
			Type:   "application/json",
		},
		Token: AuthToken{
			Header: "",
			Data:   "",
		},
		Basic: AuthBasic{
			Username: "",
			Password: "",
		},
		Queries: []*QueryParam{},
	}
}

/* params.go ends here. */
