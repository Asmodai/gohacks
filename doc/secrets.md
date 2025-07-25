<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

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
	ErrFileNotFound   = errors.Base("file not found")
	ErrNotAFile       = errors.Base("not a file")
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

#### func (*Secret) Value

```go
func (s *Secret) Value() string
```
