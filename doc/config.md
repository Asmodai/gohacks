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

#### type MockConfig

```go
type MockConfig struct {
}
```

MockConfig is a mock of Config interface.

#### func  NewMockConfig

```go
func NewMockConfig(ctrl *gomock.Controller) *MockConfig
```
NewMockConfig creates a new mock instance.

#### func (*MockConfig) AddBoolFlag

```go
func (m *MockConfig) AddBoolFlag(p *bool, name string, value bool, usage string)
```
AddBoolFlag mocks base method.

#### func (*MockConfig) AddFloat64Flag

```go
func (m *MockConfig) AddFloat64Flag(p *float64, name string, value float64, usage string)
```
AddFloat64Flag mocks base method.

#### func (*MockConfig) AddInt64Flag

```go
func (m *MockConfig) AddInt64Flag(p *int64, name string, value int64, usage string)
```
AddInt64Flag mocks base method.

#### func (*MockConfig) AddIntFlag

```go
func (m *MockConfig) AddIntFlag(p *int, name string, value int, usage string)
```
AddIntFlag mocks base method.

#### func (*MockConfig) AddStringFlag

```go
func (m *MockConfig) AddStringFlag(p *string, name, value, usage string)
```
AddStringFlag mocks base method.

#### func (*MockConfig) AddUint64Flag

```go
func (m *MockConfig) AddUint64Flag(p *uint64, name string, value uint64, usage string)
```
AddUint64Flag mocks base method.

#### func (*MockConfig) AddUintFlag

```go
func (m *MockConfig) AddUintFlag(p *uint, name string, value uint, usage string)
```
AddUintFlag mocks base method.

#### func (*MockConfig) AddValidator

```go
func (m *MockConfig) AddValidator(name string, fn any)
```
AddValidator mocks base method.

#### func (*MockConfig) ConfFile

```go
func (m *MockConfig) ConfFile() string
```
ConfFile mocks base method.

#### func (*MockConfig) EXPECT

```go
func (m *MockConfig) EXPECT() *MockConfigMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockConfig) IsDebug

```go
func (m *MockConfig) IsDebug() bool
```
IsDebug mocks base method.

#### func (*MockConfig) LogFile

```go
func (m *MockConfig) LogFile() string
```
LogFile mocks base method.

#### func (*MockConfig) LookupFlag

```go
func (m *MockConfig) LookupFlag(name string) *flag.Flag
```
LookupFlag mocks base method.

#### func (*MockConfig) Name

```go
func (m *MockConfig) Name() string
```
Name mocks base method.

#### func (*MockConfig) Parse

```go
func (m *MockConfig) Parse()
```
Parse mocks base method.

#### func (*MockConfig) String

```go
func (m *MockConfig) String() string
```
String mocks base method.

#### func (*MockConfig) Validate

```go
func (m *MockConfig) Validate() []error
```
Validate mocks base method.

#### func (*MockConfig) Version

```go
func (m *MockConfig) Version() *semver.SemVer
```
Version mocks base method.

#### type MockConfigMockRecorder

```go
type MockConfigMockRecorder struct {
}
```

MockConfigMockRecorder is the mock recorder for MockConfig.

#### func (*MockConfigMockRecorder) AddBoolFlag

```go
func (mr *MockConfigMockRecorder) AddBoolFlag(p, name, value, usage any) *gomock.Call
```
AddBoolFlag indicates an expected call of AddBoolFlag.

#### func (*MockConfigMockRecorder) AddFloat64Flag

```go
func (mr *MockConfigMockRecorder) AddFloat64Flag(p, name, value, usage any) *gomock.Call
```
AddFloat64Flag indicates an expected call of AddFloat64Flag.

#### func (*MockConfigMockRecorder) AddInt64Flag

```go
func (mr *MockConfigMockRecorder) AddInt64Flag(p, name, value, usage any) *gomock.Call
```
AddInt64Flag indicates an expected call of AddInt64Flag.

#### func (*MockConfigMockRecorder) AddIntFlag

```go
func (mr *MockConfigMockRecorder) AddIntFlag(p, name, value, usage any) *gomock.Call
```
AddIntFlag indicates an expected call of AddIntFlag.

#### func (*MockConfigMockRecorder) AddStringFlag

```go
func (mr *MockConfigMockRecorder) AddStringFlag(p, name, value, usage any) *gomock.Call
```
AddStringFlag indicates an expected call of AddStringFlag.

#### func (*MockConfigMockRecorder) AddUint64Flag

```go
func (mr *MockConfigMockRecorder) AddUint64Flag(p, name, value, usage any) *gomock.Call
```
AddUint64Flag indicates an expected call of AddUint64Flag.

#### func (*MockConfigMockRecorder) AddUintFlag

```go
func (mr *MockConfigMockRecorder) AddUintFlag(p, name, value, usage any) *gomock.Call
```
AddUintFlag indicates an expected call of AddUintFlag.

#### func (*MockConfigMockRecorder) AddValidator

```go
func (mr *MockConfigMockRecorder) AddValidator(name, fn any) *gomock.Call
```
AddValidator indicates an expected call of AddValidator.

#### func (*MockConfigMockRecorder) ConfFile

```go
func (mr *MockConfigMockRecorder) ConfFile() *gomock.Call
```
ConfFile indicates an expected call of ConfFile.

#### func (*MockConfigMockRecorder) IsDebug

```go
func (mr *MockConfigMockRecorder) IsDebug() *gomock.Call
```
IsDebug indicates an expected call of IsDebug.

#### func (*MockConfigMockRecorder) LogFile

```go
func (mr *MockConfigMockRecorder) LogFile() *gomock.Call
```
LogFile indicates an expected call of LogFile.

#### func (*MockConfigMockRecorder) LookupFlag

```go
func (mr *MockConfigMockRecorder) LookupFlag(name any) *gomock.Call
```
LookupFlag indicates an expected call of LookupFlag.

#### func (*MockConfigMockRecorder) Name

```go
func (mr *MockConfigMockRecorder) Name() *gomock.Call
```
Name indicates an expected call of Name.

#### func (*MockConfigMockRecorder) Parse

```go
func (mr *MockConfigMockRecorder) Parse() *gomock.Call
```
Parse indicates an expected call of Parse.

#### func (*MockConfigMockRecorder) String

```go
func (mr *MockConfigMockRecorder) String() *gomock.Call
```
String indicates an expected call of String.

#### func (*MockConfigMockRecorder) Validate

```go
func (mr *MockConfigMockRecorder) Validate() *gomock.Call
```
Validate indicates an expected call of Validate.

#### func (*MockConfigMockRecorder) Version

```go
func (mr *MockConfigMockRecorder) Version() *gomock.Call
```
Version indicates an expected call of Version.

#### type ValidatorsMap

```go
type ValidatorsMap map[string]interface{}
```
