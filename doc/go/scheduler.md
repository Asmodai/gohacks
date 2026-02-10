<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# scheduler -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/scheduler"
```

## Usage

```go
const (
	// Default health ticker period.
	DefaultHealthTickPeriod time.Duration = 1 * time.Second

	// Default channel buffer size.
	//
	// This should be sufficiently big enough to prevent blocking.
	DefaultChannelBufferSize int = 256
)
```

```go
const (
	// Number of seconds before a task is considered late.
	LateTaskDelay time.Duration = 5 * time.Second
)
```

```go
var (
	ErrJobRunAtZero       = errors.Base("job run-at is zero")
	ErrJobNoTarget        = errors.Base("no job target specified")
	ErrJobAmbiguousTarget = errors.Base("job target is ambiguous")
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
	Name             string                // Scheduler name.
	HealthTickPeriod time.Duration         // Health tick period.
	Health           health.Reporter       // Health status reporter.
	Prometheus       prometheus.Registerer // Prometheus registerer.
	AddBuffer        int                   // Size of `add` buffer.
	WorkBuffer       int                   // Size of `work` buffer.
}
```

Scheduler configuration.

#### func  NewConfig

```go
func NewConfig(
	name string,
	tick time.Duration,
	hlth health.Reporter,
	addBuff, workBuff int,
) *Config
```
Create a new scheduler configuration instance.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a new default scheduler configuration instance.

#### type Job

```go
type Job interface {
	Validate() error
	Resolve(context.Context) error

	Object() Task
	Function() JobFn
}
```


#### func  InsertJob

```go
func InsertJob(jobs []Job, njob Job) ([]Job, error)
```
Insert a job into a list of jobs.

#### func  MakeJob

```go
func MakeJob(obj Task, fn JobFn) Job
```
Create a new job.

#### type JobFn

```go
type JobFn func(context.Context) error
```

Job function.

#### type Priority

```go
type Priority struct {
}
```

Priority scheduler

Instance is single-use; create a new instance to restart.

#### func  NewPriority

```go
func NewPriority(ctx context.Context, cnf *Config) *Priority
```
Return a new priority scheduler instance.

#### func (*Priority) Done

```go
func (s *Priority) Done() <-chan struct{}
```
Has the priority scheduler done processing?

#### func (*Priority) Health

```go
func (s *Priority) Health() health.Reporter
```
Return the health reporter for the priority scheduler.

#### func (*Priority) Name

```go
func (s *Priority) Name() string
```
Return the name of the priority scheduler.

#### func (*Priority) Next

```go
func (s *Priority) Next(ctx context.Context) (TimedJob, bool)
```
Get the next task in the work channel.

#### func (*Priority) Start

```go
func (s *Priority) Start()
```
Start the scheduler.

Must be called once before use.

#### func (*Priority) Stop

```go
func (s *Priority) Stop()
```
Stop the priority scheduler.

#### func (*Priority) Submit

```go
func (s *Priority) Submit(task TimedJob) error
```
Add a task to the priority scheduler.

#### func (*Priority) Wait

```go
func (s *Priority) Wait(ctx context.Context) error
```
Wait for the goroutine created by `Start`.

#### func (*Priority) Work

```go
func (s *Priority) Work() <-chan TimedJob
```
Get the current work channel for the priority scheduler.

#### type Task

```go
type Task interface {
	Execute(context.Context) error
}
```


#### type TimedJob

```go
type TimedJob interface {
	Job

	RunAt() time.Time
}
```


#### func  InsertTimedJob

```go
func InsertTimedJob(jobs []TimedJob, njob TimedJob) ([]TimedJob, error)
```
Insert a timed job into a list of jobs.

#### func  MakeTimedJob

```go
func MakeTimedJob(runAt time.Time, obj Task, fn JobFn) TimedJob
```
Create a new timed job.
