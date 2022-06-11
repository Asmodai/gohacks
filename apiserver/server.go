/*
 * server.go --- Wrapper around http.Server.
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

package apiserver

import (
	"github.com/gin-gonic/gin"

	"context"
	"crypto/tls"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string, router *gin.Engine) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func NewDefaultServer() *Server {
	return &Server{
		srv: &http.Server{},
	}
}

func (s *Server) ListenAndServeTLS(cert, key string) error {
	return s.srv.ListenAndServeTLS(cert, key)
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) SetTLSConfig(conf *tls.Config) {
	s.srv.TLSConfig = conf
}

/* server.go ends here. */
