-*- Mode: gfm -*-

# utils -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/utils"
```

## Usage

```go
const (
	ElideSuffix    string = "..."
	ElideSuffixLen int    = 3
)
```

```go
const (
	PadPadding string = " "
)
```

#### func  Elide

```go
func Elide(str string, max int) string
```

#### func  FormatDuration

```go
func FormatDuration(d time.Duration) string
```

#### func  GetEnv

```go
func GetEnv(key, def string) string
```

#### func  Pad

```go
func Pad(str string, padding int) string
```

#### func  Pop

```go
func Pop(array []string) (string, []string)
```

#### func  Substr

```go
func Substr(input string, start int, length int) string
```
