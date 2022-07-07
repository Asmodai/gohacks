-*- Mode: gfm -*-

# semver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/semver"
```

## Usage

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
Convert numeric version to components

#### func (*SemVer) String

```go
func (s *SemVer) String() string
```

#### func (*SemVer) Version

```go
func (s *SemVer) Version() int
```
