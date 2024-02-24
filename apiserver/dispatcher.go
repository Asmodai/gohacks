/*
 * dispatcher.go --- Route dispatch wrapper.
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

package apiserver

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/utils"

	"github.com/gin-gonic/gin"

	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type Dispatcher struct {
	config *Config
	srv    IServer
	router *gin.Engine
	lgr    logger.ILogger
}

func NewDispatcher(lgr logger.ILogger, config *Config) *Dispatcher {
	obj := &Dispatcher{
		config: config,
		router: gin.New(),
		lgr:    lgr,
	}

	ginlogger := utils.GetEnv("GIN_LOGGER", "file")
	switch ginlogger {
	case "file":
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

func NewDefaultDispatcher() *Dispatcher {
	return NewDispatcher(
		logger.NewDefaultLogger(),
		&Config{
			Addr:   ":8080",
			UseTLS: false,
		},
	)
}

func (d *Dispatcher) GetRouter() *gin.Engine {
	return d.router
}

func (d *Dispatcher) Start() {
	go func() {
		var err error

		switch d.config.UseTLS {
		case true:
			d.srv.SetTLSConfig(&tls.Config{
				MaxVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
			})
			err = d.srv.ListenAndServeTLS(d.config.Cert, d.config.Key)

		case false:
			err = d.srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			d.lgr.Fatal(
				"listen() failed.",
				"err", err.Error(),
			)
		}
	}()
}

func (d *Dispatcher) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (d *Dispatcher) notFound(c *gin.Context) {
	NewErrorDocument(
		404,
		fmt.Sprintf("%s was not found.", c.Request.RequestURI),
	).Write(c)
}

/* dispatcher.go ends here. */
