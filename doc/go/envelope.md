<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# envelope -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/envelope"
```

## Usage

#### func  WriteJSON

```go
func WriteJSON(env Envelope) ([]byte, error)
```

#### type Envelope

```go
type Envelope interface {
	// Return a status code.
	//
	// This could be a HTTP status code or any other integer that makes
	// sense in this context; e.g. Unix error code.
	Status() int

	// Return HTTP headers.
	Headers() http.Header

	// Return the envelope's body.
	Body() any
}
```

Envelope interface.

#### type Error

```go
type Error struct {
	Error   error
	Elapsed time.Duration
}
```


#### func  NewError

```go
func NewError(status int, err error) *Error
```

#### func (*Error) Body

```go
func (ee *Error) Body() any
```

#### func (*Error) Headers

```go
func (ee *Error) Headers() http.Header
```

#### func (*Error) MarshalJSON

```go
func (ee *Error) MarshalJSON() ([]byte, error)
```

#### func (*Error) Status

```go
func (ee *Error) Status() int
```

#### type Success

```go
type Success struct {
	Data    any
	Count   int64
	Elapsed time.Duration
}
```


#### func  NewSuccess

```go
func NewSuccess(status int, data any) *Success
```

#### func (*Success) Body

```go
func (se *Success) Body() any
```

#### func (*Success) Headers

```go
func (se *Success) Headers() http.Header
```

#### func (*Success) MarshalJSON

```go
func (se *Success) MarshalJSON() ([]byte, error)
```

#### func (*Success) Status

```go
func (se *Success) Status() int
```
