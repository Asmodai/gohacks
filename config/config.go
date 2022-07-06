/*
 * config.go --- Configuration.
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
	"github.com/Asmodai/gohacks/types"

	"github.com/goccy/go-json"

	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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
type Config struct {
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

	App        interface{}   `config_hide:"true"`
	Validators ValidatorsMap `config_hide:"true"`

	flags *flag.FlagSet `config_hide:"true"`
}

// Create a new empty `Config` instance.
func NewConfig() *Config {
	return &Config{
		Validators: make(ValidatorsMap),
	}
}

// Is debug mode enabled?
func (c *Config) IsDebug() bool {
	return c.ConfigCLI.Debug
}

// Add a validator function for a tag value.
//
// This is so one can add validators after instance creation.
func (c *Config) AddValidator(name string, fn interface{}) {
	c.Validators[name] = fn
}

// Pretty-print the configuration.
func (c *Config) String() string {
	// Attempt to initialise the user config.
	c.checkCanInit(reflect.ValueOf(c.App))

	s := "Configuration:"
	s += c.recursePrint("    ", reflect.ValueOf(c), make(map[interface{}]bool))
	s += c.recursePrint("    ", reflect.ValueOf(c.App), make(map[interface{}]bool))

	return s
}

// Validate configuration.
//
// Should validation fail, then a list of errors is returned.
// Should validation pass, an empty list is returned.
func (c *Config) Validate() []error {
	sref := reflect.ValueOf(c.App).Elem()

	return c.recurseValidate(sref)
}

// Add a boolean flag.
func (c *Config) AddBoolFlag(p *bool, name string, value bool, usage string) {
	c.flags.BoolVar(p, name, value, usage)
}

// Add a float64 flag.
func (c *Config) AddFloat64Flag(p *float64, name string, value float64, usage string) {
	c.flags.Float64Var(p, name, value, usage)
}

// Add a integer flag.
func (c *Config) AddIntFlag(p *int, name string, value int, usage string) {
	c.flags.IntVar(p, name, value, usage)
}

// Add a 64-bit integer flag.
func (c *Config) AddInt64Flag(p *int64, name string, value int64, usage string) {
	c.flags.Int64Var(p, name, value, usage)
}

// Add a string flag.
func (c *Config) AddStringFlag(p *string, name string, value string, usage string) {
	c.flags.StringVar(p, name, value, usage)
}

// Add a unsigned integer flag.
func (c *Config) AddUintFlag(p *uint, name string, value uint, usage string) {
	c.flags.UintVar(p, name, value, usage)
}

// Add a unsigned 64-bit integer flag.
func (c *Config) AddUint64Flag(p *uint64, name string, value uint64, usage string) {
	c.flags.Uint64Var(p, name, value, usage)
}

// Look up a flag by its name.
func (c *Config) LookupFlag(name string) *flag.Flag {
	return c.flags.Lookup(name)
}

// Parse config and CLI flags.
func (c *Config) Parse() {
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

		os.Exit(1)
	}

only_handle:
	c.handleCLI()
}

// Return the path to the specified logging file.
func (c *Config) LogFile() string {
	return c.ConfigCLI.LogFile
}

// Call a given function with arguments, and return any error.
func (c *Config) call(field string, name string, params ...interface{}) error {
	if _, ok := c.Validators[name]; !ok {
		return types.NewError(
			"CONFIG",
			"Validator '%s' was not found.",
			name,
		)
	}
	fn := reflect.ValueOf(c.Validators[name])

	if len(params) != fn.Type().NumIn() {
		return types.NewError(
			"CONFIG",
			"Validator '%s' expects %d arguments but only %d given.",
			name,
			fn.Type().NumIn(),
			len(params),
		)
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result := fn.Call(in)
	if result[0].Interface() == nil {
		return nil
	}

	err := result[0].Interface().(error)
	return types.NewError(
		"CONFIG",
		"[%s] Validation failed on `%s`: %v",
		name,
		field,
		err,
	)
}

// Sometimes I think Go is fail.  I mean, even C# does reflection better.
func (c *Config) callMethod(value reflect.Value, method string) (interface{}, bool) {
	var final reflect.Value
	var ptr reflect.Value

	// If we're a pointer, then use the value of the pointee.
	if value.Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	}

	// Are we valid?
	if value.IsValid() {
		meth := value.MethodByName(method)

		// Better check the method is valid too
		if meth.IsValid() {
			final = meth
		}
	}

	// Are we a valid pointer?
	if ptr.IsValid() {
		meth := ptr.MethodByName(method)

		// Check the method too
		if meth.IsValid() {
			final = meth
		}
	}

	// Finally, double-check the method and invoke.
	if final.IsValid() {
		return final.Call([]reflect.Value{})[0].Interface(), true
	}

	return nil, false
}

// Attempt to call an `Init` method on a specific thing.
func (c *Config) checkCanInit(val reflect.Value) bool {
	_, ok := c.callMethod(val, "Init")

	return ok
}

// Ugly recursive nasty pretty printing.
func (c *Config) recursePrint(prefix string, val reflect.Value, visited map[interface{}]bool) string {
	var s string = ""

	toString, toStringFound := c.callMethod(val, "ToString")

	// Reflect over pointers and interfaces.
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.Kind() == reflect.Ptr {
			// If we're a pointer, check if we've visited the pointee.
			if visited[val.Interface()] {
				return s
			}

			// Tag it as visited.
			visited[val.Interface()] = true
		}

		// We want the pointee.
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		if toStringFound {
			// Thing has a `toString` method, so we use its output.
			s += fmt.Sprintf("%s%s", prefix, toString.(string))
		} else {
			// Sigh.
			t := val.Type()

			// Iterate over fields.
			for i := 0; i < val.NumField(); i++ {
				if t.Field(i).Tag.Get("config_hide") == "true" {
					// Ignore fields with the `config_hide` tag.
					continue
				}

				s += fmt.Sprintf("\n%s%s:", prefix, t.Field(i).Name)

				// Is the field exported?
				if !val.Field(i).CanSet() {
					// No, mark it so and ignore it.
					s += " <unexported>"
					continue
				}

				// Should we obscure the field's value?
				if t.Field(i).Tag.Get("config_obscure") == "true" {
					s += " [**********]"
				} else {
					// Not obscuring, recurse-print.
					s += c.recursePrint(prefix+"    ", val.Field(i), visited)
				}
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			s += c.recursePrint("\n"+prefix, val.Index(i), visited)
		}

	case reflect.Invalid:
		s += " nil"

	default:
		s += fmt.Sprintf(" [%v]", val.Interface())
	}

	return s
}

// Recursive ugly reflective validation.
func (c *Config) recurseValidate(v reflect.Value) []error {
	sref := v
	errs := []error{}

	for i := 0; i < sref.NumField(); i++ {
		field := sref.Field(i)
		ftype := sref.Type().Field(i)
		validate := ftype.Tag.Get("config_validator")

		// Nested structure?
		if field.Kind() == reflect.Struct {
			// Yep, recurse.
			nested := reflect.ValueOf(field.Interface())
			errs = append(errs, c.recurseValidate(nested)...)
		}

		// Is validation function valid?
		if validate != "" {
			result := c.call(ftype.Name, validate, field.Interface())
			if result != nil {
				errs = append(errs, []error{result}...)
			}
		}
	}

	return errs
}

// Parse CLI options.
func (c *Config) addFlags() {
	c.flags.BoolVar(&c.ConfigCLI.Debug, "debug", false, "Debug mode")
	c.flags.BoolVar(&c.ConfigCLI.Version, "version", false, "Print version and exit")
	c.flags.BoolVar(&c.ConfigCLI.Dump, "dump", false, "Dump config to stdout and exit")
	c.flags.StringVar(&c.ConfigCLI.ConfFile, "config", "", "Configuration file")
	c.flags.StringVar(&c.ConfigCLI.LogFile, "log", "", "Log file")
}

// Perform ations for specific CLI options.
func (c *Config) handleCLI() {
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

// Load JSON config file.
func (c *Config) load() {
	if c.ConfigCLI.ConfFile == "" {
		fmt.Println("Not loading a configuration file.")
		return
		//log.Fatal("A configuration file must be provided via `-config`.")
	}

	file, err := os.Open(c.ConfigCLI.ConfFile)
	if err != nil {
		panic(fmt.Errorf("Error loading config file: %s", err.Error()))
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	if len(bytes) == 0 {
		return
	}

	err = json.Unmarshal(bytes, c.App)
	if err != nil {
		fmt.Printf("Error parsing configuration file: %s\n", err.Error())
		os.Exit(1)
	}
}

// Init a new configuration instance.
func Init(name string, version *semver.SemVer, data interface{}, fns ValidatorsMap) *Config {
	inst := NewConfig()

	inst.ConfigApp.Name = name
	inst.ConfigApp.Version = version

	inst.App = data
	inst.Validators = fns

	inst.flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	inst.addFlags()

	return inst
}

/* config.go ends here. */
