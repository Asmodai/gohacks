<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->

``` go
// usersrepo_example.go --- Example repository using Database + Runner.
//
// This file is intended for documentation/examples. It shows how to write a
// repository that:
//   - accepts a `Runner` for every query (either *sqlx.DB or *sqlx.Tx);
//   - keeps SQL visible and explicit;
//   - uses `db.Rebind` to be cross-driver (MySQL '?' vs Postgres '$1');
//   - normalizes driver errors via `db.GetError`;
//   - participates in transactions via `db.WithTransaction`.
//
// Notes:
//   * Replace `your/module/database` with your real module path.
//   * You can split this into separate files in a real project.

package users

import (
    "context"
    "time"

    "gitlab.com/tozd/go/errors"
    "github.com/Asmodai/gohacks/database"
)

// User is the row model for the `users` table.
// Field tags (`db:"..."`) must match column names returned by your queries.
type User struct {
    ID       int64     `db:"id"`
    Name     string    `db:"name"`
    Email    string    `db:"email"`
    Active   bool      `db:"active"`
    Created  time.Time `db:"created_at"`
    Modified time.Time `db:"updated_at"`
}

// Repo provides a thin, explicit SQL layer for the `users` table.
// It holds a handle to your Database so it can call
// Rebind/GetError/WithTransaction.
type Repo struct {
    db database.Database
}

// New constructs a users Repo bound to a Database.
func New(db database.Database) *Repo {
    return &Repo{db: db}
}

// ByID returns a single user by ID.
// The caller supplies a Runner (either the connection pool or a transaction).
func (r *Repo) ByID(ctx context.Context, q database.Runner, id int64) (User, error) {
    var u User

    // Rebind transforms '?' placeholders to the correct style for the active
    // driver.
    //
    // MySQL:  INSERT ... VALUES (?, ?)
    // PG:     INSERT ... VALUES ($1, $2)
    query := r.db.Rebind(`
        SELECT id, name, email, active, created_at, updated_at
        FROM users
        WHERE id = ?`)

    if err := q.GetContext(ctx, &u, query, id); err != nil {
        return User{}, r.db.GetError(err)
    }

    return u, nil
}

// ListActive returns up to `limit` active users, offset by `offset`.
// This shows a typical "get many" using SelectContext into a slice.
func (r *Repo) ListActive(ctx context.Context, q database.Runner, limit, offset int) ([]User, error) {
    var out []User

    query := r.db.Rebind(`
        SELECT id, name, email, active, created_at, updated_at
        FROM users
        WHERE active = TRUE
        ORDER BY id
        LIMIT ? OFFSET ?`)

    if err := q.SelectContext(ctx, &out, query, limit, offset); err != nil {
        return nil, r.db.GetError(err)
    }

    return out, nil
}

// CountActive returns the number of active users (scalar query example).
func (r *Repo) CountActive(ctx context.Context, q database.Runner) (int64, error) {
    query := r.db.Rebind(`SELECT COUNT(*) FROM users WHERE active = TRUE`)

    var n int64

    row := q.QueryRowxContext(ctx, query)

    if err := row.Scan(&n); err != nil {
        return 0, r.db.GetError(err)
    }

    return n, nil
}

// Create inserts a new user (generic style using LastInsertId, ideal for
// MySQL/SQLite).
//
// On Postgres, prefer the CreatePG variant with RETURNING below.
func (r *Repo) Create(ctx context.Context, q database.Runner, name, email string) (int64, error) {
    query := r.db.Rebind(`
        INSERT INTO users (name, email, active, created_at, updated_at)
        VALUES (?, ?, TRUE, NOW(), NOW())`)

    res, err := q.ExecContext(ctx, query, name, email)
    if err != nil {
        return 0, r.db.GetError(err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, errors.WithStack(err)
    }

    return id, nil
}

// CreatePG inserts a new user and returns its ID using RETURNING (preferred
// on Postgres).
//
// This shows that you can also write native SQL for a single driver if
// desired.
func (r *Repo) CreatePG(ctx context.Context, q database.Runner, name, email string) (int64, error) {
    const query = `
        INSERT INTO users (name, email, active, created_at, updated_at)
        VALUES ($1, $2, TRUE, NOW(), NOW())
        RETURNING id`

    var id int64

    if err := q.GetContext(ctx, &id, query, name, email); err != nil {
        return 0, r.db.GetError(err)
    }

    return id, nil
}

// UpdateName updates the user's display name.
// ExecContext returns sql.Result so you can check RowsAffected if you need to.
func (r *Repo) UpdateName(ctx context.Context, q database.Runner, id int64, newName string) error {
    query := r.db.Rebind(`
        UPDATE users
        SET name = ?, updated_at = NOW()
        WHERE id = ?`)

    if _, err := q.ExecContext(ctx, query, newName, id); err != nil {
        return r.db.GetError(err)
    }

    return nil
}

// Delete removes a user by ID. Returns rows affected for convenience.
func (r *Repo) Delete(ctx context.Context, q database.Runner, id int64) (int64, error) {
    query := r.db.Rebind(`DELETE FROM users WHERE id = ?`)

    res, err := q.ExecContext(ctx, query, id)
    if err != nil {
        return 0, r.db.GetError(err)
    }

    aff, err := res.RowsAffected()
    if err != nil {
        return 0, errors.WithStack(err)
    }

    return aff, nil
}

// ---------------------------------------------------------------------------
// Usage examples
// ---------------------------------------------------------------------------
//
// The two examples below show how to call the repository with and without a
// transaction. The repository code is identical in both cases because it
// always speaks to a `Runner`. The *caller* decides whether that runner is
// the pool (*sqlx.DB) or a transaction (*sqlx.Tx).

// Example: calling the repo WITHOUT a transaction.
// Assumes your concrete database exposes a method `Runner() Runner` that
// returns the connection pool; if not, pass the pool from where you
// constructed it.
func Example_noTransaction(ctx context.Context, db database.Database) error {
    repo := New(db)

    // By convention we use the pool (no transaction) as the runner.
    // If your Database doesn't expose Runner(), pass the *sqlx.DB you
    // created.
    runner := db.Runner()

    // Create a user.
    id, err := repo.Create(ctx, runner, "Ada Lovelace", "ada@example.org")
    if err != nil {
        return err
    }

    // Fetch it back.
    _, err = repo.ByID(ctx, runner, id)
    if err != nil {
        return err
    }

    // Count actives.
    _, err = repo.CountActive(ctx, runner)

    return err
}

// Example: calling the repo WITH a transaction.
// `WithTransaction` provides a Runner that is a *sqlx.Tx.
// On error or panic, the transaction will be rolled back; on success,
// committed.
//
// Deadlocks/serialization failures will be retried per your DB layer’s
// policy.
func Example_withTransaction(ctx context.Context, db database.Database) error {
    repo := New(db)

    return db.WithTransaction(ctx, func(ctx context.Context, q database.Runner) error {
        // Insert
        id, err := repo.Create(ctx, q, "Grace Hopper", "grace@example.org")
        if err != nil {
            return err // rollback by WithTransaction
        }

        // Update
        if err := repo.UpdateName(ctx, q, id, "Rear Admiral Grace Hopper"); err != nil {
            return err // rollback by WithTransaction
        }

        // Optional: verify within the same tx
        if _, err := repo.ByID(ctx, q, id); err != nil {
            return err // rollback by WithTransaction
        }

        // No error => commit by WithTransaction
        return nil
    })
}
```
