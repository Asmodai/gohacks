<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/Asmodai/gohacks/database"
)

// --- A single-item job (each call is its own transaction) ---

type CreateUserJob struct {
    DB   database.Database // only for Rebind; Runner does ExecContext
    ID   int64
    Name string
}

func (j CreateUserJob) Run(ctx context.Context, r database.Runner) error {
    const q = `
        INSERT INTO users (id, name)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE name = VALUES(name)`
    sql := j.DB.Rebind(q)
    _, err := r.ExecContext(ctx, sql, j.ID, j.Name)
    return err
}

// --- A trivial batch handler (not used in job mode, but required by NewWorker) ---

type NoopBatchHandler struct{ DB database.Database }

func (h NoopBatchHandler) Run(ctx context.Context, r database.Runner, data []any) error {
    // No-op; documentation example focuses on SubmitJob path.
    return nil
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
    db.SetMaxOpenConns(4)
    db.SetMaxIdleConns(4)

    // 2) Build a DB worker (pool) — job mode uses SubmitJob
    workerCfg := &database.Config{ // this is your "worker" config with pool knobs
        UsePool:          true,
        Database:         dsnCfg.Database,
        PoolMinWorkers:   1,
        PoolMaxWorkers:   4,
        BatchSize:        100, // irrelevant for SubmitJob, but required
        BatchTimeout:     database.Duration{100 * time.Millisecond},
        PoolIdleTimeout:  database.Duration{30 * time.Second},
        PoolDrainTimeout: database.Duration{500 * time.Millisecond},
    }

    handler := NoopBatchHandler{DB: db} // not used for SubmitJob, but required

    w := database.NewWorker(ctx, workerCfg, db, handler)
    w.Start()
    defer w.Stop()

    // 3) Submit a few single-item jobs
    _ = w.SubmitJob(CreateUserJob{DB: db, ID: 1, Name: "Ada"})
    _ = w.SubmitJob(CreateUserJob{DB: db, ID: 2, Name: "Grace"})
    _ = w.SubmitJob(CreateUserJob{DB: db, ID: 3, Name: "Hedy"})

    // Let the pool process (for demo)
    time.Sleep(500 * time.Millisecond)
    fmt.Println("submitted 3 single-item jobs")
}
```
