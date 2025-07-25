// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// dispatcher.go --- Route dispatch wrapper.
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
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/utils"

	"github.com/gin-gonic/gin"

	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

const (
	MinimumTimeout = 5
)

// API route dispatcher.
type Dispatcher struct {
	config *Config
	srv    Server
	router *gin.Engine
	lgr    logger.Logger
}

// Create a new API route dispatcher.
func NewDispatcher(lgr logger.Logger, config *Config) *Dispatcher {
	obj := &Dispatcher{
		config: config,
		router: gin.New(),
		lgr:    lgr,
	}

	ginlogger := utils.GetEnv("GIN_LOGGER", "file")
	if ginlogger == "file" {
		obj.router.Use(gin.LoggerWithConfig(obj.initLog()))
	}

	obj.router.NoRoute(obj.notFound)
	obj.router.Use(gin.Recovery())
	obj.router.Use(CORSMiddleware())
	obj.router.Use(func(c *gin.Context) {
		c.Set("start_time", time.Now())
		c.Next()
	})

	obj.srv = NewServer(obj.config.Addr, obj.router)

	lgr.Info(
		"Dispatcher initialised.",
		"addr", config.Addr,
	)

	return obj
}

// Create a new API route dispatcher with default values.
//
// The dispatcher returned by this function will listen on port 8080 and
// bind to all available addresses on the host machine.
func NewDefaultDispatcher() *Dispatcher {
	return NewDispatcher(
		logger.NewDefaultLogger(),
		&Config{
			Addr:   ":8080",
			UseTLS: false,
		},
	)
}

// Return the router used by this dispatcher.
func (d *Dispatcher) GetRouter() *gin.Engine {
	return d.router
}

// Start the API route dispatcher.
func (d *Dispatcher) Start() {
	go func() {
		var err error

		switch d.config.UseTLS {
		case true:
			// Configure TLS with a minimum version.
			// Do not set a max version unless there is a specific reason.
			// The minimum version should NEVER be below v1.2.
			d.srv.SetTLSConfig(&tls.Config{
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
			})

			err = d.srv.ListenAndServeTLS(d.config.Cert, d.config.Key)

		case false:
			err = d.srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			d.lgr.Fatal(
				"listen() failed.",
				"err", err.Error(),
			)
		}
	}()
}

// Stop the API route dispatcher.
func (d *Dispatcher) Stop() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		MinimumTimeout*time.Second,
	)

	defer func() {
		cancel()
	}()

	if err := d.srv.Shutdown(ctx); err != nil {
		d.lgr.Fatal(
			"API dispatcher server shutdown failure.",
			"err", err.Error(),
		)
	}
}

// Default function invoked when the API dispatcher does not have a handler
// available for a route.
func (d *Dispatcher) notFound(c *gin.Context) {
	NewErrorDocument(
		http.StatusNotFound,
		c.Request.RequestURI+" was not found.",
	).Write(c)
}

// dispatcher.go ends here.
