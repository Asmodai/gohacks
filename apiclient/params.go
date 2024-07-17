/*
 * params.go --- API client parameters.
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

// API client request parameters.
type Params struct {
	URL string // API URL.

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

// Create a new API parameters object.
func NewParams() *Params {
	return &Params{
		URL:      "",
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
