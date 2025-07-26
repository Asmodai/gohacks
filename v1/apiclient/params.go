// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// params.go --- API client parameters.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:

//
//
//

// * Package:

package apiclient

// * Code:

// ** Types:

// Basic authentication configuration.
type AuthBasic struct {
	Username string
	Password string
}

// Authentication token configuration.
//
// This type does not care about the token type.  If you intend to make use of
// token types such as JWTs then you must implement that code yourself.
type AuthToken struct {
	// The HTTP header added to the request that contains the token data.
	Header string

	// The authentication token.
	Data string
}

// Content types configuration.
type ContentType struct {
	// The MIME type that we wish to accept.  Sent in the request via the `Accept`
	// header.
	Accept string

	// The MIME type of the data we are sending.  Sent in the request via the
	// `Content-Type` header.
	Type string
}

// API client request parameters.
//
// Although we support both Basic Auth and authentication tokens, only one
// should be used.  If both are specified, an error will be triggered.
type Params struct {
	// Request URL.
	URL string

	// Use the Basic Auth schema to authenticate to the remote server.
	UseBasic bool

	// Use an 'authentication token' to authenticate to the remote server.
	UseToken bool

	// MIME content type of the data we are sending and wish to accept from the
	// remote server.
	Content ContentType

	// Authentication token configuration.
	Token AuthToken

	// Basic Auth configuration.
	Basic AuthBasic

	// Queries that are sent via the HTTP request.
	Queries []*QueryParam
}

// ** Methods:

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

// ** Functions:

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

// * params.go ends here.
