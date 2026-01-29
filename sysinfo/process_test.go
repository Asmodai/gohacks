// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// process_test.go --- SysInfo process tests.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

package sysinfo

// * Imports:

import (
	"context"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"
)

// * Code:

// ** Tests:

func TestProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mgr := process.NewManager()

	nctx, err := process.SetManager(ctx, mgr)
	if err != nil {
		t.Fatalf("Could not register process manager: %#v", err)
	}

	testSIProc, err := Spawn(nctx, types.Duration(1*time.Second))
	if err != nil {
		t.Fatalf("Spawn: %s", err.Error())
		return
	}

	time.Sleep(2 * time.Second)

	if !testSIProc.Running() {
		t.Fatal("Process is not running.")
		testSIProc.Stop()
		return
	}

	testSIProc.Stop()
}

// * process_test.go ends here.
