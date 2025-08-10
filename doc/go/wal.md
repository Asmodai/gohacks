<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# wal -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/wal"
```

## Usage

```go
var (
	ErrKeyTooLarge       = errors.Base("key too large")
	ErrValueTooLarge     = errors.Base("value too large")
	ErrRecordTooLarge    = errors.Base("record too large")
	ErrShortWrite        = errors.Base("short write")
	ErrInvalidHeader     = errors.Base("invalid WAL header")
	ErrInvalidLog        = errors.Base("WAL log file is invalid")
	ErrTimestampNegative = errors.Base("negative value timestamp")
	ErrTimestampTooBig   = errors.Base("timestamp too big")
)
```

#### type ApplyCallbackFn

```go
type ApplyCallbackFn func(lsn uint64, tstamp int64, key, value []byte) error
```


#### type Policy

```go
type Policy struct {
	// SyncEveryBytes: fsync after this many bytes appended since the
	// last sync.
	//
	// 0 disables byte-based syncing (only time-based or explicit
	// sync/reset/close).
	SyncEveryBytes int64

	// SyncEvery: also fsync on this cadence if there is unsynced data.
	//
	// 0 disables time-based syncing.
	SyncEvery time.Duration

	// Maximum size of a value in bytes.
	MaxValueBytes uint32

	// Maximum size of a key in bytes.
	MaxKeyBytes uint32
}
```


#### type WriteAheadLog

```go
type WriteAheadLog interface {
	// Append writes one record to the log.
	//
	// Errors:
	//   • ErrKeyTooLarge / ErrValueTooLarge / ErrRecordTooLarge when
	//     limits are exceeded.
	//   • ErrTimestampNegative / ErrTimestampTooBig for invalid
	//     timestamps.
	//   • I/O errors wrapped with stack context
	//
	// Durability:
	//   • Subject to Policy; may buffer in page cache until
	//     SyncEvery/SyncEveryBytes hit.
	//   • Call Sync() to force durability.
	Append(lsn uint64, tstamp int64, key, value []byte) error

	// Replay re-applies records with LSN > baseLSN in log order.
	//
	// The supplied callback is invoked for each record; if the callback
	// returns an error, replay aborts and that error is returned. On
	// success, the returned uint64 is the highest LSN applied.
	//
	// Robustness:
	//   • Stops at first malformed/incomplete/CRC-mismatched record
	//     (treats as a safely truncatable tail). Earlier valid records
	//     are preserved.
	Replay(baseLSN uint64, applyCb ApplyCallbackFn) (uint64, error)

	// SetPolicy atomically updates durability and size limits at runtime.
	//
	// Notes:
	//   • Time-based flushing may start/stop a background ticker.
	//   • Limits (MaxKeyBytes/MaxValueBytes) apply to subsequent appends
	//     only.
	SetPolicy(Policy)

	// Sync forces an fsync if there are dirty bytes; otherwise it is a
	// no-op.
	//
	// Returns any I/O error encountered.
	Sync() error

	// Reset truncates the log back to just the header (discarding all
	// records), then fsyncs.
	//
	// Use with care.
	Reset() error

	// Close stops background flushing, flushes if dirty, and closes the
	// file.
	//
	// Returns any flush/close error encountered. Close is
	// idempotent-safe to call once; do not use the instance after Close.
	Close() error
}
```

WriteAheadLog is a single-writer, crash-safe, append-only log.

Records are written in this binary layout:

    [size:u32][lsn:u64][ts:u64][klen:u32][vlen:u32][key][value][crc:u32]

Where `size` is the number of bytes after the size field (including the CRC),
and `lsn` is a caller-supplied logical sequence number, `ts` is a non-negative
UNIX timestamp (seconds), and `crc` is a CRC32C of the payload.

Concurrency & safety:

    - Append/Sync/Reset/SetPolicy/Close are safe to call from multiple
      goroutines, but only one append will make progress at a time.
    - Replay may run concurrently with appends; it snapshots the current
      end and stops at the first incomplete/corrupt record (safe truncation
      behavior).

Durability:

    - Durability is controlled by Policy (time-based fsync via SyncEvery,
      byte-based via SyncEveryBytes). Sync() forces an fsync immediately.
    - Close() stops background flush, flushes if dirty, then closes the file.

Limits & validation:

    - MaxKeyBytes and MaxValueBytes (from Policy) bound record sizes.
    - Timestamps must be in [0, math.MaxInt64]. Violations return errors.

LSN semantics:

    - Caller is responsible for monotonically increasing LSNs.
    - Replay(baseLSN, ...) applies records with LSN > baseLSN and returns the
      highest LSN applied.

Example:

```go

    w, _ := wal.OpenWAL(ctx, "data.wal", 4<<20)
    defer w.Close()
    _ = w.Append(next, time.Now().Unix(), []byte("k"), []byte("v"))
    _ = w.Sync()
    _, _ = w.Replay(0, func(lsn uint64, ts int64, k, v []byte) error {
    	return nil
    })

```

This is probably not the best implementation of a WAL.

#### func  OpenWAL

```go
func OpenWAL(ctx context.Context, path string, syncEveryBytes int64) (WriteAheadLog, error)
```

#### func  OpenWALWithPolicy

```go
func OpenWALWithPolicy(parent context.Context, path string, pol Policy) (WriteAheadLog, error)
```
