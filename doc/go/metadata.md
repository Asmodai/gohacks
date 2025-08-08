<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# metadata -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/metadata"
```

## Usage

```go
const (
	TagDelimiter = ","

	KeyDoc        = "doc"
	KeySince      = "since"
	KeyVersion    = "version"
	KeyDeprecated = "deprecated"
	KeyProtocol   = "protocol"
	KeyVisibility = "visibility"
	KeyExample    = "example"
	KeyTags       = "tags"
	KeyAuthor     = "author"
)
```

```go
var (
	ErrKeyIsInvalid   = errors.Base("metadata key is invalid")
	ErrKeyIsAmbiguous = errors.Base("metadata key is ambiguous")
	ErrKeyIsReserved  = errors.Base("reserved metadata key")

	// Reserved selector metadata keys and their intended semantics.
	//
	// These keys are reserved for internal use and documentation
	// purposes.
	//
	// While not enforced, they are expected to follow consistent
	// formatting and be used by tooling for introspection, generation,
	// and validation.
	//
	// - "doc": Short description of what the selector does. (Markdown
	//          allowed.)
	// - "since": First version or date this selector was introduced.
	//            Format: "v1.2.3" or ISO date (e.g. "2025-08-07").
	// - "version": Current version of the selector logic.
	//              Useful if selector behavior has changed over time.
	// - "deprecated": Optional deprecation notice or replacement advice.
	// - "protocol": Protocol(s) this selector belongs to
	//               (comma-separated).
	// - "visibility": One of: "public", "internal", "private".
	//                 Useful for generating user-facing docs or
	//                 restricting UI tools.
	// - "example": A usage example (inline or structured Markdown).
	// - "tags": Comma-separated labels for filtering or grouping.
	//           E.g. "filesystem,experimental,fastpath"
	// - "author": Who wrote or maintains the selector logic.
	//             Useful for blame or kudos.
	//
	// Tools can recognize and use these for generating CLI docs, debug
	// dumps, live inspector UIs, etc.
	//
	//nolint:gochecknoglobals
	ReservedMetadataKeys = map[string]struct{}{
		KeyDoc:        {},
		KeySince:      {},
		KeyVersion:    {},
		KeyDeprecated: {},
		KeyProtocol:   {},
		KeyVisibility: {},
		KeyExample:    {},
		KeyTags:       {},
		KeyAuthor:     {},
	}

	//nolint:gochecknoglobals
	VisibilityLevels = map[string]struct{}{
		"public":   {},
		"internal": {},
		"private":  {},
	}
)
```

#### type Metadata

```go
type Metadata interface {
	List() map[string]string

	Set(string, string) error

	Get(string) (string, bool)

	SetDoc(string) Metadata
	SetSince(string) Metadata
	SetVersion(string) Metadata
	SetDeprecated(string) Metadata
	SetProtocol(string) Metadata
	SetVisibility(string) Metadata
	SetExample(string) Metadata
	SetTags(string) Metadata
	SetTagsFromSlice([]string) Metadata
	SetAuthor(string) Metadata

	GetDoc() string
	GetSince() string
	GetVersion() string
	GetDeprecated() string
	GetProtocol() string
	GetVisibility() string
	GetExample() string
	GetTags() string
	GetAuthor() string

	Clone() Metadata
	Merge(Metadata, bool) Metadata

	Tags() []string
	TagsNormalised() []string
}
```


#### func  NewMetadata

```go
func NewMetadata() Metadata
```
