// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// cli.go --- CLI flags and parsing.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package config

import (
	"flag"
	"fmt"
	"os"
)

// Add a boolean flag.
func (c *config) AddBoolFlag(p *bool, name string, value bool, usage string) {
	c.flags.BoolVar(p, name, value, usage)
}

// Add a float64 flag.
func (c *config) AddFloat64Flag(p *float64, name string, value float64, usage string) {
	c.flags.Float64Var(p, name, value, usage)
}

// Add an integer flag.
func (c *config) AddIntFlag(p *int, name string, value int, usage string) {
	c.flags.IntVar(p, name, value, usage)
}

// Add a 64-bit integer flag.
func (c *config) AddInt64Flag(p *int64, name string, value int64, usage string) {
	c.flags.Int64Var(p, name, value, usage)
}

// Add a string flag.
func (c *config) AddStringFlag(p *string, name string, value string, usage string) {
	c.flags.StringVar(p, name, value, usage)
}

// Add an unsigned integer flag.
func (c *config) AddUintFlag(p *uint, name string, value uint, usage string) {
	c.flags.UintVar(p, name, value, usage)
}

// Add an unsigned 64-bit integer flag.
func (c *config) AddUint64Flag(p *uint64, name string, value uint64, usage string) {
	c.flags.Uint64Var(p, name, value, usage)
}

// Look up a flag by its name.
func (c *config) LookupFlag(name string) *flag.Flag {
	return c.flags.Lookup(name)
}

// Parse CLI options.
func (c *config) addFlags() {
	c.flags.BoolVar(&c.ConfigCLI.Debug, "debug", false, "Debug mode")
	c.flags.BoolVar(&c.ConfigCLI.Version, "version", false, "Print version and exit")
	c.flags.BoolVar(&c.ConfigCLI.Dump, "dump", false, "Dump config to stdout and exit")
	c.flags.StringVar(&c.ConfigCLI.ConfFile, "config", "", "Configuration file")
	c.flags.StringVar(&c.ConfigCLI.LogFile, "log", "", "Log file")
}

// Perform ations for specific CLI options.
//
//nolint:forbidigo
func (c *config) handleCLI() {
	if c.ConfigCLI.Version {
		fmt.Printf(
			"This is %s, version %s (%s)\n",
			c.ConfigApp.Name,
			c.ConfigApp.Version,
			c.ConfigApp.Version.Commit,
		)
		os.Exit(0)
	}

	if c.ConfigCLI.Dump {
		fmt.Printf("%s\n", c.String())
		os.Exit(0)
	}
}

// Parse config and CLI flags.
//
//nolint:forbidigo
func (c *config) Parse() {
	var err []error

	//nolint:errcheck
	c.flags.Parse(os.Args[1:])

	// Check if we're something that should just print and exit here.
	if c.ConfigCLI.Version {
		goto only_handle
	}

	c.load()

	err = c.Validate()
	if len(err) > 0 {
		fmt.Printf("Error(s) when parsing %s:\n\n", c.ConfigCLI.ConfFile)

		for i, e := range err {
			fmt.Printf("%5d) %s\n", (i + 1), e.Error())
		}

		fmt.Printf("\n")

		// XXX Maybe this shouldn't be responsible for exiting.
		// Could have a failure chain that leads back to the
		// user's code for them to make the decision on whether
		// to exit back to the OS with an error.
		os.Exit(1)
	}

only_handle:
	c.handleCLI()
}

// cli.go ends here.
