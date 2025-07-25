<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# semver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/semver"
```

## Usage

```go
const (
	MAGICMAJOR        = 10000000
	MAGICMINOR        = 10000
	MAGICMAJORTOMINOR = 1000
)
```

```go
var (
	ErrInvalidVersion = errors.Base("invalid version")
)
```

#### type SemVer

```go
type SemVer struct {
	Major  int
	Minor  int
	Patch  int
	Commit string
}
```


#### func  MakeSemVer

```go
func MakeSemVer(info string) (*SemVer, error)
```

#### func  NewSemVer

```go
func NewSemVer() *SemVer
```

#### func (*SemVer) FromString

```go
func (s *SemVer) FromString(info string) error
```
Convert numeric version to components.

#### func (*SemVer) String

```go
func (s *SemVer) String() string
```

#### func (*SemVer) Version

```go
func (s *SemVer) Version() int
```
