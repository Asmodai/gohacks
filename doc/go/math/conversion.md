<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# conversion -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/conversion"
```

## Usage

#### func  BToGiB

```go
func BToGiB(b uint64) uint64
```
Convert bytes to gibibytes.

#### func  BToKiB

```go
func BToKiB(b uint64) uint64
```
Convert bytes to kibibytes.

#### func  BToMiB

```go
func BToMiB(b uint64) uint64
```
Convert bytes to mebibytes.

#### func  BToTiB

```go
func BToTiB(b uint64) uint64
```
Convert bytes to tebibytes.

#### func  GiBToB

```go
func GiBToB(b uint64) uint64
```
Convert gibibytes to bytes.

#### func  KiBToB

```go
func KiBToB(b uint64) uint64
```
Convert kibibytes to bytes.

#### func  MiBToB

```go
func MiBToB(b uint64) uint64
```
Convert mebibytes to bytes.

#### func  TiBToB

```go
func TiBToB(b uint64) uint64
```
Convert tebibytes to bytes.

#### func  ToFloat64

```go
func ToFloat64(val any) (float64, bool)
```
Convert a value to a 64-bit floating-point value.
