// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config_test.go --- Config tests.
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

package config

import (
	"github.com/Asmodai/gohacks/semver"

	"fmt"
	"os"
	"testing"
)

var DummyConfigJson = "{\"testing\":{\"opt1\":\"testing\",\"opt2\":42}}"

var version = &semver.SemVer{
	Major:  0,
	Minor:  1,
	Patch:  0,
	Commit: "test",
}

var StringRep = `Configuration:
    Opt1: [testing]
    Opt2: [42]
    Flags:
        BoolFlag: [false]
        F64Flag: [0]
        IntFlag: [0]
        I64Flag: [0]
        StringFlag: []
        UIntFlag: [0]
        UI64Flag: [0]`

type DummyConfig struct {
	Opt1 string `json:"opt1"`
	Opt2 int    `json:"opt2" config_validator:"ValidOpt2"`

	Flags struct {
		BoolFlag   bool
		F64Flag    float64
		IntFlag    int
		I64Flag    int64
		StringFlag string
		UIntFlag   uint
		UI64Flag   uint64
	}
}

type AccessorConfig struct {
}

func ValidOpt2(value int) error {
	if value != 42 {
		return fmt.Errorf("Not 42")
	}

	return nil
}

func MakeFns() map[string]interface{} {
	return map[string]interface{}{
		"ValidOpt2": ValidOpt2,
	}
}

func BenchmarkConfig(b *testing.B) {
	var str string

	for i := 0; i < 100; i++ {
		cnf := Init("Test", version, &DummyConfig{}, nil, false)
		str = cnf.String()
	}
	fmt.Printf("Config: %s\n", str)
}

func TestSimple(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Errorf("getwd: %s", err.Error())
		return
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"wibble", "-config", path + "/../testing/conf.json"}

	opts := &DummyConfig{}
	conf := Init("Test", version, opts, MakeFns(), false)

	t.Run("Construct config object", func(t *testing.T) {
		if conf == nil {
			t.Error("Could not build config object!")
		}
	})
	conf.Parse()

	t.Run("Check debug is false", func(t *testing.T) {
		if conf.IsDebug() != false {
			t.Error("Debug is true!")
		}
	})

	t.Run("String representation is as expected", func(t *testing.T) {
		if conf.String() != StringRep {
			t.Errorf("No, got:\n%v\bWanted:\n%v\n", conf.String(), StringRep)
		}
	})
}

func TestValidators(t *testing.T) {
	t.Run("Add validators on the fly", func(t *testing.T) {
		c := NewDefaultConfig(false)
		f := func() {}
		c.AddValidator("test", f)

		if c.(*config).Validators["test"] == nil {
			t.Error("Could not add validators at runtime.")
		}
	})
}

func TestCLIFlags(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Errorf("getwd: %s", err.Error())
		return
	}

	oldArgs := os.Args
	os.Args = []string{
		"wibble",
		"-config", path + "/../testing/conf.json",
		"-bool=true",
		"-float64=1.23",
		"-int=-2",
		"-int64=3456",
		"-string=seven",
		"-uint=8",
		"-uint64=90210",
	}

	o := &DummyConfig{}
	c := Init("Test", version, o, MakeFns(), false)

	c.AddBoolFlag(&o.Flags.BoolFlag, "bool", false, "bool")
	c.AddFloat64Flag(&o.Flags.F64Flag, "float64", 0.0, "float64")
	c.AddIntFlag(&o.Flags.IntFlag, "int", 0, "int")
	c.AddInt64Flag(&o.Flags.I64Flag, "int64", 0, "int64")
	c.AddStringFlag(&o.Flags.StringFlag, "string", "", "string")
	c.AddUintFlag(&o.Flags.UIntFlag, "uint", 0, "uint")
	c.AddUint64Flag(&o.Flags.UI64Flag, "uint64", 0, "uint64")

	c.Parse()
	os.Args = oldArgs

	t.Run("Look up uint value", func(t *testing.T) {
		flag := c.LookupFlag("uint")
		if flag == nil {
			t.Error("`LookupFlag` did not work.")
		}
	})

	t.Run("Look up bool value", func(t *testing.T) {
		if o.Flags.BoolFlag != true {
			t.Errorf("Unexpected boolean value %#v", o.Flags.BoolFlag)
		}
	})

	t.Run("Look up 64-bit float value", func(t *testing.T) {
		if o.Flags.F64Flag != 1.23 {
			t.Errorf("Unexpected float64 value %#v", o.Flags.F64Flag)
		}
	})

	t.Run("Look up 32-bit integer value", func(t *testing.T) {
		if o.Flags.IntFlag != -2 {
			t.Errorf("Unexpected int value %#v", o.Flags.IntFlag)
		}
	})

	t.Run("Look up 64-bit integer value", func(t *testing.T) {
		if o.Flags.I64Flag != 3456 {
			t.Errorf("Unexpected int64 value %#v", o.Flags.I64Flag)
		}
	})

	t.Run("Look up string value", func(t *testing.T) {
		if o.Flags.StringFlag != "seven" {
			t.Errorf("Unexpected string value %#v", o.Flags.StringFlag)
		}
	})

	t.Run("Look up 32-bit unsigned value", func(t *testing.T) {
		if o.Flags.UIntFlag != 8 {
			t.Errorf("Unexpected uint32 value %#v", o.Flags.UIntFlag)
		}
	})

	t.Run("Look up 64-bit unsigned value", func(t *testing.T) {
		if o.Flags.UI64Flag != 90210 {
			t.Errorf("Unexpected uint64 value %#v", o.Flags.UI64Flag)
		}
	})
}

func TestAccessors(t *testing.T) {
	var name string = "TestApp"
	var vers *semver.SemVer
	var path string
	var file string = "/../testing/conf.json"
	var log string = "/../logs/test.log"
	var err error

	path, err = os.Getwd()
	if err != nil {
		t.Errorf("getwd: %s", err.Error())
	}

	vers, err = semver.MakeSemVer("10020003:derpy")
	if err != nil {
		t.Errorf("Could not compose semantic version: %s", err.Error())
	}

	oldArgs := os.Args
	os.Args = []string{
		"wibble",
		"-config", path + file,
		"-log", path + log,
		"-debug",
	}

	acnf := &AccessorConfig{}
	conf := NewConfig(name, vers, acnf, nil, false)
	conf.Parse()
	os.Args = oldArgs

	t.Run("IsDebug", func(t *testing.T) {
		if conf.IsDebug() != true {
			t.Errorf("Unexpected debug value: %#v", conf.IsDebug())
		}
	})

	t.Run("Name", func(t *testing.T) {
		if conf.Name() != name {
			t.Errorf("Unexpected name value: %#v", conf.Name())
		}
	})

	t.Run("Version", func(t *testing.T) {
		if conf.Version() != vers {
			t.Errorf("Unexpected version value: %#v", conf.Version())
		}
	})

	t.Run("ConfFile", func(t *testing.T) {
		if conf.ConfFile() != path+file {
			t.Errorf("Unexpected config file value: %#v", conf.ConfFile())
		}
	})

	t.Run("LogFile", func(t *testing.T) {
		if conf.LogFile() != path+log {
			t.Errorf("Unexpected log file value: %#v", conf.LogFile())
		}
	})
}

// config_test.go ends here.
