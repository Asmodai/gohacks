/* mock:yes */
/*
 * server.go --- Wrapper around http.Server.
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
	"github.com/gin-gonic/gin"
	"gitlab.com/tozd/go/errors"

	"context"
	"crypto/tls"
	"net/http"
	"time"
)

const (
	defaultReadHeaderTimeout = 2 * time.Second
)

type Server interface {
	ListenAndServeTLS(string, string) error
	ListenAndServe() error
	Shutdown(context.Context) error
	SetTLSConfig(*tls.Config)
}

type server struct {
	srv *http.Server
}

func NewServer(addr string, router *gin.Engine) Server {
	return &server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
	}
}

func NewDefaultServer() Server {
	return &server{
		srv: &http.Server{
			// Avoid a Slowlaris issue and set a read header timeout.
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
	}
}

func (s *server) ListenAndServeTLS(cert, key string) error {
	if err := s.srv.ListenAndServeTLS(cert, key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *server) ListenAndServe() error {
	if err := s.srv.ListenAndServe(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *server) SetTLSConfig(conf *tls.Config) {
	s.srv.TLSConfig = conf
}

/* server.go ends here. */
