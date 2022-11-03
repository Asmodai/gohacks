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

#### func  All

```go
func All[T Numeric](vs []T, fn func(T) bool) bool
```
Run predicate `fn` on all elems of `vs` and return true if all elems of `vs`
match the predicate.

#### func  Any

```go
func Any[T Numeric](vs []T, fn func(T) bool) bool
```
Run predicate `fn` on all elems of `vs` and return true if any elems of `vs`
match the predicate.

#### func  FormatDuration

```go
func FormatDuration(d time.Duration) string
```

#### func  HasFraction

```go
func HasFraction(num float64) bool
```
Does the given float have a faction?

#### func  Member

```go
func Member[T constraints.Ordered](vs []T, elt T) bool
```

#### func  Number

```go
func Number[V Numeric](val V) interface{}
```
If the given thing does not have a fraction, convert it to an int.

#### func  Pop

```go
func Pop(array []string) (string, []string)
```

#### func  Substr

```go
func Substr(input string, start int, length int) string
```

#### func  ValueOf

```go
func ValueOf(thing interface{}) interface{}
```
Return the value of an interface.

#### type Elidable

```go
type Elidable string
```


#### func (Elidable) Elide

```go
func (s Elidable) Elide(max int) string
```

#### type Numeric

```go
type Numeric interface {
	~int | ~int64 | ~float64
}
```


#### type Padable

```go
type Padable string
```


#### func (Padable) Pad

```go
func (p Padable) Pad(padding int) string
```
