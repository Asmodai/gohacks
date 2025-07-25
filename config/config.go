// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- Configuration.
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
//
// mock:yes

package config

import (
	"github.com/Asmodai/gohacks/semver"

	json "github.com/goccy/go-json"
	"gitlab.com/tozd/go/errors"

	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
)

var (
	ErrInvalidObject = errors.Base("invalid configuration object")
)

type ValidatorsMap map[string]interface{}

/*
Config structure.

User-defined configuration options are placed into `App`, with any
custom validators stuffed into `Validators`.

The magic happens like this.

1) Define the structure you want your options to live in:

```go

	type Options struct {
	    Option1 string `json:"option1" config_validator:"ValidateOption1"
	    // ...
	}

```

The `config_validator` tag informs the Config module that you wish to validate
the `Option1` field using the `ValidateOption1` function.

2) Define your validators:

```go

	func ValidateOption1(value string) error {
	    if value == "" {
		return fmt.Errorf("Noooooooooooo!")
	    }

	    return nil
	}

```

The validator *must* return an `error` or `nil`.

3) Set it all up:

```go

	func main() {
	    // ...

	    opts := &Options{
		// ...
	    }

	    fns := map[string]interface{}{
		"ValidateOption1": ValidateOption1,
	    }

	    vers := &semver.SemVer{1, 2, 3, "herpderp"}
	    conf := config.Init("App Name", vers, opts, fns)
	    conf.Parse()

	    // ...
	}

```

Options will be parsed during `init`.  Any validation or JSON errors will
result in the program exiting with error information dumped to stdout.

It is worth noting that there are three special structure tags:

* `config_hide`:      Field is hidden when the config is dumped to string,

* `config_obscure`:   Field is obscured with asterisks when dumped,

* `config_validator`: The validation function for the field.
*/
type Config interface {
	// Have we detected `-debug` in the CLI arguments?
	IsDebug() bool

	// Returns the application name.
	Name() string

	// Returns the application's version.
	Version() *semver.SemVer

	// Return the pathname to a logfile passed via the `-log` CLI
	// argument.
	//
	// Will be empty if no log file has been specified.
	LogFile() string

	// Return the pathname to the configuration file passed via the
	// '-config' CLI argument.
	//
	// Will be empty if no configuration file has been specified.
	ConfFile() string

	// Return the user-provided configuration structure.
	AppConfig() any

	// Add a validator function for an option named `name`.  `fn` must be
	// a function that can be funcalled and must return either `nil` or an
	// error.
	AddValidator(name string, fn any)

	// Return the string representation of this object.
	String() string

	// Perform validation on the application configuration structure that
	// is parsed from the config file.
	Validate() []error

	// Add a Boolean CLI flag.
	AddBoolFlag(p *bool, name string, value bool, usage string)

	// Add a 64-bit floating point CLI flag.
	AddFloat64Flag(p *float64, name string, value float64, usage string)

	// Add a 32-bit signed integer CLI flag.
	AddIntFlag(p *int, name string, value int, usage string)

	// Add a 64-bit signed integer CLI flag.
	AddInt64Flag(p *int64, name string, value int64, usage string)

	// Add a string CLI flag.
	AddStringFlag(p *string, name string, value string, usage string)

	// Add a 32-bit unsigned integer CLI flag.
	AddUintFlag(p *uint, name string, value uint, usage string)

	// Add a 64-bit unsigned integer CLI flag.
	AddUint64Flag(p *uint64, name string, value uint64, usage string)

	// Look up a given flag value.  Looks for a flag called `name`.
	LookupFlag(name string) *flag.Flag

	// Parse CLI arguments.
	Parse()
}

// Internal structure.
type config struct {
	// Application information.
	ConfigApp struct {
		Name    string
		Version *semver.SemVer
	} `config_hide:"true"`

	// CLI flags
	ConfigCLI struct {
		Debug    bool   `config_hide:"true"`
		Version  bool   `config_hide:"true"`
		Dump     bool   `config_hide:"true"`
		ConfFile string `config_hide:"true"`
		LogFile  string `config_hide:"true"`
	} `config_hide:"true"`

	App        any           `config_hide:"true"`
	Validators ValidatorsMap `config_hide:"true"`

	flags       *flag.FlagSet `config_hide:"true"`
	mustHaveCLI bool          `config_hide:"true"`
}

// Create a new empty `Config` instance.
func NewDefaultConfig(required bool) Config {
	return &config{
		Validators:  make(ValidatorsMap),
		mustHaveCLI: required,
	}
}

// Create a new configuration instance.
func NewConfig(
	name string,
	version *semver.SemVer,
	data any,
	fns ValidatorsMap,
	required bool,
) Config {
	obj, ok := NewDefaultConfig(required).(*config)
	if !ok {
		panic(ErrInvalidObject)
	}

	obj.ConfigApp.Name = name
	obj.ConfigApp.Version = version
	obj.App = data
	obj.Validators = fns
	obj.flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	obj.addFlags()

	return obj
}

// Return the application's name.
func (c *config) Name() string {
	return c.ConfigApp.Name
}

// Return the application's version.
func (c *config) Version() *semver.SemVer {
	return c.ConfigApp.Version
}

// Is debug mode enabled?
func (c *config) IsDebug() bool {
	return c.ConfigCLI.Debug
}

// Return the path to the specified logging file.
func (c *config) LogFile() string {
	return c.ConfigCLI.LogFile
}

// Return the user-provided configuration structure.
func (c *config) AppConfig() any {
	return c.App
}

// Return the path to the configuration file.
func (c *config) ConfFile() string {
	return c.ConfigCLI.ConfFile
}

// Pretty-print the configuration.
func (c *config) String() string {
	// Attempt to initialise the user config.
	c.checkCanInit(reflect.ValueOf(c.App))

	s := "Configuration:"
	s += c.recursePrint("    ", reflect.ValueOf(c), make(map[any]bool))
	s += c.recursePrint("    ", reflect.ValueOf(c.App), make(map[any]bool))

	return s
}

// Load JSON config file.
//
//nolint:forbidigo
func (c *config) load() {
	if c.ConfigCLI.ConfFile == "" {
		if c.mustHaveCLI {
			fmt.Printf("A configuration file must be provided via `-config`.\n")
			os.Exit(1)
		}

		return
	}

	file, err := os.Open(c.ConfigCLI.ConfFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	if len(bytes) == 0 {
		return
	}

	err = json.Unmarshal(bytes, c.App)
	if err != nil {
		fmt.Printf("Error parsing configuration file: %s\n", err.Error())
		file.Close()
		os.Exit(1) //nolint:gocritic
	}
}

// Init a new configuration instance.
func Init(
	name string,
	version *semver.SemVer,
	data any,
	fns ValidatorsMap,
	required bool,
) Config {
	return NewConfig(name, version, data, fns, required)
}

// config.go ends here.
