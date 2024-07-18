// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// cors.go --- CORS middleware.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package apiserver

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
)

//nolint:gochecknoglobals
var (
	// List of allowed HTTP headers.
	AllowedHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Accept",
		"Origin",
		"Cache-Control",
		"X-Requested-With",
	}

	// List of allowed HTTP methods.
	AllowedMethods = []string{
		"POST",
		"OPTIONS",
		"GET",
		"PUT",
		"DELETE",
		"PATCH",
	}
)

// A gin-gonic handler for handling CORS.
func CORSMiddleware() gin.HandlerFunc {
	hdrs := strings.Join(AllowedHeaders, ",")
	mths := strings.Join(AllowedMethods, ",")

	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(
			"Access-Control-Allow-Origin", "*",
		)

		ctx.Writer.Header().Set(
			"Access-Control-Allow-Credentials", "true",
		)

		ctx.Writer.Header().Set(
			"Access-Control-Allow-Headers", hdrs,
		)

		ctx.Writer.Header().Set(
			"Access-Control-Allow-Methods", mths,
		)

		ctx.Writer.Header().Set(
			"Access-Control-Max-Age", "600",
		)

		// If we're using the OPTIONS method, then return no content.
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)

			return
		}

		// Invoke next handler.
		ctx.Next()
	}
}

// cors.go ends here.
