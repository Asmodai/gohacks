<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# amqp -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/amqp"
```

## Usage

```go
var (
	// Signalled when there is no hostname in the AMQP configuration.
	ErrNoHostname = errors.Base("no AMQP hostname provided")

	// Signalled when there is no queue name in the AMQP configuration.
	ErrNoQueueName = errors.Base("no AMQP queue name provided")
)
```

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
func NewClient(ctx context.Context, cfg *Config, pool dynworker.WorkerPool) Client
```

#### type Config

```go
type Config struct {
	Username              string         `json:"username"`
	Password              string         `config_obscure:"true" json:"password"`
	Hostname              string         `json:"hostname"`
	Port                  int            `json:"port"`
	VirtualHost           string         `json:"vhost"`
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
func NewConfig(hostname, virtualhost, queuename string) *Config
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
Generate a worker pool configuration.

#### func (*Config) IsValidated

```go
func (obj *Config) IsValidated() bool
```
Has the configuration been validated?

#### func (*Config) MakeWorkerPool

```go
func (obj *Config) MakeWorkerPool(ctx context.Context) dynworker.WorkerPool
```
Generate a worker pool.

#### func (*Config) SetDialer

```go
func (obj *Config) SetDialer(dialer DialFn)
```
Set the dialer function.

This is useful for mocking.

#### func (*Config) SetMessageHandler

```go
func (obj *Config) SetMessageHandler(callback dynworker.TaskFn)
```
Set the message handler worker function.

#### func (*Config) URL

```go
func (obj *Config) URL() string
```
Compose the AMQP URL.

#### func (*Config) Validate

```go
func (obj *Config) Validate() []error
```
Validate the AMQP configuration.

This *must* be called before any attempt to use the AMQP configuration with a
client is made.

The idea here is that we use the `config` package and its `Validate` methods.

#### type DialFn

```go
type DialFn func(url string) (amqpshim.Connection, error)
```
