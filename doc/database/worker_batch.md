<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/Asmodai/gohacks/dynworker"
    "github.com/Asmodai/gohacks/database"
)

// Per-item payload you submit via SubmitBatch.
type UserRow struct {
    ID   int64
    Name string
}

// Batch handler: receives a slice of UserData (your payloads) and runs one tx.
type UserInsertHandler struct {
    DB database.Database // for Rebind
}

func (h UserInsertHandler) Run(ctx context.Context, r database.Runner, data []dynworker.UserData) error {
    if len(data) == 0 {
        return nil
    }

    // Build a multi-VALUES UPSERT for the whole batch.
    const base = `INSERT INTO users (id, name) VALUES `
    sql := base
    args := make([]any, 0, len(data)*2)

    for i, ud := range data {
        row := ud.(UserRow) // your own type
        if i > 0 {
            sql += ", "
        }
        sql += "(?, ?)"
        args = append(args, row.ID, row.Name)
    }

    sql += " ON DUPLICATE KEY UPDATE name = VALUES(name)"

    sql = h.DB.Rebind(sql)
    _, err := r.ExecContext(ctx, sql, args...)

    return err
}

func main() {
    ctx := context.Background()

    // 1) Open DB (fill in your real credentials)
    dsnCfg := &database.Config{
        Driver:   "mysql",
        Username: "user",
        Password: "pass",
        Hostname: "127.0.0.1",
        Port:     3306,
        Database: "app",
    }
    if errs := dsnCfg.Validate(); len(errs) != nil && len(errs) > 0 {
        panic(fmt.Errorf("bad DSN config: %v", errs))
    }

    db, err := database.Open(dsnCfg.Driver, dsnCfg.ToDSN())
    if err != nil {
        panic(err)
    }

    defer db.Close()
    db.SetMaxOpenConns(8)
    db.SetMaxIdleConns(8)

    // 2) Build a DB worker (pool) — batch mode uses SubmitBatch
    workerCfg := &database.Config{ // worker/pool config
        UsePool:          true,
        Database:         dsnCfg.Database,
        PoolMinWorkers:   1,
        PoolMaxWorkers:   8,                    // ≤ MaxOpenConns
        BatchSize:        200,                  // rows per transaction
        BatchTimeout:     database.Duration{100 * time.Millisecond}, // flush sparse traffic
        PoolIdleTimeout:  database.Duration{30 * time.Second},
        PoolDrainTimeout: database.Duration{500 * time.Millisecond},
    }
    handler := UserInsertHandler{DB: db}

    w := database.NewWorker(ctx, workerCfg, db, handler)
    w.Start()
    defer w.Stop()

    // 3) Submit many per-item payloads; batcher will group them
    _ = w.SubmitBatch(UserRow{ID: 10, Name: "Marie"})
    _ = w.SubmitBatch(UserRow{ID: 11, Name: "Rosalind"})
    _ = w.SubmitBatch(UserRow{ID: 12, Name: "Katherine"})

    // Let the pool process (for demo)
    time.Sleep(500 * time.Millisecond)
    fmt.Println("submitted 3 batched items (will be flushed by timeout or size)")
}
```
