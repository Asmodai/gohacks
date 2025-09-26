<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/Asmodai/gohacks/database"
    "github.com/Asmodai/gohacks/dynworker"
)

// ---------- domain ----------

type Metric struct {
    HostID     int64
    IP         []byte
    RAMBytes   int64
    CPUPercent float64
    DiskTotal  int64
    DiskUsed   int64
    TS         time.Time
}

// ---------- handler ----------

type MetricsHandler struct {
    DB database.Database
}

// Run is intentionally tiny: delegate to helpers to keep cyclomatic complexity low.
func (h MetricsHandler) Run(ctx context.Context, r database.Runner, batch []dynworker.UserData) error {
    if len(batch) == 0 {
        return nil
    }
    q := h.buildQueries()
    for _, ud := range batch {
        m := ud.(Metric)
        if err := h.handleMetric(ctx, r, q, m); err != nil {
            return err
        }
    }
    return nil
}

// ---- small helpers (each does one thing) ----

type queries struct {
    selLatestTS string
    upsertLatest string
    insertHistory string
}

func (h MetricsHandler) buildQueries() queries {
    const selLatestTS = `
        SELECT ts
        FROM metrics_latest
        WHERE host_id = ?
        FOR UPDATE
    `
    const upsertLatest = `
        INSERT INTO metrics_latest
            (host_id, ip, ram_bytes, cpu_pct, disk_total, disk_used, ts)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            ip = VALUES(ip),
            ram_bytes = VALUES(ram_bytes),
            cpu_pct = VALUES(cpu_pct),
            disk_total = VALUES(disk_total),
            disk_used = VALUES(disk_used),
            ts = VALUES(ts)
    `
    const insertHistory = `
        INSERT INTO metrics_history
            (host_id, ip, ram_bytes, cpu_pct, disk_total, disk_used, ts)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    return queries{
        selLatestTS:  h.DB.Rebind(selLatestTS),
        upsertLatest: h.DB.Rebind(upsertLatest),
        insertHistory:h.DB.Rebind(insertHistory),
    }
}

func (h MetricsHandler) handleMetric(ctx context.Context, r database.Runner, q queries, m Metric) error {
    got, latestTS, err := h.fetchLatestTS(ctx, r, q, m.HostID)
    if err != nil {
        return err
    }

    if h.isNewerOrMissing(got, latestTS, m.TS) {
        if err := h.upsertLatest(ctx, r, q, m); err != nil {
            return err
        }
    }
    return h.insertHistory(ctx, r, q, m)
}

func (h MetricsHandler) fetchLatestTS(ctx context.Context, r database.Runner, q queries, hostID int64) (got bool, ts time.Time, err error) {
    var row struct{ TS time.Time }
    if err = r.GetContext(ctx, &row, q.selLatestTS, hostID); err != nil {
        // Treat “no rows” as missing latest. If your GetContext surfaces sql.ErrNoRows,
        // map it to (got=false, nil error) here; otherwise ignore and return got=false.
        return false, time.Time{}, nil
    }
    return true, row.TS, nil
}

func (h MetricsHandler) isNewerOrMissing(got bool, latestTS, incomingTS time.Time) bool {
    if !got {
        return true
    }
    return incomingTS.After(latestTS)
}

func (h MetricsHandler) upsertLatest(ctx context.Context, r database.Runner, q queries, m Metric) error {
    _, err := r.ExecContext(ctx, q.upsertLatest,
        m.HostID, m.IP, m.RAMBytes, m.CPUPercent, m.DiskTotal, m.DiskUsed, m.TS,
    )
    return err
}

func (h MetricsHandler) insertHistory(ctx context.Context, r database.Runner, q queries, m Metric) error {
    _, err := r.ExecContext(ctx, q.insertHistory,
        m.HostID, m.IP, m.RAMBytes, m.CPUPercent, m.DiskTotal, m.DiskUsed, m.TS,
    )
    return err
}

// ---------- minimal wiring (unchanged from prior example) ----------

func main() {
    ctx := context.Background()

    dsnCfg := &database.Config{
        Driver:   "mysql",
        Username: "user",
        Password: "pass",
        Hostname: "127.0.0.1",
        Port:     3306,
        Database: "metricsdb",
    }
    if errs := dsnCfg.Validate(); len(errs) > 0 {
        panic(fmt.Errorf("bad DSN: %v", errs))
    }
    db, err := database.Open(dsnCfg.Driver, dsnCfg.ToDSN())
    if err != nil {
        panic(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(8)
    db.SetMaxIdleConns(8)

    workerCfg := &database.Config{
        UsePool:          true,
        Database:         dsnCfg.Database,
        PoolMinWorkers:   1,
        PoolMaxWorkers:   8,
        BatchSize:        200,
        BatchTimeout:     database.Duration{100 * time.Millisecond},
        PoolIdleTimeout:  database.Duration{30 * time.Second},
        PoolDrainTimeout: database.Duration{500 * time.Millisecond},
    }
    handler := MetricsHandler{DB: db}
    w := database.NewWorker(ctx, workerCfg, db, handler)
    w.Start()
    defer w.Stop()

    now := time.Now().UTC()
    _ = w.SubmitBatch(Metric{
        HostID: 101, IP: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 10, 0, 0, 1},
        RAMBytes: 16 << 30, CPUPercent: 27.5, DiskTotal: 512 << 30, DiskUsed: 128 << 30, TS: now,
    })
    _ = w.SubmitBatch(Metric{
        HostID: 101, IP: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 10, 0, 0, 1},
        RAMBytes: 16 << 30, CPUPercent: 25.0, DiskTotal: 512 << 30, DiskUsed: 120 << 30, TS: now.Add(-5 * time.Minute),
    })
    _ = w.SubmitBatch(Metric{
        HostID: 101, IP: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 10, 0, 0, 1},
        RAMBytes: 16 << 30, CPUPercent: 30.0, DiskTotal: 512 << 30, DiskUsed: 140 << 30, TS: now.Add(2 * time.Minute),
    })

    time.Sleep(500 * time.Millisecond)
    fmt.Println("submitted metric samples")
}
```
