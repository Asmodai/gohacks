/*
 * config_test.go --- Config tests.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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

package config

import (
	"github.com/Asmodai/gohacks/semver"

	"fmt"
	"os"
	"testing"
)

var DummyConfigJson = "{\"testing\":{\"opt1\":\"testing\",\"opt2\":42}}"

var version = &semver.SemVer{0, 1, 0, "test"}

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

func TestSimple(t *testing.T) {
	t.Log("Can we construct a simple config object?")

	path, err := os.Getwd()
	if err != nil {
		t.Errorf("getwd: %s", err.Error())
		return
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"wibble", "-config", path + "/../testing/conf.json"}

	opts := &DummyConfig{}
	conf := Init("Test", version, opts, MakeFns())
	if conf == nil {
		t.Error("No, nil config object!")
	}
	conf.Parse()

	t.Log("Check debug is false")
	if conf.IsDebug() == false {
		t.Log("Yes.")
	} else {
		t.Error("No, debug is true!")
		return
	}

	t.Log("Is string representation as expected?")
	if conf.String() == StringRep {
		t.Log("Yes.")
	} else {
		t.Errorf("No, got:\n%v\bWanted:\n%v\n", conf.String(), StringRep)
		return
	}
}

func TestValidators(t *testing.T) {
	t.Log("Can we add new validators on the fly?")

	c := NewConfig()
	f := func() {}
	c.AddValidator("test", f)

	if c.Validators["test"] != nil {
		t.Log("Yes.")
	} else {
		t.Error("No.")
	}
}

func TestCLIFlags(t *testing.T) {
	t.Log("Can we manipulate CLI flags?")

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
	c := Init("Test", version, o, MakeFns())

	c.AddBoolFlag(&o.Flags.BoolFlag, "bool", false, "bool")
	c.AddFloat64Flag(&o.Flags.F64Flag, "float64", 0.0, "float64")
	c.AddIntFlag(&o.Flags.IntFlag, "int", 0, "int")
	c.AddInt64Flag(&o.Flags.I64Flag, "int64", 0, "int64")
	c.AddStringFlag(&o.Flags.StringFlag, "string", "", "string")
	c.AddUintFlag(&o.Flags.UIntFlag, "uint", 0, "uint")
	c.AddUint64Flag(&o.Flags.UI64Flag, "uint64", 0, "uint64")

	c.Parse()
	os.Args = oldArgs

	flag := c.LookupFlag("uint")
	if flag == nil {
		t.Error("`LookupFlag` did not work.")
	}

	if o.Flags.BoolFlag == true {
		t.Log("Bool yes.")
	} else {
		t.Error("Bool no!")
	}

	if o.Flags.F64Flag == 1.23 {
		t.Log("Float64 yes.")
	} else {
		t.Error("Float64 no!")
	}

	if o.Flags.IntFlag == -2 {
		t.Log("Int yes.")
	} else {
		t.Error("Int no!")
	}

	if o.Flags.I64Flag == 3456 {
		t.Log("Int64 yes.")
	} else {
		t.Error("Int64 no!")
	}

	if o.Flags.StringFlag == "seven" {
		t.Log("String yes.")
	} else {
		t.Error("String no!")
	}

	if o.Flags.UIntFlag == 8 {
		t.Log("UInt yes.")
	} else {
		t.Error("UInt no!")
	}

	if o.Flags.UI64Flag == 90210 {
		t.Log("UInt64 yes.")
	} else {
		t.Error("UInt64 no!")
	}
}

/* config_test.go ends here. */
