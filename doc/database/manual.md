<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->
# Database Layer

This package provides a thin abstraction over sqlx to make working with
databases in Go simpler, testable, and driver-aware. It is not an ORM -- you
still write SQL, but you get nice ergonomics for:

* clean transaction handling,
* consistent error normalisation (MySQL + Postgres),
* placeholder rebinding (`?` -> `$1`),
* mockability via the Runner interface.

## Opening a connection

``` go
import "your/module/database"

db, err := database.Open("postgres", "postgres://user:pass@localhost/dbname")
if err != nil {
    log.Fatal(err)
}

defer db.Close()

// Ping to verify the connection is live
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

### Supported driver strings:
* "mysql"
* "postgres"
* "pgx" / "pgx/v5"

## The `Database` object

`Database` represents a live connection pool. It exposes:

* `Ping()`, `Close()`, `SetMaxIdleConns(n)`, `SetMaxOpenConns(n)`
* `WithTransaction(ctx, fn)`
* `Rebind(query)` -> rewrites placeholders for the active driver.
* `GetError(err)` -> normalizes driver-specific errors.

## The `Runner` interface

A Runner is anything that can run queries:

* `*sqlx.DB` (no transaction)
* `*sqlx.Tx` (inside a transaction)

It provides `GetContext`, `SelectContext`, `ExecContext`, `QueryxContext`,
etc.

That means repository code can be written against Runner and is agnostic to
whether it’s called inside or outside a transaction.

## Running Queries

### Simple `select`

``` go
var users []User
q := db.Rebind(`SELECT id, name FROM users WHERE active = ?`)
err := db.Real.SelectContext(ctx, &users, q)
if err != nil {
    return db.GetError(err)
}
```

### Insert (MySQL)

``` go
res, err := db.Real.ExecContext(ctx,
    db.Rebind(`INSERT INTO users(name, email) VALUES(?, ?)`),
    "Ada", "ada@example.org",
)
if err != nil {
    return db.GetError(err)
}

id, _ := res.LastInsertId()
```

### Insert (PostreSQL with `RETURNING`)

``` go
var id int64
q := `INSERT INTO users(name, email) VALUES($1, $2) RETURNING id`
err := db.Real.GetContext(ctx, &id, q, "Alan", "alan@example.org")
if err != nil {
    return db.GetError(err)
}
```

## Transactions

Use `WithTransaction` for all multi-step changes. It will:

* begin a transaction,
* run your callback with a `Runner` bound to the `*sqlx.Tx`,
* rollback on error or panic,
* commit on success,
* retry automatically on deadlocks/serialization failures.

``` go
err := db.WithTransaction(ctx, func(ctx context.Context, r database.Runner) error {
    // debit
    q := db.Rebind(`UPDATE accounts SET balance = balance - ? WHERE id = ?`)
    if _, err := r.ExecContext(ctx, q, 100, srcID); err != nil {
        return err
    }

    // credit
    q = db.Rebind(`UPDATE accounts SET balance = balance + ? WHERE id = ?`)
    if _, err := r.ExecContext(ctx, q, 100, dstID); err != nil {
        return err
    }

    return nil
})
if err != nil {
    // err has been normalized (e.g., ErrTxnDeadlock, ErrLostConn)
    return err
}
```

## Error normalisation

`GetError` converts driver-specific errors into well-known conditions:

* MySQL:
  * 1213 → `ErrTxnDeadlock`
  * 2006 → `ErrServerConnClosed`
  * 2013 → `ErrLostConn`
* Postgres:
  * 40P01 → `ErrTxnDeadlock`
  * 40001 → `ErrTxnSerialization`
  * 0800x → `ErrLostConn`
  * 57P01 → `ErrServerConnClosed`

You can `errors.Is(err, database.ErrTxnDeadlock)` to detect and handle them.

## Testing with `sqlmock`

Because everything speaks `Runner`, you can inject a mock:

``` go
db, mock, _ := sqlmock.New()
sqlxDB := sqlx.NewDb(db, "mysql")
d := &database.Database{Real: sqlxDB, Driver: "mysql"}

mock.ExpectBegin()
mock.ExpectExec("INSERT INTO users").WithArgs("Ada").
    WillReturnResult(sqlmock.NewResult(1, 1))
mock.ExpectCommit()

err := d.WithTransaction(ctx, func(ctx context.Context, r database.Runner) error {
    _, err := r.ExecContext(ctx, "INSERT INTO users(name) VALUES(?)", "Ada")
    return err
})
```

## Philosophy

This layer is deliberately thin:

* SQL stays visible. You keep control over queries and performance.
* Transactions are explicit. No hidden context keys.
* Errors are normalized. Your app doesn’t need to know driver codes.
* Testing is easy. `Runner` can be backed by `sqlmock`.

That’s it. Not an ORM, just a safe, ergonomic way to run SQL.
