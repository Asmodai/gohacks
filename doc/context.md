-*- Mode: gfm -*-

# context -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/context"
```

## Usage

```go
var (
	ErrInvalidValueMap  = errors.Base("invalid value map")
	ErrValueMapNotFound = errors.Base("value map not found")
)
```

#### func  WithValueMap

```go
func WithValueMap(ctx context.Context, valuemap ValueMap) context.Context
```
Create a context with the value map using a default key.

#### func  WithValueMapWithKey

```go
func WithValueMapWithKey(ctx context.Context, key string, valuemap ValueMap) context.Context
```
Create a context with the value map using the specified key.

#### type MockValueMap

```go
type MockValueMap struct {
}
```

MockValueMap is a mock of ValueMap interface.

#### func  NewMockValueMap

```go
func NewMockValueMap(ctrl *gomock.Controller) *MockValueMap
```
NewMockValueMap creates a new mock instance.

#### func (*MockValueMap) EXPECT

```go
func (m *MockValueMap) EXPECT() *MockValueMapMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockValueMap) Get

```go
func (m *MockValueMap) Get(arg0 string) (any, bool)
```
Get mocks base method.

#### func (*MockValueMap) Set

```go
func (m *MockValueMap) Set(key string, value any)
```
Set mocks base method.

#### type MockValueMapMockRecorder

```go
type MockValueMapMockRecorder struct {
}
```

MockValueMapMockRecorder is the mock recorder for MockValueMap.

#### func (*MockValueMapMockRecorder) Get

```go
func (mr *MockValueMapMockRecorder) Get(arg0 any) *gomock.Call
```
Get indicates an expected call of Get.

#### func (*MockValueMapMockRecorder) Set

```go
func (mr *MockValueMapMockRecorder) Set(key, value any) *gomock.Call
```
Set indicates an expected call of Set.

#### type ValueMap

```go
type ValueMap interface {
	Get(string) (key any, ok bool)
	Set(key string, value any)
}
```

A map-based storage structure to pass multiple values via contexts rather than
many invocations of `context.WithValue` and their respective copy operations.

The main caveat with this approach is that as contexts are copied by the various
`With` functions we have no means of passing changes to child contexts once the
context with the value map is copied.

This is not the main aim of this type, so such functionality should not be
considered. The main usage is to provide a means of passing a lot of values to
some top-level context in order to avoid a lot of `WithValue` calls and a
somewhat slow lookup.

#### func  GetValueMap

```go
func GetValueMap(ctx context.Context) (ValueMap, error)
```
Get the value map (if any) from the context.

Returns nil if there is no value map.

#### func  GetValueMapWithKey

```go
func GetValueMapWithKey(ctx context.Context, key string) (ValueMap, error)
```
Get the value map (if any) from the context with the specified value key.

#### func  NewValueMap

```go
func NewValueMap() ValueMap
```
Create a new value map with no data.

#### type ValueMapKey

```go
type ValueMapKey string
```

ValueMap key type for `WithValue`.
