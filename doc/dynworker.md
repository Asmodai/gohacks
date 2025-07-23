-*- Mode: gfm -*-

# dynworker -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/dynworker"
```

## Usage

```go
const (
	// Default minimum worker count.
	DefaultMinimumWorkerCount int32 = 1

	// Default maximum worker count.
	DefaultMaximumWorkerCount int32 = 10

	// Default worker count multipler.
	//
	// This is used when there is an invalid maximum worker count.
	DefaultWorkerCountMult int32 = 4

	// Default worker timeout.
	DefaultTimeout time.Duration = 30 * time.Second
)
```

#### type Config

```go
type Config struct {
	Name        string          // Worker pool name for logger and metrics.
	MinWorkers  int32           // Minimum number of workers.
	MaxWorkers  int32           // Maximum number of workers.
	Logger      logger.Logger   // Logger instance.
	Parent      context.Context // Parent context.
	IdleTimeout time.Duration   // Idle timeout duration.
}
```


#### func  NewConfig

```go
func NewConfig(
	ctx context.Context,
	lgr logger.Logger,
	name string,
	minw, maxw int32,
) *Config
```
Create a new configuration.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a new default configuration.

#### func (*Config) SetItleTimeout

```go
func (obj *Config) SetItleTimeout(timeout time.Duration)
```
Set the idle timeout value.

#### type Task

```go
type Task any
```

Task data type.

#### type TaskFn

```go
type TaskFn func(Task) error
```

Type of functions executed by workers.

#### type WorkerPool

```go
type WorkerPool interface {
	// Start the worker pool.
	Start()

	// Stop the worker pool.
	Stop()

	// Submit a task to the worker pool.
	Submit(Task) error

	// Return the number of current workers in the pool.
	WorkerCount() int32

	// Return the minimum number of workers in the pool.
	MinWorkers() int32

	// Return the maximum number of workers in the pool.
	MaxWorkers() int32
}
```

Worker pool interface.

#### func  NewWorkerPool

```go
func NewWorkerPool(config *Config, workfn TaskFn) WorkerPool
```
Create a new worker pool.
