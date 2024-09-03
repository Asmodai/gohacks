// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// dispatcher_test.go --- Dispatcher tests.
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
	mlogger "github.com/Asmodai/gohacks/mocks/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type FakeServer struct {
	LSTLSFn func(cert, key string) error
	LSFn    func() error
	SDFn    func(context.Context) error
}

func (ms *FakeServer) ListenAndServeTLS(cert, key string) error {
	if ms.LSTLSFn == nil {
		return nil
	}

	return ms.LSTLSFn(cert, key)
}

func (ms *FakeServer) ListenAndServe() error {
	if ms.LSFn == nil {
		return nil
	}

	return ms.LSFn()
}

func (ms *FakeServer) Shutdown(ctx context.Context) error {
	if ms.SDFn == nil {
		return nil
	}

	return ms.SDFn(ctx)
}

func (ms *FakeServer) SetTLSConfig(_ *tls.Config) {
}

func TestDispatch(t *testing.T) {
	mocked := gomock.NewController(t)
	defer mocked.Finish()

	gin.SetMode(gin.ReleaseMode)

	srv := &FakeServer{}
	lgr := mlogger.NewMockLogger(mocked)
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

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
		errmsg := "Forced startup failure"

		inst.config.UseTLS = false
		srv.SDFn = nil
		srv.LSTLSFn = nil
		srv.LSFn = func() error {
			return fmt.Errorf(errmsg)
		}

		lgr.EXPECT().
			Fatal("listen() failed.", gomock.Any()).
			Do(func(msg string, rest ...any) {
				got := fmt.Sprintf("%s %v", msg, rest)
				want := fmt.Sprintf("listen() failed. [err %s]", errmsg)

				if got != want {
					t.Errorf("Unexpected error: %s", got)
				}
			})

		inst.Start()
		time.Sleep(2 * time.Second)
		inst.Stop()
	})

	t.Run("Errors with invalid TLS config", func(t *testing.T) {
		errmsg := "Forced TLS failure"

		inst.config.UseTLS = true
		srv.SDFn = nil
		srv.LSFn = nil
		srv.LSTLSFn = func(_, _ string) error {
			return fmt.Errorf(errmsg)
		}

		lgr.EXPECT().
			Fatal(gomock.Any(), gomock.Any()).
			Do(func(msg string, rest ...any) {
				got := fmt.Sprintf("%s %v", msg, rest)
				want := fmt.Sprintf("listen() failed. [err %s]", errmsg)

				if got != want {
					t.Errorf("Unexpected error: %s", got)
				}
			})

		inst.Start()
		time.Sleep(2 * time.Second)
		inst.Stop()
	})

	t.Run("Errors when shutdown fails", func(t *testing.T) {
		errmsg := "Forced shutdown failure"

		inst.config.UseTLS = false
		srv.LSFn = nil
		srv.LSTLSFn = nil
		srv.SDFn = func(_ context.Context) error {
			return fmt.Errorf(errmsg)
		}

		lgr.EXPECT().
			Fatal(gomock.Any(), gomock.Any()).
			Do(func(msg string, rest ...any) {
				got := fmt.Sprintf("%s %v", msg, rest)
				want := fmt.Sprintf("API dispatcher server shutdown failure. [err %s]", errmsg)

				if got != want {
					t.Errorf("Unexpected error: %s", got)
				}
			})

		inst.Start()
		time.Sleep(1 * time.Second)
		inst.Stop()
	})
}

func TestLogWriter(t *testing.T) {
	mocked := gomock.NewController(t)
	defer mocked.Finish()

	gin.SetMode(gin.ReleaseMode)

	srv := &FakeServer{}
	lgr := mlogger.NewMockLogger(mocked)
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

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
		inst.config.LogFile = "/dev/null"

		f := inst.logWriter()
		time.Sleep(1 * time.Second)

		if f == nil {
			t.Error("No, io.Writer was not returned")
		}

		if f != nil {
			f.(*os.File).Close()
		}
	})

	t.Run("Errors with invalid files", func(t *testing.T) {
		lgr.EXPECT().
			Fatal(gomock.Any(), gomock.Any()).
			Do(func(msg string, rest ...any) {
				got := fmt.Sprintf("%s %v", msg, rest)

				switch got {
				case "Could not open file for writing. [file /NOPE err open /NOPE: permission denied]":
					// The root filesystem is not read-only (macOS), but cannot open /NOPE for writing.
					break

				case "Could not open file for writing. [file /NOPE err open /NOPE: read-only file system]":
					// The root file system is read-only (macOS)
					break

				default:
					t.Errorf("Unexpected error: %s", got)
				}
			}).
			AnyTimes()

		inst.config.LogFile = "/NOPE"

		f := inst.logWriter()
		time.Sleep(1 * time.Second)

		if f == nil {
			t.Error("No, io.Writer was not returned")
		}

		if f != nil {
			f.(*os.File).Close()
		}
	})
}

func TestFormatter(t *testing.T) {
	mocked := gomock.NewController(t)
	defer mocked.Finish()

	gin.SetMode(gin.ReleaseMode)

	srv := &FakeServer{}
	lgr := mlogger.NewMockLogger(mocked)
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

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
		req, _ := http.NewRequest("GET", "/vmunix", nil)
		req.Header.Add("User-Agent", "derp")

		params := gin.LogFormatterParams{
			Method:     "GET",
			TimeStamp:  time.Time{},
			StatusCode: 200,
			Latency:    0,
			ClientIP:   "127.0.0.1",
			Path:       "/vmunix",
			Request:    req,
		}

		// *sigh* keep the space at the end for the error message component.
		err := "127.0.0.1 - [Mon, 01 Jan 0001 00:00:00 UTC] \"GET /vmunix HTTP/1.1\" 200 0s \"derp\" "

		stuff := strings.TrimSuffix(inst.logFormatter(params), "\n")

		if stuff != err {
			t.Errorf("No, '%v'", stuff)
		}
	})
}

// dispatcher_test.go ends here.
