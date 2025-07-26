// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// connection.go --- AMQP connection mockery.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package amqpmock

// * Imports:

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/Asmodai/gohacks/v1/amqp/amqpshim"
	mamqpshim "github.com/Asmodai/gohacks/v1/mocks/amqpshim"
	goamqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/mock/gomock"
)

// * Code:

// ** Mock result types:

// ** Types

type MockConnection struct {
	channel         ChannelFn
	close           SimpleErrorFn
	closeDeadline   CloseDeadlineFn
	connectionState ConnectionStateFn
	isClosed        SimpleBoolFn
	localAddr       SimpleAddrFn
	notifyBlocked   NotifyBlockedFn
	notifyClose     NotifyCloseFn
	remoteAddr      SimpleAddrFn
	updateSecret    UpdateSecretFn

	CallLog CallLog
}

// ** Methods:

// *** Initialisation:

func (obj *MockConnection) Init() {
	obj.BuildChannel(ChannelResults{})
	obj.BuildClose(ErrorResults{})
	obj.BuildCloseDeadline(ErrorResults{})
	obj.BuildConnectionState(ConnectionStateResults{})
	obj.BuildIsClosed(BoolResults{})
	obj.BuildLocalAddr(AddrResults{})
	obj.BuildNotifyBlocked(NotifyBlockedResults{})
	obj.BuildNotifyClose(NotifyCloseResults{})
	obj.BuildRemoteAddr(AddrResults{})
	obj.BuildUpdateSecret(ErrorResults{})
}

func (obj *MockConnection) AddCallLog(fname string, values ...any) {
	if obj.CallLog == nil {
		obj.CallLog = make(CallLog)
	}

	obj.CallLog[fname] = append(obj.CallLog[fname], CallLogList{values})
}

// *** `Channel`:

func (obj *MockConnection) SetChannel(fn ChannelFn) {
	obj.channel = fn
}

func (obj *MockConnection) BuildChannel(results ChannelResults) {
	obj.SetChannel(func() (amqpshim.Channel, error) {
		return results.Channel, results.Error
	})
}

func (obj *MockConnection) MockChannel(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		Channel().
		DoAndReturn(func() (amqpshim.Channel, error) {
			obj.AddCallLog("Channel")

			return obj.channel()
		})
}

// *** `Close`:

func (obj *MockConnection) SetClose(fn SimpleErrorFn) {
	obj.close = fn
}

func (obj *MockConnection) BuildClose(results ErrorResults) {
	obj.SetClose(func() error {
		return results.Error
	})
}

func (obj *MockConnection) MockClose(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		Close().
		DoAndReturn(func() error {
			obj.AddCallLog("Close")

			return obj.close()
		})
}

// *** `CloseDeadline`:

func (obj *MockConnection) SetCloseDeadline(fn CloseDeadlineFn) {
	obj.closeDeadline = fn
}

func (obj *MockConnection) BuildCloseDeadline(results ErrorResults) {
	obj.SetCloseDeadline(func(_ time.Time) error {
		return results.Error
	})
}

func (obj *MockConnection) MockCloseDeadline(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		CloseDeadline(gomock.Any()).
		DoAndReturn(func(deadline time.Time) error {
			obj.AddCallLog("CloseDeadline", deadline)

			return obj.closeDeadline(deadline)
		})
}

// *** `ConnectionState`:

func (obj *MockConnection) SetConnectionState(fn ConnectionStateFn) {
	obj.connectionState = fn
}

func (obj *MockConnection) BuildConnectionState(result ConnectionStateResults) {
	obj.SetConnectionState(func() tls.ConnectionState {
		return result.TLS
	})
}

func (obj *MockConnection) MockConnectionState(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		ConnectionState().
		DoAndReturn(func() tls.ConnectionState {
			obj.AddCallLog("ConnectionState")

			return obj.connectionState()
		})
}

// *** `IsClosed`:

func (obj *MockConnection) SetIsClosed(fn SimpleBoolFn) {
	obj.isClosed = fn
}

func (obj *MockConnection) BuildIsClosed(result BoolResults) {
	obj.SetIsClosed(func() bool {
		return result.Value
	})
}

func (obj *MockConnection) MockIsClosed(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		IsClosed().
		DoAndReturn(func() bool {
			obj.AddCallLog("IsClosed")

			return obj.isClosed()
		})
}

// *** `LocalAddr`:

func (obj *MockConnection) SetLocalAddr(fn SimpleAddrFn) {
	obj.localAddr = fn
}

func (obj *MockConnection) BuildLocalAddr(results AddrResults) {
	obj.SetLocalAddr(func() net.Addr {
		return results.Addr
	})
}

func (obj *MockConnection) MockLocalAddr(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		LocalAddr().
		DoAndReturn(func() net.Addr {
			obj.AddCallLog("LocalAddr")

			return obj.localAddr()
		})
}

// *** `NotifyBlocked`:

func (obj *MockConnection) SetNotifyBlocked(fn NotifyBlockedFn) {
	obj.notifyBlocked = fn
}

func (obj *MockConnection) BuildNotifyBlocked(results NotifyBlockedResults) {
	obj.SetNotifyBlocked(func(_ chan goamqp.Blocking) chan goamqp.Blocking {
		return results.BlockingChan
	})
}

func (obj *MockConnection) MockNotifyBlocked(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		NotifyBlocked(gomock.Any()).
		DoAndReturn(func(c chan goamqp.Blocking) chan goamqp.Blocking {
			obj.AddCallLog("NotifyBlocked", c)

			return obj.notifyBlocked(c)
		})
}

// *** `NotifyClose`:

func (obj *MockConnection) SetNotifyClose(fn NotifyCloseFn) {
	obj.notifyClose = fn
}

func (obj *MockConnection) BuildNotifyClose(results NotifyCloseResults) {
	obj.SetNotifyClose(func(_ chan *goamqp.Error) chan *goamqp.Error {
		return results.ErrorChan
	})
}

func (obj *MockConnection) MockNotifyClose(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		NotifyClose(gomock.Any()).
		DoAndReturn(func(c chan *goamqp.Error) chan *goamqp.Error {
			obj.AddCallLog("NotifyClose", c)

			return obj.notifyClose(c)
		})
}

// *** `RemoteAddr`:

func (obj *MockConnection) SetRemoteAddr(fn SimpleAddrFn) {
	obj.localAddr = fn
}

func (obj *MockConnection) BuildRemoteAddr(results AddrResults) {
	obj.SetRemoteAddr(func() net.Addr {
		return results.Addr
	})
}

func (obj *MockConnection) MockRemoteAddr(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		RemoteAddr().
		DoAndReturn(func() net.Addr {
			obj.AddCallLog("RemoteAddr")

			return obj.remoteAddr()
		})
}

// *** `UpdateSecret`:

func (obj *MockConnection) SetUpdateSecret(fn UpdateSecretFn) {
	obj.updateSecret = fn
}

func (obj *MockConnection) BuildUpdateSecret(results ErrorResults) {
	obj.SetUpdateSecret(func(_, _ string) error {
		return results.Error
	})
}

func (obj *MockConnection) MockUpdateSecret(mock *mamqpshim.MockConnection) *gomock.Call {
	return mock.EXPECT().
		UpdateSecret(gomock.Any(), gomock.Any()).
		DoAndReturn(func(newSecret, reason string) error {
			obj.AddCallLog("UpdateSecret", newSecret, reason)

			return obj.updateSecret(newSecret, reason)
		})
}

// * connection.go ends here.
