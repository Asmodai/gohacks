-*- Mode: gfm -*-

# exec -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/exec"
```

## Usage

```go
const (
	STRING_LIT_PORT string = "-port"
	STRING_FMT_PORT string = "%d"
)
```

```go
const (
	CMD_SET_CTX = iota
	CMD_SET_PATH
	CMD_SET_ARGS
	CMD_SPAWN
	CMD_CHECK
	CMD_KILL_ALL
)
```

```go
var (
	CheckDelaySleep time.Duration = 250 * time.Millisecond
)
```

#### func  CheckProcs

```go
func CheckProcs(mgr process.IManager) error
```

#### func  KillAllProcs

```go
func KillAllProcs(mgr process.IManager) error
```

#### func  SetArgs

```go
func SetArgs(mgr process.IManager, args ...string) error
```

#### func  SetContext

```go
func SetContext(mgr process.IManager, ctx context.Context) error
```

#### func  SetPath

```go
func SetPath(mgr process.IManager, path string) error
```

#### func  Spawn

```go
func Spawn(mgr process.IManager, lgr logger.ILogger, cnf *Config) (*process.Process, error)
```

#### func  SpawnProcs

```go
func SpawnProcs(mgr process.IManager) error
```

#### type Args

```go
type Args struct {
}
```


#### func  NewArgs

```go
func NewArgs(port int, args []string) Args
```

#### func (Args) Get

```go
func (a Args) Get(port int) []string
```

#### type Config

```go
type Config struct {
	Count int // Number of processes.
	Base  int // Base port for RPC.
}
```


#### func  NewConfig

```go
func NewConfig(count, base int) *Config
```

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```

#### type Manager

```go
type Manager struct {
}
```


#### func  NewManager

```go
func NewManager(lgr logger.ILogger, count, base int) *Manager
```

#### func (*Manager) Check

```go
func (m *Manager) Check()
```

#### func (*Manager) Dump

```go
func (m *Manager) Dump()
```

#### func (*Manager) KillAll

```go
func (m *Manager) KillAll()
```

#### func (*Manager) SetArgs

```go
func (m *Manager) SetArgs(val Args)
```

#### func (*Manager) SetContext

```go
func (m *Manager) SetContext(val context.Context)
```

#### func (*Manager) SetCount

```go
func (m *Manager) SetCount(val int)
```

#### func (*Manager) SetPath

```go
func (m *Manager) SetPath(val string)
```

#### func (*Manager) Spawn

```go
func (m *Manager) Spawn()
```

#### type Process

```go
type Process struct {
}
```


#### func  NewProcess

```go
func NewProcess(lgr logger.ILogger, cnf *Config) *Process
```

#### func (*Process) Action

```go
func (p *Process) Action(state **process.State)
```

#### func (*Process) SetArgs

```go
func (p *Process) SetArgs(args ...string)
```

#### func (*Process) SetContext

```go
func (p *Process) SetContext(val context.Context)
```

#### func (*Process) SetCount

```go
func (p *Process) SetCount(val int)
```

#### func (*Process) SetPath

```go
func (p *Process) SetPath(val string)
```
