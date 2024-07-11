-*- Mode: gfm -*-

# utils -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/utils"
```

## Usage

```go
const (
	ElideSuffix string = "..."
)
```

```go
const (
	PadPadding string = " "
)
```

#### func  FormatDuration

```go
func FormatDuration(d time.Duration) string
```

#### func  GetEnv

```go
func GetEnv(key, def string) string
```

#### func  Member

```go
func Member[T constraints.Ordered](vs []T, elt T) bool
```

#### func  Pop

```go
func Pop(array []string) (string, []string)
```

#### func  Substr

```go
func Substr(input string, start int, length int) string
```

#### type Elidable

```go
type Elidable string
```


#### func (Elidable) Elide

```go
func (s Elidable) Elide(max int) string
```

#### type Padable

```go
type Padable string
```


#### func (Padable) Pad

```go
func (p Padable) Pad(padding int) string
```
