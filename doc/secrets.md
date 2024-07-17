-*- Mode: gfm -*-

# secrets -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/secrets"
```

## Usage

```go
const (
	SecretsPath string = "/run/secrets"
)
```

```go
var (
	ErrNoPathSet      = errors.Base("no secrets path set")
	ErrZeroLengthFile = errors.Base("file has zero length")
)
```

#### type Secret

```go
type Secret struct {
}
```


#### func  Make

```go
func Make(file string) *Secret
```

#### func  New

```go
func New() *Secret
```

#### func (*Secret) Path

```go
func (s *Secret) Path() string
```

#### func (*Secret) Probe

```go
func (s *Secret) Probe() error
```

#### func (*Secret) SetPath

```go
func (s *Secret) SetPath(val string) error
```

#### func (*Secret) Value

```go
func (s *Secret) Value() string
```
