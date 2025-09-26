<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# dynworker -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/dynworker"
```

## Usage

```go
var (
	ErrChannelClosed = errors.Base("channel closed")
)
```

```go
var (
	ErrNotTask error = errors.Base("task pool entity is not a task")
)
```

#### func  InitPrometheus

```go
func InitPrometheus(reg prometheus.Registerer)
```
Initialise Prometheus metrics for this module.

#### type Config

```go
type Config struct {

	// Custom queue to use.
	InputQueue TaskQueue

	// Prometheus registerer.
	Prometheus prometheus.Registerer

	// Function to use as the worker.
	WorkerFunc TaskFn

	// Function to use to determine scaling.
	ScalerFunc ScalerFn

	// Worker pool name for logger and metrics.
	Name string

	// Minimum number of workers.
	MinWorkers int64

	// Maximum number of workers.
	MaxWorkers int64

	// Idle timeout duration.
	IdleTimeout time.Duration

	// Drain target duration.
	DrainTarget time.Duration
}
```


#### func  NewConfig

```go
func NewConfig(name string, minw, maxw int64) *Config
```

#### func  NewConfigWithQueue

```go
func NewConfigWithQueue(name string, minw, maxw int64, queue TaskQueue) *Config
```
Create a new configuration with a custom queue.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a new default configuration.

#### func (*Config) SetIdleTimeout

```go
func (obj *Config) SetIdleTimeout(timeout time.Duration)
```
Set the idle timeout value.

#### func (*Config) SetScalerFunction

```go
func (obj *Config) SetScalerFunction(scalefn ScalerFn)
```
Set the scaler function.

#### func (*Config) SetWorkerFunction

```go
func (obj *Config) SetWorkerFunction(workfn TaskFn)
```
Set the worker function.

#### type ScalerFn

```go
type ScalerFn func() int
```

Scaler callback function type.

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

#### type TaskQueue

```go
type TaskQueue interface {
	Put(context.Context, *Task) error
	Get(context.Context) (*Task, error)
	Len() int
}
```


#### func  NewChanTaskQueue

```go
func NewChanTaskQueue(size int) TaskQueue
```

#### func  NewQueueTaskQueue

```go
func NewQueueTaskQueue(q *types.Queue) TaskQueue
```

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

	// Set the task scaler function.
	SetScalerFunction(ScalerFn)

	// Return the name of the pool.
	Name() string
}
```

Worker pool interface.

#### func  NewWorkerPool

```go
func NewWorkerPool(ctx context.Context, config *Config) WorkerPool
```
Create a new worker pool.

The provided context must have `logger.Logger` in its user value. See
`contextdi` and `logger.SetLogger`.
