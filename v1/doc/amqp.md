<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# amqp -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/amqp"
```

## Usage

```go
var (
	ErrNoWorkerPool error = errors.Base("no worker pool available")
)
```

#### func  InitPrometheus

```go
func InitPrometheus()
```
Initialise Prometheus metrics for this module.

#### type Client

```go
type Client interface {
	Connect() error
	IsConnected() bool
	Consume() error
	Publish(goamqp.Publishing) error
	QueueStats() (goamqp.Queue, error)
	GetMessageCount() int64
	Disconnect()
	Close() error
}
```


#### func  NewClient

```go
func NewClient(cfg *Config, pool dynworker.WorkerPool) Client
```

#### type Config

```go
type Config struct {
	URL                   string         `json:"url"`
	QueueName             string         `json:"queue_name"`
	QueueIsDurable        bool           `json:"queue_is_durable"`
	QueueDeleteWhenUnused bool           `json:"queue_delete_when_unused"`
	QueueIsExclusive      bool           `json:"queue_is_exclusive"`
	QueueNoWait           bool           `json:"queue_no_wait"`
	PrefetchCount         int64          `json:"prefetch_count"`
	PollInterval          types.Duration `json:"poll_interval"`
	ReconnectDelay        types.Duration `json:"reconnect_delay"`
	ConsumerName          string         `json:"consumer_name"`
	MaxRetryConnect       int            `json:"max_retry_connect"`
	MaxWorkers            int64          `json:"max_workers"`
	MinWorkers            int64          `json:"min_workers"`
	WorkerIdleTimeout     types.Duration `json:"worker_idle_timeout"`
}
```


#### func  NewConfig

```go
func NewConfig(
	parent context.Context,
	lgr logger.Logger,
	url, queuename string,
) *Config
```
Generate a new configuration object.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Generate a new default configuration object.

#### func (*Config) ConfigureWorkerPool

```go
func (obj *Config) ConfigureWorkerPool() *dynworker.Config
```

#### func (*Config) SetDialer

```go
func (obj *Config) SetDialer(dialer DialFn)
```

#### func (*Config) SetLogger

```go
func (obj *Config) SetLogger(lgr logger.Logger)
```

#### func (*Config) SetMessageHandler

```go
func (obj *Config) SetMessageHandler(callback dynworker.TaskFn)
```

#### func (*Config) SetParent

```go
func (obj *Config) SetParent(ctx context.Context)
```

#### func (*Config) Validate

```go
func (obj *Config) Validate()
```

#### type DialFn

```go
type DialFn func(url string) (amqpshim.Connection, error)
```
