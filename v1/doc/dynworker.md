<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# dynworker -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/dynworker"
```

## Usage

```go
var (
	ErrNotTask error = errors.Base("task pool entity is not a task")
)
```

#### func  InitPrometheus

```go
func InitPrometheus()
```
Initialise Prometheus metrics for this module.

#### type Config

```go
type Config struct {
	Name        string          // Worker pool name for logger and metrics.
	MinWorkers  int64           // Minimum number of workers.
	MaxWorkers  int64           // Maximum number of workers.
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
	minw, maxw int64,
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
type Task struct {
}
```

Task structure.

#### func  NewTask

```go
func NewTask(ctx context.Context, lgr logger.Logger, data UserData) *Task
```

#### func (*Task) Data

```go
func (obj *Task) Data() UserData
```
Get the user-supplied data for the task.

#### func (*Task) Logger

```go
func (obj *Task) Logger() logger.Logger
```
Get the logger instance for the task.

#### func (*Task) Parent

```go
func (obj *Task) Parent() context.Context
```
Get the parent context for the task.

#### type TaskFn

```go
type TaskFn func(*Task) error
```

Task callback function type.

#### type UserData

```go
type UserData any
```

User-supplied data.

#### type WorkerPool

```go
type WorkerPool interface {
	// Start the worker pool.
	Start()

	// Stop the worker pool.
	Stop()

	// Submit a task to the worker pool.
	Submit(UserData) error

	// Return the number of current workers in the pool.
	WorkerCount() int64

	// Return the minimum number of workers in the pool.
	MinWorkers() int64

	// Return the maximum number of workers in the pool.
	MaxWorkers() int64

	// Set the minimum number of workers to the given value.
	SetMinWorkers(int64)

	// Set the maximum number of workers to the given value.
	SetMaxWorkers(int64)

	// Set the task callback function.
	SetTaskFunction(TaskFn)
}
```

Worker pool interface.

#### func  NewWorkerPool

```go
func NewWorkerPool(config *Config, workfn TaskFn) WorkerPool
```
Create a new worker pool.
