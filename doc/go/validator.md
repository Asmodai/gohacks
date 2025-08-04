<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# validator -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/validator"
```

## Usage

#### func  BuildPredicateDict

```go
func BuildPredicateDict() dag.PredicateDict
```

#### func  KindToString

```go
func KindToString(kind FieldKind) string
```

#### type Bindings

```go
type Bindings struct {
	Bindings map[reflect.Type]*StructDescriptor
}
```


#### func  NewBindings

```go
func NewBindings() *Bindings
```

#### func (*Bindings) Bind

```go
func (b *Bindings) Bind(object Reflectable) (*BoundObject, bool)
```

#### func (*Bindings) BindWithReflection

```go
func (b *Bindings) BindWithReflection(object any) (*BoundObject, bool)
```

#### func (*Bindings) Build

```go
func (b *Bindings) Build(object Reflectable) (*StructDescriptor, bool)
```

#### func (*Bindings) BuildWithReflection

```go
func (b *Bindings) BuildWithReflection(object any) (*StructDescriptor, bool)
```

#### func (*Bindings) Register

```go
func (b *Bindings) Register(object *StructDescriptor) bool
```

#### type BoundObject

```go
type BoundObject struct {
	Descriptor *StructDescriptor
	Binding    any
}
```


#### func (*BoundObject) Get

```go
func (bo *BoundObject) Get(key string) (any, bool)
```

#### func (*BoundObject) GetValue

```go
func (bo *BoundObject) GetValue(key string) (any, bool)
```
Get the value for the given key from the bound object.

This works by using the accessor obtained via reflection during the predicate
building phase.

#### func (*BoundObject) Keys

```go
func (bo *BoundObject) Keys() []string
```

#### func (*BoundObject) Set

```go
func (bo *BoundObject) Set(_ string, _ any) bool
```

#### func (*BoundObject) String

```go
func (bo *BoundObject) String() string
```

#### type FTEQBuilder

```go
type FTEQBuilder struct{}
```


#### func (*FTEQBuilder) Build

```go
func (bld *FTEQBuilder) Build(key string, val any) dag.Predicate
```

#### func (*FTEQBuilder) Token

```go
func (bld *FTEQBuilder) Token() string
```

#### type FTEQPredicate

```go
type FTEQPredicate struct {
	MetaPredicate
}
```

Field Type Equality.

This predicate compares the type of the structure's field. If it is equal then
the predicate returns true.

#### func (*FTEQPredicate) Eval

```go
func (pred *FTEQPredicate) Eval(input dag.Filterable) bool
```

#### func (*FTEQPredicate) String

```go
func (pred *FTEQPredicate) String() string
```

#### type FVEQBuilder

```go
type FVEQBuilder struct{}
```


#### func (*FVEQBuilder) Build

```go
func (bld *FVEQBuilder) Build(key string, val any) dag.Predicate
```

#### func (*FVEQBuilder) Token

```go
func (bld *FVEQBuilder) Token() string
```

#### type FVEQPredicate

```go
type FVEQPredicate struct {
	MetaPredicate
}
```

Field Value Equality.

This predicate compares the value to that in the structure. If they are equal
then the predicate returns true.

The predicate will take various circumstances into consideration while checking
the value:

If the field is `any` then the comparison will match just the type of the value
rather than using the type of the field along with the value.

If the field is integer, then the structure's field must have a bid width large
enough to hold the value.

#### func (*FVEQPredicate) Eval

```go
func (pred *FVEQPredicate) Eval(input dag.Filterable) bool
```

#### func (*FVEQPredicate) String

```go
func (pred *FVEQPredicate) String() string
```

#### type FieldAccessorFn

```go
type FieldAccessorFn func(any) any
```


#### type FieldInfo

```go
type FieldInfo struct {
	Name        string
	Type        reflect.Type
	TypeKind    reflect.Kind
	TypeName    string
	Accessor    FieldAccessorFn
	Tags        reflect.StructTag
	Kind        FieldKind
	ElementType reflect.Type
	ElementKind FieldKind
}
```


#### func (*FieldInfo) Debug

```go
func (fi *FieldInfo) Debug(params ...any) *debug.Debug
```

#### func (*FieldInfo) String

```go
func (fi *FieldInfo) String() string
```

#### type FieldKind

```go
type FieldKind int
```


```go
const (
	KindPrimitive FieldKind = iota
	KindStruct
	KindSlice
	KindMap
	KindUnknown
)
```

#### type MetaPredicate

```go
type MetaPredicate struct {
}
```


#### func (*MetaPredicate) GetKeyAsFieldInfo

```go
func (meta *MetaPredicate) GetKeyAsFieldInfo(input dag.Filterable) (*FieldInfo, bool)
```
Return the `Filterable`'s field information.

This is directed through to `BoundObject.Description.Fields`.

#### func (*MetaPredicate) GetKeyAsString

```go
func (meta *MetaPredicate) GetKeyAsString(input dag.Filterable) (string, bool)
```

#### func (*MetaPredicate) GetKeyAsValue

```go
func (meta *MetaPredicate) GetKeyAsValue(input dag.Filterable) (any, bool)
```
Get the value from the `Filterable`.

This equates to `BoundObject.Descriptor.Field[key].Accessor` being called with
`BoundObject.Binding`.

See `BoundObject.GetValue` for more.

#### func (*MetaPredicate) GetValueAsAny

```go
func (meta *MetaPredicate) GetValueAsAny() (any, bool)
```

#### func (*MetaPredicate) GetValueAsBool

```go
func (meta *MetaPredicate) GetValueAsBool() (bool, bool)
```

#### func (*MetaPredicate) GetValueAsComplex128

```go
func (meta *MetaPredicate) GetValueAsComplex128() (complex128, bool)
```

#### func (*MetaPredicate) GetValueAsFloat64

```go
func (meta *MetaPredicate) GetValueAsFloat64() (float64, bool)
```

#### func (*MetaPredicate) GetValueAsInt64

```go
func (meta *MetaPredicate) GetValueAsInt64() (int64, bool)
```

#### func (*MetaPredicate) GetValueAsString

```go
func (meta *MetaPredicate) GetValueAsString() (string, bool)
```
Return the condition value as a string.

#### func (*MetaPredicate) GetValueAsUint64

```go
func (meta *MetaPredicate) GetValueAsUint64() (uint64, bool)
```

#### type Reflectable

```go
type Reflectable interface {
	// Return the reflected type for a given object.
	//
	// An example of how this could work is:
	//
	// ```go
	//     var (
	//         typeForYourStruct reflect.Type
	//         onceForYourStruct sync.Once
	//     )
	//
	//     type YourStruct struct {
	//         // ...
	//     }
	//
	//     func (ys *YourStruct) ReflectType() reflect.Type {
	//         onceForYourStruct.Do(func() {
	//             typeForYourStruct = reflect.TypeOf(ys).Elem()
	//         })
	//
	//         return typeForYourStruct
	//     }
	// ```
	//
	// You might need to do things differently for non-pointer types.
	ReflectType() reflect.Type
}
```


#### type StructDescriptor

```go
type StructDescriptor struct {
	Type     reflect.Type
	TypeName string
	Fields   map[string]*FieldInfo
}
```


#### func  BuildDescriptor

```go
func BuildDescriptor(typ reflect.Type) *StructDescriptor
```

#### func  NewStructDescriptor

```go
func NewStructDescriptor() *StructDescriptor
```

#### func (*StructDescriptor) Debug

```go
func (sd *StructDescriptor) Debug(params ...any) *debug.Debug
```

#### func (*StructDescriptor) Find

```go
func (sd *StructDescriptor) Find(what string) (any, bool)
```

#### func (*StructDescriptor) Get

```go
func (sd *StructDescriptor) Get(key string) (any, bool)
```

#### func (*StructDescriptor) Keys

```go
func (sd *StructDescriptor) Keys() []string
```

#### func (*StructDescriptor) String

```go
func (sd *StructDescriptor) String() string
```
