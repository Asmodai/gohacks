<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# database -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/database"
```

## Usage

```go
const ContextKeyManager = "gohacks/database@v1"
```
Key used to store the instance in the context's user value.

```go
var (
	// Signalled when there is no SQL driver.
	ErrNoDriver error = errors.Base("no driver provided")

	// Signalled when there is no SQL username given.
	ErrNoUsername error = errors.Base("no username provided")

	// Signalled when there is no SQL password given.
	ErrNoPassword error = errors.Base("no password provided")

	// Signalled when there is no SQL server hostname given.
	ErrNoHostname error = errors.Base("no hostname provided")

	// Signalled when there is no SQL database name provided.
	ErrNoDatabase error = errors.Base("no database name provided")
)
```

```go
var (
	ErrServerConnClosed = errors.Base("server connection closed")
	ErrLostConn         = errors.Base("lost connection during query")
	ErrTxnDeadlock      = errors.Base("deadlock found when trying to get lock")
	ErrTxnSerialization = errors.Base("serialization failure")
)
```

```go
var (
	ErrNotAWorkerPool = errors.Base("not configured for worker pools")
	ErrNoPoolWorker   = errors.Base("no pool worker function provided")
)
```

```go
var (
	// The empty cursor.
	EmptyCursor = &Cursor{Offset: 0, Limit: 0}
)
```

```go
var ErrValueNotManager = errors.Base("value is not Manager")
```
Signalled if the instance associated with the context key is not of type
Manager.

#### func  SetManager

```go
func SetManager(ctx context.Context, inst Manager) (context.Context, error)
```
Set Manager stores the instance in the context map.

#### func  SetManagerIfAbsent

```go
func SetManagerIfAbsent(ctx context.Context, inst Manager) (context.Context, error)
```
SetManagerIfAbsent sets only if not already present.

#### func  WithManager

```go
func WithManager(ctx context.Context, fn func(Manager))
```
WithManager calls fn with the instance or fallback.

#### type BatchJob

```go
type BatchJob interface {
	Run(ctx context.Context, runner Runner, data []dynworker.UserData) error
}
```

BatchJob provides a means to invoke a user-supplied function with a batch of
jobs.

#### type Config

```go
type Config struct {
	Driver           string         `json:"driver"`
	Username         string         `json:"username"`
	Password         string         `config_obscure:"true"     json:"password"`
	Hostname         string         `json:"hostname"`
	Port             int            `json:"port"`
	Database         string         `json:"database"`
	BatchSize        int            `json:"batch_size"`
	BatchTimeout     types.Duration `json:"batch_timeout"`
	SetPoolLimits    bool           `json:"set_pool_limits"`
	MaxIdleConns     int            `json:"max_idle_conns"`
	MaxOpenConns     int            `json:"max_open_conns"`
	UsePool          bool           `json:"use_worker_pool"`
	PoolMinWorkers   int            `json:"pool_min_workers"`
	PoolMaxWorkers   int            `json:"pool_max_workers"`
	PoolIdleTimeout  types.Duration `json:"pool_idle_timeout"`
	PoolDrainTimeout types.Duration `json:"pool_drain_timeout"`
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
Validate the configuration.

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

	// Rebind query placeholders to the chosen SQL backend.
	Rebind(string) string

	// Parses the given error looking for common MySQL error conditions.
	//
	// If one is found, then a Golang error describing the condition is
	// raised.
	//
	// If nothing interesting is found, then the original error is
	// returned.
	GetError(error) error

	// Run a query function within the context of a database transaction.
	//
	// If there is no error, then the transaction is committed.
	//
	// If there is an error, then the transaction is rolled back.
	WithTransaction(context.Context, TxnFn) error

	// Exposes the database's pool as a `Runner`.
	Runner() Runner
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
	OpenWorker(context.Context, *Config, BatchJob) (Worker, error)
	CheckDB(Database) error
}
```

Database management.

This is a series of wrappers around Go's internal DB stuff to ensure that we set
up max idle/open connections et al.

#### func  FromManager

```go
func FromManager(ctx context.Context) Manager
```
FromManager returns the instance or the fallback.

#### func  GetManager

```go
func GetManager(ctx context.Context) (Manager, error)
```
Get the instance from the given context.

Will return ErrValueNotManager if the value in the context is not of type
Manager.

#### func  MustGetManager

```go
func MustGetManager(ctx context.Context) Manager
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  NewManager

```go
func NewManager() Manager
```
Create a new manager.

#### func  TryGetManager

```go
func TryGetManager(ctx context.Context) (Manager, bool)
```
TryGetManager returns the instance and true if present and typed.

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

#### type Runner

```go
type Runner interface {
	sqlx.ExtContext

	GetContext(context.Context, any, string, ...any) error
	SelectContext(context.Context, any, string, ...any) error
}
```

Runner is "anything that can run sqlx queries with context". Both *sqlx.DB and
*sqlx.Tx satisfy this.

#### type TxnFn

```go
type TxnFn func(context.Context, Runner) error
```


#### type TxnProvider

```go
type TxnProvider interface {
	Txn(context.Context, Runner) error
}
```

Any object that contains a `Txn` function can be used for callbacks.

#### type Worker

```go
type Worker interface {
	Name() string
	Database() Database
	Start()
	Stop()
	SubmitBatch(dynworker.UserData) error
	SubmitJob(WorkerJob) error
}
```


#### func  NewWorker

```go
func NewWorker(parent context.Context, cfg *Config, dbase Database, handler BatchJob) Worker
```

#### type WorkerJob

```go
type WorkerJob interface {
	Run(ctx context.Context, runner Runner) error
}
```

WorkerJob is a user-supplied unit of work which will be executed inside a
database transaction.

Implement `Run' with your SQL using the provided `Runner' (`*sqlx.DB` or
`*sqlx.Tx`).
