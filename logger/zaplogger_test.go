// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// zaplogger_test.go --- Zap logger testing.
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
//

// * Package:

package logger

// * Imports:

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// * Code:

// ** Tests:

func TestLoggerShim_Info(t *testing.T) {
	var buf bytes.Buffer

	// Create an encoder config that gives us something readable
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "" // Disable time for easier asserts

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)

	logger := zap.New(core)
	shim := &zapLogger{logger: logger}

	// Act
	shim.Info("Scaling down workers.",
		"pool", "alpha",
		"current", 8,
		"required", 4,
		"new", 4,
	)

	// Now check what got written
	out := buf.String()
	if !strings.Contains(out, `"msg":"Scaling down workers."`) {
		t.Errorf("expected message not found in log output: %s", out)
	}

	if !strings.Contains(out, `"pool":"alpha"`) {
		t.Errorf("expected field not found in log output: %s", out)
	}

	// Optional: decode JSON and assert fields more surgically
	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}

	if parsed["new"] != float64(4) { // JSON numbers decode as float64
		t.Errorf("unexpected value for 'new': %v", parsed["new"])
	}
}

// * zaplogger_test.go ends here.
