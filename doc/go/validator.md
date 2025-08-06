<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# validator -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/validator"
```

## Usage

```go
var (
	ErrInvalidSlice       = errors.Base("invalid slice")
	ErrNotCanonicalisable = errors.Base("value cannot be canonicalised")
	ErrNotComparable      = errors.Base("value cannot be compared")
)
```

```go
var (
	ErrValueNotString = errors.Base("value is not a string")
)
```

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
func (bld *FTEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
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
func (pred *FTEQPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FTEQPredicate) String

```go
func (pred *FTEQPredicate) String() string
```

#### type FTINBuilder

```go
type FTINBuilder struct{}
```


#### func (*FTINBuilder) Build

```go
func (bld *FTINBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FTINBuilder) Token

```go
func (bld *FTINBuilder) Token() string
```

#### type FTINPredicate

```go
type FTINPredicate struct {
	MetaPredicate
}
```

Field Type In.

This predicate returns true of the type of a field in the input structure is one
of the provided values in the predicate.

#### func (*FTINPredicate) Eval

```go
func (pred *FTINPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FTINPredicate) String

```go
func (pred *FTINPredicate) String() string
```

#### type FVEQBuilder

```go
type FVEQBuilder struct{}
```


#### func (*FVEQBuilder) Build

```go
func (bld *FVEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
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

If the field is integer, then the structure's field must have a bit width large
enough to hold the value.

#### func (*FVEQPredicate) Eval

```go
func (pred *FVEQPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVEQPredicate) String

```go
func (pred *FVEQPredicate) String() string
```

#### type FVFALSEBuilder

```go
type FVFALSEBuilder struct{}
```


#### func (*FVFALSEBuilder) Build

```go
func (bld *FVFALSEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVFALSEBuilder) Token

```go
func (bld *FVFALSEBuilder) Token() string
```

#### type FVFALSEPredicate

```go
type FVFALSEPredicate struct {
	MetaPredicate
}
```

Field Value is Logically False.

This predicate returns true if the value of the filtered field is logically
false.

A logical false value is any value that is empty or zero. The following are
examples of this:

    "" string, 0 numeric, [] array

Logical falsehood is not the same as `nil`, so if you are looking for nil values
then you should look at `FVNIL` instead.

Structures are a special case. They are never logically false. This is because
the validator does not recurse into structures. If you wish to deal with
structures within structures, then those sub-structures require validation by
themselves. How you do that is up to you.

Interfaces are also a special case. An interface can be considered logically
false if it is `nil`, but it can also be considered logically false if the
wrapped value is zero or empty.

#### func (*FVFALSEPredicate) Eval

```go
func (pred *FVFALSEPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVFALSEPredicate) String

```go
func (pred *FVFALSEPredicate) String() string
```

#### type FVINBuilder

```go
type FVINBuilder struct{}
```


#### func (*FVINBuilder) Build

```go
func (bld *FVINBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVINBuilder) Token

```go
func (bld *FVINBuilder) Token() string
```

#### type FVINPredicate

```go
type FVINPredicate struct {
	MetaPredicate
}
```

Field Value In.

This predicate returns true if the value in the structure is one of the provided
values in the predicate.

#### func (*FVINPredicate) Eval

```go
func (pred *FVINPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVINPredicate) String

```go
func (pred *FVINPredicate) String() string
```

#### type FVNEQBuilder

```go
type FVNEQBuilder struct{}
```


#### func (*FVNEQBuilder) Build

```go
func (bld *FVNEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVNEQBuilder) Token

```go
func (bld *FVNEQBuilder) Token() string
```

#### type FVNEQPredicate

```go
type FVNEQPredicate struct {
	FVEQPredicate
}
```

Field Valie Inequality.

This predicate compares the value to that in the structure. If they are not
equal then the predicate returns true.

The predicate will take various circumstances into consideration while checking
the value:

If the field is `any` then the comparison will match just the type of the value
rather than using the type of the field along with the value.

If the field is integer then the structure's field must have a bit width large
enough to hold the value.

#### func (*FVNEQPredicate) Eval

```go
func (pred *FVNEQPredicate) Eval(ctx context.Context, input dag.Filterable) bool
```

#### func (*FVNEQPredicate) String

```go
func (pred *FVNEQPredicate) String() string
```

#### type FVNILBuilder

```go
type FVNILBuilder struct{}
```


#### func (*FVNILBuilder) Build

```go
func (bld *FVNILBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVNILBuilder) Token

```go
func (bld *FVNILBuilder) Token() string
```

#### type FVNILPredicate

```go
type FVNILPredicate struct {
	MetaPredicate
}
```

Field Value Is Nil.

This predicate returns true if and only if the **reference** value of the
filtered field is `nil`.

If the field value is a concrete type (e.g. string, int, float, bool etc), then
the predicate will return false.

It only applies to types that can be `nil` in Go--e.g. pointers, slices, maps,
interfaces, et al. If you're looking to test whether a field is "logically nil"
(e.g. zero, false, empty) then consider using `FVFALSE` instead.

#### func (*FVNILPredicate) Eval

```go
func (pred *FVNILPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVNILPredicate) String

```go
func (pred *FVNILPredicate) String() string
```

#### type FVTRUEBuilder

```go
type FVTRUEBuilder struct{}
```


#### func (*FVTRUEBuilder) Build

```go
func (bld *FVTRUEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVTRUEBuilder) Token

```go
func (bld *FVTRUEBuilder) Token() string
```

#### type FVTRUEPredicate

```go
type FVTRUEPredicate struct {
	FVFALSEPredicate
}
```

Field Value is Logically True.

This predicate returns true if the value of the filtered field is logically
true.

A logical true value is any value that is not empty or zero.

For more details on how this works, see `FVFALSE`.

#### func (*FVTRUEPredicate) Eval

```go
func (pred *FVTRUEPredicate) Eval(ctx context.Context, input dag.Filterable) bool
```

#### func (*FVTRUEPredicate) String

```go
func (pred *FVTRUEPredicate) String() string
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
	KindInterface
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
