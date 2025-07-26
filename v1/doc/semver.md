<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# semver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/semver"
```

## Usage

```go
var (
	ErrInvalidVersion = errors.Base("invalid version")
)
```

#### type SemVer

```go
type SemVer struct {
	Major  int    // Major version number.
	Minor  int    // Minor version number.
	Patch  int    // Patch number.
	Commit string // VCS commit identifier.
}
```

Semantic version structure.

#### func  MakeSemVer

```go
func MakeSemVer(info string) (*SemVer, error)
```
Make a new semantic version from the given string.

#### func  NewSemVer

```go
func NewSemVer() *SemVer
```
Create a new empty semantic version object.

#### func (*SemVer) FromString

```go
func (s *SemVer) FromString(info string) error
```
Convert numeric version to components.

#### func (*SemVer) String

```go
func (s *SemVer) String() string
```
Return a string representation of the semantic version.

#### func (*SemVer) Version

```go
func (s *SemVer) Version() int
```
Return an integer version.
