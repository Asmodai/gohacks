-*- Mode: gfm -*-

# database -- Asmodai's Go Hacks

```go
    import "github.com/Asmodai/gohacks/database"
```

## Usage

#### type Config

```go
type Config struct {
	Driver        string `json:"driver"`
	Username      string `json:"username"`
	Password      string `json:"password" config_obscure:"true"`
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

#### type Database

```go
type Database struct {
}
```

SQL proxy object.

This trainwreck exists so that we can make use of database interfaces.

It might be 100% useless, as `sql.DB` will most likely conform to `IDatabase`,
so this file might vanish at some point.

#### func (*Database) Begin

```go
func (db *Database) Begin() (*sql.Tx, error)
```

#### func (*Database) Beginx

```go
func (db *Database) Beginx() (*sqlx.Tx, error)
```

#### func (*Database) Close

```go
func (db *Database) Close() error
```

#### func (*Database) Exec

```go
func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error)
```

#### func (*Database) Get

```go
func (db *Database) Get(what interface{}, query string, args ...interface{}) error
```

#### func (*Database) MustBegin

```go
func (db *Database) MustBegin() *sqlx.Tx
```

#### func (*Database) NamedExec

```go
func (db *Database) NamedExec(query string, args interface{}) (sql.Result, error)
```

#### func (*Database) Ping

```go
func (db *Database) Ping() error
```

#### func (*Database) Prepare

```go
func (db *Database) Prepare(query string) (*sql.Stmt, error)
```

#### func (*Database) Query

```go
func (db *Database) Query(query string, args ...interface{}) (IRows, error)
```

#### func (*Database) QueryRowx

```go
func (db *Database) QueryRowx(query string, args ...interface{}) IRow
```

#### func (*Database) Queryx

```go
func (db *Database) Queryx(query string, args ...interface{}) (IRowsx, error)
```

#### func (*Database) Select

```go
func (db *Database) Select(what interface{}, query string, args ...interface{}) error
```

#### func (*Database) SetMaxIdleConns

```go
func (db *Database) SetMaxIdleConns(limit int)
```

#### func (*Database) SetMaxOpenConns

```go
func (db *Database) SetMaxOpenConns(limit int)
```

#### type DatabaseMgr

```go
type DatabaseMgr struct {
}
```

Database management.

This is a series of wrappers around Go's internal DB stuff to ensure that we set
up max idle/open connections et al.

#### func (*DatabaseMgr) CheckDB

```go
func (dbm *DatabaseMgr) CheckDB(db IDatabase) error
```
Check the db connection.

#### func (*DatabaseMgr) Open

```go
func (dbm *DatabaseMgr) Open(driver string, dsn string) (IDatabase, error)
```
Open a connection to the database specified in the DSN string.

#### func (*DatabaseMgr) OpenConfig

```go
func (dbm *DatabaseMgr) OpenConfig(conf *Config) (IDatabase, error)
```
Open and configure a database connection.

#### type IDatabase

```go
type IDatabase interface {
	MustBegin() *sqlx.Tx
	Begin() (*sql.Tx, error)
	Beginx() (*sqlx.Tx, error)
	Close() error
	Exec(string, ...interface{}) (sql.Result, error)
	NamedExec(string, interface{}) (sql.Result, error)
	Ping() error
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (IRows, error)
	Queryx(string, ...interface{}) (IRowsx, error)
	QueryRowx(string, ...interface{}) IRow
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}
```

Interface for `sql.DB` objects.

#### func  Open

```go
func Open(driver string, dsn string) (IDatabase, error)
```

#### type IDatabaseMgr

```go
type IDatabaseMgr interface {
	Open(string, string) (IDatabase, error)
	OpenConfig(*Config) (IDatabase, error)
	CheckDB(IDatabase) error
}
```


#### type IRow

```go
type IRow interface {
	Err() error
	Scan(dest ...interface{}) error
}
```

Interface for `sql.Row` objects.

#### type IRows

```go
type IRows interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...interface{}) error
}
```

Interface for `sql.Rows` objects.

#### type IRowsx

```go
type IRowsx interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...interface{}) error
	StructScan(interface{}) error
}
```

Interface for `sqlx.Rows` objects.

#### type ITx

```go
type ITx interface {
	NamedExec(string, interface{}) (sql.Result, error)
	Commit()
}
```


#### type Tx

```go
type Tx struct {
}
```


#### func (*Tx) Commit

```go
func (tx *Tx) Commit()
```

#### func (*Tx) NamedExec

```go
func (tx *Tx) NamedExec(query string, arg interface{}) (sql.Result, error)
```
