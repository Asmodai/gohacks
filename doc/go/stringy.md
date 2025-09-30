<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# stringy -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/stringy"
```

## Usage

#### func  BMH

```go
func BMH(needle, haystack []byte) int
```
Returns the byte index of the first occurrence of needle in haystack, or -1.

#### func  BMHRunes

```go
func BMHRunes(needle, haystack []rune) int
```

#### func  FindAllBM

```go
func FindAllBM(pattern, text []byte) []int
```
FindAllBM returns all (overlapping) byte indices where pattern occurs.

#### func  IndexBM

```go
func IndexBM(pattern, text []byte) int
```
Returns the byte index of the first occurrence of pattern in text, or -1.

#### func  IsHexadecimal

```go
func IsHexadecimal(thing rune) bool
```
Is the given rune a valid component of a hexadecimal number?

#### func  Levenshtein

```go
func Levenshtein(a, b string) int
```
Levenshtein returns the edit distance between a and b, counting insertions,
deletions, and substitutions (all cost = 1). It operates on runes (Unicode code
points), not bytes.
