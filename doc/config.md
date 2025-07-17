-*- Mode: gfm -*-

# config -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/config"
```

## Usage

```go
var (
	ErrValidatorNotFound  = errors.Base("validator not found")
	ErrIncorrectArguments = errors.Base("incorrect arguments")
	ErrValidationFailed   = errors.Base("validation failed")
)
```

```go
var (
	ErrInvalidObject = errors.Base("invalid configuration object")
)
```

#### type Config

```go
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
```

Config structure.

User-defined configuration options are placed into `App`, with any custom
validators stuffed into `Validators`.

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

Options will be parsed during `init`. Any validation or JSON errors will result
in the program exiting with error information dumped to stdout.

It is worth noting that there are three special structure tags:

* `config_hide`: Field is hidden when the config is dumped to string,

* `config_obscure`: Field is obscured with asterisks when dumped,

* `config_validator`: The validation function for the field.

#### func  Init

```go
func Init(
	name string,
	version *semver.SemVer,
	data any,
	fns ValidatorsMap,
	required bool,
) Config
```
Init a new configuration instance.

#### func  NewConfig

```go
func NewConfig(
	name string,
	version *semver.SemVer,
	data any,
	fns ValidatorsMap,
	required bool,
) Config
```
Create a new configuration instance.

#### func  NewDefaultConfig

```go
func NewDefaultConfig(required bool) Config
```
Create a new empty `Config` instance.

#### type ValidatorsMap

```go
type ValidatorsMap map[string]interface{}
```
