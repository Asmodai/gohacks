<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# crypto -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/crypto"
```

## Usage

#### func  GenerateRandomBytes

```go
func GenerateRandomBytes(n int) ([]byte, error)
```
Generate n number of random bytes from a cryptographic randomiser.

#### func  GenerateRandomSafeBytes

```go
func GenerateRandomSafeBytes(count int) ([]byte, error)
```
Operates the same way as `GenerateRandomBytes` but encodes the result using
base64 encoding.

#### func  GenerateRandomSafeString

```go
func GenerateRandomSafeString(count int) (string, error)
```
Operates the same way as `GenerateRandomString` but encodes the result using
base64 encoding.

#### func  GenerateRandomString

```go
func GenerateRandomString(count int) (string, error)
```
Generate a random string of the given length using bytes from the cryptographic
randomiser.
