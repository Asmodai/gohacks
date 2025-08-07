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
	ErrInvalidRegexp = errors.Base("invalid regexp")
	ErrRegexpParse   = errors.Base("error parsing regexp")
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

#### func (*FTEQPredicate) Debug

```go
func (pred *FTEQPredicate) Debug() string
```

#### func (*FTEQPredicate) Eval

```go
func (pred *FTEQPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FTEQPredicate) Instruction

```go
func (pred *FTEQPredicate) Instruction() string
```

#### func (*FTEQPredicate) String

```go
func (pred *FTEQPredicate) String() string
```

#### func (*FTEQPredicate) Token

```go
func (pred *FTEQPredicate) Token() string
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

#### func (*FTINPredicate) Debug

```go
func (pred *FTINPredicate) Debug() string
```

#### func (*FTINPredicate) Eval

```go
func (pred *FTINPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FTINPredicate) Instruction

```go
func (pred *FTINPredicate) Instruction() string
```

#### func (*FTINPredicate) String

```go
func (pred *FTINPredicate) String() string
```

#### func (*FTINPredicate) Token

```go
func (pred *FTINPredicate) Token() string
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

#### func (*FVEQPredicate) Debug

```go
func (pred *FVEQPredicate) Debug() string
```

#### func (*FVEQPredicate) Eval

```go
func (pred *FVEQPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVEQPredicate) Instruction

```go
func (pred *FVEQPredicate) Instruction() string
```

#### func (*FVEQPredicate) String

```go
func (pred *FVEQPredicate) String() string
```

#### func (*FVEQPredicate) Token

```go
func (pred *FVEQPredicate) Token() string
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

#### func (*FVFALSEPredicate) Debug

```go
func (pred *FVFALSEPredicate) Debug() string
```

#### func (*FVFALSEPredicate) Eval

```go
func (pred *FVFALSEPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVFALSEPredicate) Instruction

```go
func (pred *FVFALSEPredicate) Instruction() string
```

#### func (*FVFALSEPredicate) String

```go
func (pred *FVFALSEPredicate) String() string
```

#### func (*FVFALSEPredicate) Token

```go
func (pred *FVFALSEPredicate) Token() string
```

#### type FVGTBuilder

```go
type FVGTBuilder struct{}
```


#### func (*FVGTBuilder) Build

```go
func (bld *FVGTBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVGTBuilder) Token

```go
func (bld *FVGTBuilder) Token() string
```

#### type FVGTEBuilder

```go
type FVGTEBuilder struct{}
```


#### func (*FVGTEBuilder) Build

```go
func (bld *FVGTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVGTEBuilder) Token

```go
func (bld *FVGTEBuilder) Token() string
```

#### type FVGTEPredicate

```go
type FVGTEPredicate struct {
	MetaPredicate
}
```


#### func (*FVGTEPredicate) Debug

```go
func (pred *FVGTEPredicate) Debug() string
```

#### func (*FVGTEPredicate) Eval

```go
func (pred *FVGTEPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVGTEPredicate) Instruction

```go
func (pred *FVGTEPredicate) Instruction() string
```

#### func (*FVGTEPredicate) String

```go
func (pred *FVGTEPredicate) String() string
```

#### func (*FVGTEPredicate) Token

```go
func (pred *FVGTEPredicate) Token() string
```

#### type FVGTPredicate

```go
type FVGTPredicate struct {
	MetaPredicate
}
```


#### func (*FVGTPredicate) Debug

```go
func (pred *FVGTPredicate) Debug() string
```

#### func (*FVGTPredicate) Eval

```go
func (pred *FVGTPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVGTPredicate) Instruction

```go
func (pred *FVGTPredicate) Instruction() string
```

#### func (*FVGTPredicate) String

```go
func (pred *FVGTPredicate) String() string
```

#### func (*FVGTPredicate) Token

```go
func (pred *FVGTPredicate) Token() string
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

#### func (*FVINPredicate) Debug

```go
func (pred *FVINPredicate) Debug() string
```

#### func (*FVINPredicate) Eval

```go
func (pred *FVINPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVINPredicate) Instruction

```go
func (pred *FVINPredicate) Instruction() string
```

#### func (*FVINPredicate) String

```go
func (pred *FVINPredicate) String() string
```

#### func (*FVINPredicate) Token

```go
func (pred *FVINPredicate) Token() string
```

#### type FVLTBuilder

```go
type FVLTBuilder struct{}
```


#### func (*FVLTBuilder) Build

```go
func (bld *FVLTBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVLTBuilder) Token

```go
func (bld *FVLTBuilder) Token() string
```

#### type FVLTEBuilder

```go
type FVLTEBuilder struct{}
```


#### func (*FVLTEBuilder) Build

```go
func (bld *FVLTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVLTEBuilder) Token

```go
func (bld *FVLTEBuilder) Token() string
```

#### type FVLTEPredicate

```go
type FVLTEPredicate struct {
	MetaPredicate
}
```


#### func (*FVLTEPredicate) Debug

```go
func (pred *FVLTEPredicate) Debug() string
```

#### func (*FVLTEPredicate) Eval

```go
func (pred *FVLTEPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVLTEPredicate) Instruction

```go
func (pred *FVLTEPredicate) Instruction() string
```

#### func (*FVLTEPredicate) String

```go
func (pred *FVLTEPredicate) String() string
```

#### func (*FVLTEPredicate) Token

```go
func (pred *FVLTEPredicate) Token() string
```

#### type FVLTPredicate

```go
type FVLTPredicate struct {
	MetaPredicate
}
```


#### func (*FVLTPredicate) Debug

```go
func (pred *FVLTPredicate) Debug() string
```

#### func (*FVLTPredicate) Eval

```go
func (pred *FVLTPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVLTPredicate) Instruction

```go
func (pred *FVLTPredicate) Instruction() string
```

#### func (*FVLTPredicate) String

```go
func (pred *FVLTPredicate) String() string
```

#### func (*FVLTPredicate) Token

```go
func (pred *FVLTPredicate) Token() string
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

#### func (*FVNEQPredicate) Debug

```go
func (pred *FVNEQPredicate) Debug() string
```

#### func (*FVNEQPredicate) Eval

```go
func (pred *FVNEQPredicate) Eval(ctx context.Context, input dag.Filterable) bool
```

#### func (*FVNEQPredicate) Instruction

```go
func (pred *FVNEQPredicate) Instruction() string
```

#### func (*FVNEQPredicate) String

```go
func (pred *FVNEQPredicate) String() string
```

#### func (*FVNEQPredicate) Token

```go
func (pred *FVNEQPredicate) Token() string
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

#### func (*FVNILPredicate) Debug

```go
func (pred *FVNILPredicate) Debug() string
```

#### func (*FVNILPredicate) Eval

```go
func (pred *FVNILPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVNILPredicate) Instruction

```go
func (pred *FVNILPredicate) Instruction() string
```

#### func (*FVNILPredicate) String

```go
func (pred *FVNILPredicate) String() string
```

#### func (*FVNILPredicate) Token

```go
func (pred *FVNILPredicate) Token() string
```

#### type FVREMBuilder

```go
type FVREMBuilder struct{}
```


#### func (*FVREMBuilder) Build

```go
func (bld *FVREMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error)
```

#### func (*FVREMBuilder) Token

```go
func (bld *FVREMBuilder) Token() string
```

#### type FVREMPredicate

```go
type FVREMPredicate struct {
	MetaPredicate
}
```


#### func (*FVREMPredicate) Debug

```go
func (pred *FVREMPredicate) Debug() string
```

#### func (*FVREMPredicate) Eval

```go
func (pred *FVREMPredicate) Eval(_ context.Context, input dag.Filterable) bool
```

#### func (*FVREMPredicate) Instruction

```go
func (pred *FVREMPredicate) Instruction() string
```

#### func (*FVREMPredicate) String

```go
func (pred *FVREMPredicate) String() string
```

#### func (*FVREMPredicate) Token

```go
func (pred *FVREMPredicate) Token() string
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

#### func (*FVTRUEPredicate) Debug

```go
func (pred *FVTRUEPredicate) Debug() string
```

#### func (*FVTRUEPredicate) Eval

```go
func (pred *FVTRUEPredicate) Eval(ctx context.Context, input dag.Filterable) bool
```

#### func (*FVTRUEPredicate) Instruction

```go
func (pred *FVTRUEPredicate) Instruction() string
```

#### func (*FVTRUEPredicate) String

```go
func (pred *FVTRUEPredicate) String() string
```

#### func (*FVTRUEPredicate) Token

```go
func (pred *FVTRUEPredicate) Token() string
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


#### func (*MetaPredicate) Debug

```go
func (meta *MetaPredicate) Debug(isn, token string) string
```

#### func (*MetaPredicate) GetKeyAsFieldInfo

```go
func (meta *MetaPredicate) GetKeyAsFieldInfo(input dag.Filterable) (*FieldInfo, bool)
```
Return the `Filterable`'s field information.

This is directed through to `BoundObject.Description.Fields`.

#### func (*MetaPredicate) GetKeyAsFloat64

```go
func (meta *MetaPredicate) GetKeyAsFloat64(input dag.Filterable) (float64, bool)
```

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

#### func (*MetaPredicate) String

```go
func (meta *MetaPredicate) String(token string) string
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
