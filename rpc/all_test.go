/*
 * all_test.go --- RPC tests.
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

package rpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/logger"
)

var (
	TestInfo Info = Info{Name: "Test", Version: 12}
)

// ==================================================================
// {{{ Test RPC structures:

type TestArgs struct {
	X, Y int
}

type TestReply struct {
	X, Y, Result int
}

type TestRpc struct{}

func (t *TestRpc) RpcInfo(_ NoArgs, reply *Info) error {
	*reply = TestInfo

	return nil
}

func (t *TestRpc) Invoke(args TestArgs, reply *TestReply) error {
	reply.X = args.X
	reply.Y = args.Y
	reply.Result = args.X * args.Y

	return nil
}

// }}}
// ==================================================================

func TestAll(t *testing.T) {
	ctx := context.Background()
	lgr := logger.NewDefaultLogger()
	cnf := NewConfig("tcp4", "127.0.0.1:5432")

	var mgr *Manager = NewManager(cnf, ctx, lgr)
	var clt *Client = NewClient(cnf, lgr)

	mgr.Add(reflect.TypeOf(TestRpc{}))
	defer mgr.Shutdown()

	go func() {
		mgr.Start()
	}()
	time.Sleep(1 * time.Second)

	t.Run("Connects to RPC server", func(t *testing.T) {
		if err := clt.Dial(); err != nil {
			t.Error(err.Error())
			return
		}
	})
	defer clt.Close()

	var r *TestReply = &TestReply{}
	var i *Info = &Info{}

	args := TestArgs{3, 4}

	t.Run("Calls remote 'invoke' method", func(t *testing.T) {
		err := clt.Call("TestRpc.Invoke", args, &r)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})

	t.Run("Calls remote 'rpcinfo' method", func(t *testing.T) {
		err := clt.Call("TestRpc.RpcInfo", NoArgs{}, &i)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})

	t.Run("Reply is as expected", func(t *testing.T) {
		if r.X != args.X {
			t.Errorf("X is different:  %v != %v", r.X, args.X)
		}

		if r.Y != args.Y {
			t.Errorf("Y is different:  %v != %v", r.Y, args.Y)
		}

		if r.Result != 12 {
			t.Errorf("Result is wrong:  %v", r.Result)
		}
	})

	t.Run("Info is as expected", func(t *testing.T) {
		if i.Name != "Test" {
			t.Errorf("RPC name is wrong:  %v", i.Name)
		}

		if i.Version != 12 {
			t.Errorf("RPC version is wrong:  %v", i.Version)
		}
	})
}

/* all_test.go ends here. */
