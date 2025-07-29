<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# database -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/database"
```

## Usage

```go
const (
	KeyTransaction   string = "_DB_TXN"
	StringDeadlock   string = "Error 1213" // Deadlock detected.
	StringConnClosed string = "Error 2006" // MySQL server connection closed.
	StringLostConn   string = "Error 2013" // Lost connection during query.
)
```

```go
const (
	ContextKeyDBManager = "_DI_DB_MGR"
)
```

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
	ErrTxnKeyNotFound error = errors.Base("transaction key not found")
	ErrTxnKeyNotTxn   error = errors.Base("key value is not a transaction")
	ErrTxnContext     error = errors.Base("could not create transaction context")
	ErrTxnStart       error = errors.Base("could not start transaction")

	ErrTxnDeadlock      error = errors.Base("deadlock found when trying to get lock")
	ErrServerConnClosed error = errors.Base("server connection closed")
	ErrLostConn         error = errors.Base("lost connection during query")
)
```

```go
var (
	// The empty cursor.
	EmptyCursor = &Cursor{Offset: 0, Limit: 0}
)
```

```go
var (
	ErrValueNotDBManager = errors.Base("value is not database.Manager")
)
```

#### func  Exec

```go
func Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
```
Wrapper around `Tx.Exec`.

The transaction should be passed via a context value.

#### func  ExecStmt

```go
func ExecStmt(ctx context.Context, stmt *stmt, args ...any) (sql.Result, error)
```
Wrapper around `Tx.ExecStmt`.

The transaction should be passed via a context value.

#### func  Get

```go
func Get(ctx context.Context, dest any, query string, args ...any) error
```
Wrapper around `Tx.Get`.

The transaction should be passed via a context value.

#### func  NamedExec

```go
func NamedExec(ctx context.Context, query string, arg any) (sql.Result, error)
```
Wrapper around `Tx.NamedExec`.

The transaction should be passed via a context value.

#### func  Prepare

```go
func Prepare(ctx context.Context, query string, args ...any) (*stmt, error)
```
Wrapper around `Tx.Prepare`.

The transaction should be passed via a context value.

#### func  Queryx

```go
func Queryx(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
```
Wrapper around `Tx.Queryx`.

The transaction should be passed via a context value.

#### func  QueryxContext

```go
func QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
```
Wrapper around `Tx.QueryxContext`.

The transaction should be passed via a context value.

#### func  Select

```go
func Select(ctx context.Context, dest any, query string, args ...any) error
```
Wrapper around `Tx.Select`.

The transaction should be passed via a context value.

#### func  SetManager

```go
func SetManager(ctx context.Context, inst Manager) (context.Context, error)
```
Set the database manager value to the context map.

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
func (c *Config) Validate() []error
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
	// Pings the database connection to ensure it is alive and connected.
	Ping() error

	// Close a database connection.  This does nothing if the connection
	// is already closed.
	Close() error

	// Set the maximum idle connections.
	SetMaxIdleConns(int)

	// Set the maximum open connections.
	SetMaxOpenConns(int)

	// Return the transaction (if any) from the given context.
	Tx(context.Context) (*sqlx.Tx, error)

	// Initiate a transaction.  Returns a new context that contains the
	// database transaction session as a value.
	Begin(context.Context) (context.Context, error)

	// Initiate a transaction commit.
	Commit(context.Context) error

	// Initiate a transaction rollback.
	Rollback(context.Context) error

	// Parses the given error looking for common MySQL error conditions.
	//
	// If one is found, then a Golang error describing the condition is
	// raised.
	//
	// If nothing interesting is found, then the original error is
	// returned.
	GetError(error) error
}
```


#### func  FromDB

```go
func FromDB(db *sql.DB, driver string) Database
```
Create a new database object using an existing `sql` object.

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

#### func  GetManager

```go
func GetManager(ctx context.Context) (Manager, error)
```
Get the database manager from the given context.

Will return `ErrValueNoDBManager` if the value in the context is not of type
`database.Manager`.

#### func  MustGetManager

```go
func MustGetManager(ctx context.Context) Manager
```
Attempt to get the database manager from the given context. Panics if the
operation fails.

#### func  NewManager

```go
func NewManager() Manager
```
Create a new manager.

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
