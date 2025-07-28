<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# generics -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/generics"
```

## Usage

#### func  All

```go
func All[T any](vs []T, fn func(T) bool) bool
```
Run predicate `fn` on all elems of `vs` and return true if all elems of `vs`
match the predicate.

#### func  Any

```go
func Any[T any](vs []T, fn func(T) bool) bool
```
Run predicate `fn` on all elems of `vs` and return true if any elems of `vs`
match the predicate.

#### func  CoerceInt

```go
func CoerceInt[N Numeric](val N) any
```
Coerce a numeric to an integer value of appropriate signage and size.

#### func  HasFraction

```go
func HasFraction(num float64) bool
```
Does the given float have a faction?

#### func  MapList

```go
func MapList[T any](lst []T, fn func(T) bool) []T
```
Map a function on all elements of a list returning a new list containing all the
values from the initial list for which the function returns `true`.

#### func  MapMap

```go
func MapMap[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V
```
Apply a function on all key/value pairs in a map and return a new map of
key/value pairs containing all elements for which the provided function returns
`true`.

#### func  Member

```go
func Member[T cmp.Ordered](ordered []T, elt T) bool
```
Returns true if the container specified by `ordered` contains the member
specified by `elt`.

#### func  Number

```go
func Number[V Numeric](val V) any
```
If the given thing does not have a fraction, convert it to an int.

#### func  ValueOf

```go
func ValueOf(thing any) any
```
Attempt to ascertain the numeric value (if any) for a given thing.

If the type of `thing` is a floating-point number, then the value will be
converted via `Number` to either an explicit float or to an integer type if
there is no fraction.

If the type of `thing` is any other type, then it will be returned with no
conversion.

#### type Numeric

```go
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}
```

Generic "numeric" type constraint.
