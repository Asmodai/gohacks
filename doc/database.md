-*- Mode: gfm -*-

# database -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/database"
```

## Usage

```go
var (
	ErrNoDriver   error = errors.Base("no driver provided")
	ErrNoUsername error = errors.Base("no username provided")
	ErrNoPassword error = errors.Base("no password provided")
	ErrNoHostname error = errors.Base("no hostname provided")
	ErrNoDatabase error = errors.Base("no database name provided")
)
```

```go
var (
	ErrNoContextKey       error = errors.Base("no context key given")
	ErrValueIsNotDatabase error = errors.Base("not a database")
)
```

```go
var (
	// The empty cursor.
	EmptyCursor = &Cursor{Offset: 0, Limit: 0}
)
```

#### func  ToContext

```go
func ToContext(ctx context.Context, inst Database, key string) (context.Context, error)
```

#### type Config

```go
type Config struct {
	Driver        string `json:"driver"`
	Username      string `json:"username"`
	Password      string `config_obscure:"true"  json:"password"`
	Hostname      string `json:"hostname"`
	Port          int    `json:"port"`
	Database      string `json:"database"`
	BatchSize     int    `json:"batch_size"`
	SetPoolLimits bool   `json:"set_pool_limits"`
	MaxIdleConns  int    `json:"max_idle_conns"`
	MaxOpenConns  int    `json:"max_open_conns"`
}
```

SQL configuration structure.

#### func  NewConfig

```go
func NewConfig() *Config
```
Create a new configuration object.

#### func (*Config) ToDSN

```go
func (c *Config) ToDSN() string
```
Return the DSN for this database configuration.

#### func (*Config) Validate

```go
func (c *Config) Validate() error
```

#### type Cursor

```go
type Cursor struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}
```

TODO: Look, this sucks... offset/limit cursors are just fail. TODO: Rework this
to be a proper cursor!

#### func  NewCursor

```go
func NewCursor(offset, limit int64) *Cursor
```
Create a new cursor.

#### func (Cursor) Valid

```go
func (c Cursor) Valid() bool
```
Is the cursor valid?

#### type Database

```go
type Database interface {
}
```


#### func  FromContext

```go
func FromContext(ctx context.Context, key string) (Database, error)
```

#### func  FromDB

```go
func FromDB(db *sql.DB, driver string) Database
```

#### func  Open

```go
func Open(driver string, dsn string) (Database, error)
```
Open a connection using the relevant driver to the given data source name.

#### type Manager

```go
type Manager interface {
	Open(string, string) (Database, error)
	OpenConfig(*Config) (Database, error)
	CheckDB(Database) error
}
```

Database management.

This is a series of wrappers around Go's internal DB stuff to ensure that we set
up max idle/open connections et al.

#### type NullBool

```go
type NullBool struct {
	sql.NullBool
}
```


#### func (NullBool) MarshalJSON

```go
func (x NullBool) MarshalJSON() ([]byte, error)
```

#### type NullByte

```go
type NullByte struct {
	sql.NullByte
}
```


#### func (NullByte) MarshalJSON

```go
func (x NullByte) MarshalJSON() ([]byte, error)
```

#### type NullFloat64

```go
type NullFloat64 struct {
	sql.NullFloat64
}
```


#### func (NullFloat64) MarshalJSON

```go
func (x NullFloat64) MarshalJSON() ([]byte, error)
```

#### type NullInt16

```go
type NullInt16 struct {
	sql.NullInt16
}
```


#### func (NullInt16) MarshalJSON

```go
func (x NullInt16) MarshalJSON() ([]byte, error)
```

#### type NullInt32

```go
type NullInt32 struct {
	sql.NullInt32
}
```


#### func (NullInt32) MarshalJSON

```go
func (x NullInt32) MarshalJSON() ([]byte, error)
```

#### type NullInt64

```go
type NullInt64 struct {
	sql.NullInt64
}
```


#### func (NullInt64) MarshalJSON

```go
func (x NullInt64) MarshalJSON() ([]byte, error)
```

#### type NullString

```go
type NullString struct {
	sql.NullString
}
```


#### func (NullString) MarshalJSON

```go
func (x NullString) MarshalJSON() ([]byte, error)
```

#### type NullTime

```go
type NullTime struct {
	sql.NullTime
}
```


#### func (NullTime) MarshalJSON

```go
func (x NullTime) MarshalJSON() ([]byte, error)
```

#### type Row

```go
type Row interface {
	Err() error
	Scan(...any) error
}
```


#### type Rows

```go
type Rows interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...any) error
}
```


#### type Rowsx

```go
type Rowsx interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...any) error
	StructScan(any) error
}
```


#### type SQL

```go
type SQL struct {
}
```


#### func  NewSQL

```go
func NewSQL(impl implementation) *SQL
```
