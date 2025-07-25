-*- Mode: gfm -*-

# utils -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/utils"
```

## Usage

```go
const (
	// Suffix used to denote elided strings.
	ElideSuffix string = "..."

	// Length of the elide suffix.
	ElideSuffixLen int = 3
)
```

```go
const (
	// Character to use when padding strings.
	PadPadding string = " "
)
```

#### func  Elide

```go
func Elide(str string, maxima int) string
```
Return a string that has been elided to the given length.

#### func  FormatDuration

```go
func FormatDuration(dur time.Duration) string
```
Format a time duration in pretty format.

Example, a duration of 72 minutes becomes "1 hour(s), 12 minute(s)".

#### func  GetEnv

```go
func GetEnv(key, def string) string
```
Look up the given key in the underlying operating system's environment and
return its value.

Should the key not exist, then the given default value is returned instead.

#### func  Pad

```go
func Pad(str string, padding int) string
```
Pad the given string with the given number of spaces.

#### func  Pop

```go
func Pop(array []string) (string, []string)
```
Pop the last element from an array, returning the element and a copy of the
original array with the last element removed.

#### func  Substr

```go
func Substr(input string, start int, length int) (string, bool)
```
Return a substring of the given input string that starts at the given start
point and has the given length.

#### type Filesystem

```go
type Filesystem struct {
}
```

Filesystem

Allow fort he checking of filesystem entities without having to open, stat,
close continuously.

This is a wrapper around various methods from the `os` package that is designed
to allow repeated querying without having to re-open and re-stat the entity.

If used within a struct, this will not provide any sort of live updates and it
does not utilise `epoll` (or any other similar facility) to update if a file or
directory structure changes.

Maybe one day it will facilitate live updates et al to track file changes, but
that day is not today.

#### func  NewFilesystem

```go
func NewFilesystem(path string) *Filesystem
```
Create a new filesystem object.

#### func (*Filesystem) Exists

```go
func (fs *Filesystem) Exists() bool
```
Does the entity exist?

#### func (*Filesystem) IsDirectory

```go
func (fs *Filesystem) IsDirectory() bool
```
Is the entity a directory?

#### func (*Filesystem) IsExecutable

```go
func (fs *Filesystem) IsExecutable() bool
```
Does the entity's permission bits include any type of "executable"?

#### func (*Filesystem) IsFile

```go
func (fs *Filesystem) IsFile() bool
```
Is the entity a file?

#### func (*Filesystem) IsGroupExecutable

```go
func (fs *Filesystem) IsGroupExecutable() bool
```
Does the entity's permission bits include "group executable"?

#### func (*Filesystem) IsOtherExecutable

```go
func (fs *Filesystem) IsOtherExecutable() bool
```
Does the entity's permission bits include "other executable"?

#### func (*Filesystem) IsOwnerExecutable

```go
func (fs *Filesystem) IsOwnerExecutable() bool
```
Does the entity's permission bits include "owner executable"?

#### func (*Filesystem) Mode

```go
func (fs *Filesystem) Mode() os.FileMode
```
Return the entity's mode.

#### func (*Filesystem) Name

```go
func (fs *Filesystem) Name() string
```
Return the entity's name.

#### func (*Filesystem) Size

```go
func (fs *Filesystem) Size() int64
```
Return the entity's size.
