<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# math -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/math"
```

## Usage

#### func  AbsI64

```go
func AbsI64(val int64) int64
```
Return the absolute value of a 64-bit signed integer value.

#### func  ClampI64

```go
func ClampI64(val, minVal, maxVal int64) int64
```
Clamp a signed 64-bit value to a minima and maxima.

#### func  FormatFloat64

```go
func FormatFloat64(num float64) string
```
Format a float.

If the float is NaN or infinite, then those are explicitly returned.

#### func  MaxI

```go
func MaxI(lhs, rhs int) int
```

#### func  MaxI32

```go
func MaxI32(lhs, rhs int32) int32
```

#### func  MaxI64

```go
func MaxI64(lhs, rhs int64) int64
```

#### func  MinI

```go
func MinI(lhs, rhs int) int
```

#### func  MinI32

```go
func MinI32(lhs, rhs int32) int32
```

#### func  MinI64

```go
func MinI64(lhs, rhs int64) int64
```

#### func  RoundF

```go
func RoundF(num float64, precision uint) float64
```
Rounds num to the given number of decimal places and returns the result as a
float64, using math.Round for IEEE-754 compliant rounding.

#### func  RoundI

```go
func RoundI(num float64) int
```
Rounds a 64-bit floating point number to the nearest integer, returning it as an
int. Values halfway between integers are rounded away from zero.

#### func  ToFixed

```go
func ToFixed(num float64, precision uint) float64
```
Rounds num to the given number of decimal places and returns the result as a
float64. Unlike RoundF, it uses integer rounding logic (via RoundI), which may
behave slightly differently around half-values.

#### func  WithinPlatform

```go
func WithinPlatform(value, defValue int64) int
```
