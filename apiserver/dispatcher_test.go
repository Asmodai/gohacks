/*
 * dispatcher_test.go --- Dispatcher tests.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/gin-gonic/gin"

	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"
)

type MockServer struct {
	LSTLSFn func(cert, key string) error
	LSFn    func() error
	SDFn    func(context.Context) error
}

func (ms *MockServer) ListenAndServeTLS(cert, key string) error {
	if ms.LSTLSFn == nil {
		return nil
	}

	return ms.LSTLSFn(cert, key)
}

func (ms *MockServer) ListenAndServe() error {
	if ms.LSFn == nil {
		return nil
	}

	return ms.LSFn()
}

func (ms *MockServer) Shutdown(ctx context.Context) error {
	if ms.SDFn == nil {
		return nil
	}

	return ms.SDFn(ctx)
}

func (ms *MockServer) SetTLSConfig(_ *tls.Config) {
}

func TestShit(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	srv := &MockServer{}

	lgr := logger.NewMockLogger("")
	lgr.Test = t

	inst := NewDispatcher(
		lgr,
		&Config{
			Addr:   ":8080",
			UseTLS: false,
		},
	)

	// Inject our fake HTTP server.
	inst.srv = srv

	t.Run("Starts without TLS", func(t *testing.T) {
		inst.config.UseTLS = false

		inst.Start()
		time.Sleep(1 * time.Second)
		inst.Stop()
	})

	t.Run("Errors when startup fails", func(t *testing.T) {
		inst.config.UseTLS = false
		srv.SDFn = nil
		srv.LSTLSFn = nil
		srv.LSFn = func() error {
			return fmt.Errorf("Synthetic error")
		}

		inst.Start()
		time.Sleep(2 * time.Second)
		inst.Stop()

		if lgr.LastFatal != "FATAL: listen() failed.  [err Synthetic error]" {
			t.Errorf("No, '%v'", lgr.LastFatal)
		}
	})

	t.Run("Errors with invalid TLS config", func(t *testing.T) {
		inst.config.UseTLS = true
		srv.SDFn = nil
		srv.LSFn = nil
		srv.LSTLSFn = func(_, _ string) error {
			return fmt.Errorf("Synthetic error")
		}

		inst.Start()
		time.Sleep(2 * time.Second)
		inst.Stop()

		if lgr.LastFatal != "FATAL: listen() failed.  [err Synthetic error]" {
			t.Errorf("No, '%v'", lgr.LastFatal)
		}
	})

	t.Run("Errors when shutdown fails", func(t *testing.T) {
		inst.config.UseTLS = false
		srv.LSFn = nil
		srv.LSTLSFn = nil
		srv.SDFn = func(_ context.Context) error {
			return fmt.Errorf("OH NO")
		}

		inst.Start()
		time.Sleep(1 * time.Second)
		inst.Stop()

		if lgr.LastFatal != "FATAL: API dispatcher server shutdown failure.  [err OH NO]" {
			t.Errorf("No, '%v'", lgr.LastFatal)
		}
	})
}

func TestLogWriter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	srv := &MockServer{}

	lgr := logger.NewMockLogger("")
	lgr.Test = t

	inst := NewDispatcher(
		lgr,
		&Config{
			Addr:   ":8080",
			UseTLS: false,
		},
	)

	// Inject our fake HTTP server.
	inst.srv = srv

	t.Run("Can open valid files", func(t *testing.T) {
		lgr.LastFatal = ""
		inst.config.LogFile = "/dev/null"

		f := inst.logWriter()
		time.Sleep(1 * time.Second)

		if lgr.LastFatal != "" {
			t.Errorf("No, '%v'", lgr.LastFatal)
		}

		if f == nil {
			t.Error("No, io.Writer was not returned")
		}

		if f != nil {
			f.(*os.File).Close()
		}
	})

	t.Run("Errors with invalid files", func(t *testing.T) {
		lgr.LastFatal = ""
		inst.config.LogFile = "/NOPE"

		f := inst.logWriter()
		time.Sleep(1 * time.Second)

		if lgr.LastFatal != "FATAL: Could not open file for writing.  [file /NOPE err open /NOPE: permission denied]" {
			t.Errorf("No, '%v'", lgr.LastFatal)
		}

		if f == nil {
			t.Error("No, io.Writer was not returned")
		}

		if f != nil {
			f.(*os.File).Close()
		}
	})
}

func TestFormatter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	srv := &MockServer{}

	lgr := logger.NewMockLogger("")
	lgr.Test = t

	inst := NewDispatcher(
		lgr,
		&Config{
			Addr:   ":8080",
			UseTLS: false,
		},
	)

	// Inject our fake HTTP server.
	inst.srv = srv

	t.Run("Outputs things", func(t *testing.T) {
		params := gin.LogFormatterParams{
			TimeStamp:  time.Time{},
			StatusCode: 200,
			Latency:    0,
			ClientIP:   "127.0.0.1",
			Path:       "/vmunix",
		}

		// *sigh*
		err := "0001/01/01 00:00:00 | 200 |            0s |       127.0.0.1 |          \"/vmunix\"\n"

		stuff := inst.logFormatter(params)

		if stuff != err {
			t.Errorf("No, '%v'", stuff)
		}
	})
}

/* dispatcher_test.go ends here. */
