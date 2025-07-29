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
