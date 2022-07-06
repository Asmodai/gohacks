-*- Mode: gfm -*-

# config -- Asmodai's Go Hacks

```go
    import "github.com/Asmodai/gohacks/config"
```

## Usage

#### type Config

```go
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
}
```

Config structure.

User-defined configuration options are placed into `App`, with any custom
validators stuffed into `Validators`.

The magic happens like this.

1) Define the structure you want your options to live in: ```go

    type Options struct {
        Option1 string `json:"option1" config_validator:"ValidateOption1"
        // ...
    }

``` The `config_validator` tag informs the Config module that you wish to
validate the `Option1` field using the `ValidateOption1` function.

2) Define your validators: ```go

    func ValidateOption1(value string) error {
        if value == "" {
            return fmt.Errorf("Noooooooooooo!")
        }

        return nil
    }

``` The validator *must* return an `error` or `nil`.

3) Set it all up: ```go

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

``` Options will be parsed during `init`. Any validation or JSON errors will
result in the program exiting with error information dumped to stdout.

It is worth noting that there are three special structure tags: * `config_hide`:
Field is hidden when the config is dumped to string, * `config_obscure`: Field
is obscured with asterisks when dumped, * `config_validator`: The validation
function for the field.

#### func  Init

```go
func Init(name string, version *semver.SemVer, data interface{}, fns ValidatorsMap) *Config
```
Init a new configuration instance.

#### func  NewConfig

```go
func NewConfig() *Config
```
Create a new empty `Config` instance.

#### func (*Config) AddBoolFlag

```go
func (c *Config) AddBoolFlag(p *bool, name string, value bool, usage string)
```
Add a boolean flag.

#### func (*Config) AddFloat64Flag

```go
func (c *Config) AddFloat64Flag(p *float64, name string, value float64, usage string)
```
Add a float64 flag.

#### func (*Config) AddInt64Flag

```go
func (c *Config) AddInt64Flag(p *int64, name string, value int64, usage string)
```
Add a 64-bit integer flag.

#### func (*Config) AddIntFlag

```go
func (c *Config) AddIntFlag(p *int, name string, value int, usage string)
```
Add a integer flag.

#### func (*Config) AddStringFlag

```go
func (c *Config) AddStringFlag(p *string, name string, value string, usage string)
```
Add a string flag.

#### func (*Config) AddUint64Flag

```go
func (c *Config) AddUint64Flag(p *uint64, name string, value uint64, usage string)
```
Add a unsigned 64-bit integer flag.

#### func (*Config) AddUintFlag

```go
func (c *Config) AddUintFlag(p *uint, name string, value uint, usage string)
```
Add a unsigned integer flag.

#### func (*Config) AddValidator

```go
func (c *Config) AddValidator(name string, fn interface{})
```
Add a validator function for a tag value.

This is so one can add validators after instance creation.

#### func (*Config) IsDebug

```go
func (c *Config) IsDebug() bool
```
Is debug mode enabled?

#### func (*Config) LogFile

```go
func (c *Config) LogFile() string
```
Return the path to the specified logging file.

#### func (*Config) LookupFlag

```go
func (c *Config) LookupFlag(name string) *flag.Flag
```
Look up a flag by its name.

#### func (*Config) Parse

```go
func (c *Config) Parse()
```
Parse config and CLI flags.

#### func (*Config) String

```go
func (c *Config) String() string
```
Pretty-print the configuration.

#### func (*Config) Validate

```go
func (c *Config) Validate() []error
```
Validate configuration.

Should validation fail, then a list of errors is returned. Should validation
pass, an empty list is returned.

#### type IConfig

```go
type IConfig interface {
	IsDebug() bool
	AddValidator(name string, fn interface{})
	String() string
	Validate() []error
	AddBoolFlag(p *bool, name string, value bool, usage string)
	AddFloat64Flag(p *float64, name string, value float64, usage string)
	AddIntFlag(p *int, name string, value int, usage string)
	AddInt64Flag(p *int64, name string, value int64, usage string)
	AddStringFlag(p *string, name string, value string, usage string)
	AddUintFlag(p *uint, name string, value uint, usage string)
	AddUint64Flag(p *uint64, name string, value uint64, usage string)
	LookupFlag(name string) *flag.Flag
	Parse()
	LogFile() string
}
```

Config interface.

#### type ValidatorsMap

```go
type ValidatorsMap map[string]interface{}
```
