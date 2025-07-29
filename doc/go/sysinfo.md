<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# sysinfo -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/sysinfo"
```

## Usage

#### func  Spawn

```go
func Spawn(ctx context.Context, interval types.Duration) (*process.Process, error)
```
Spawn a system information process.

The provided context must have a `process.Manager` entry in its user value. See
`contextdi` and `process.SetProcessManager`.

#### type Proc

```go
type Proc struct {
}
```


#### func  NewProc

```go
func NewProc() *Proc
```
Create a new system information process with default values.

#### func (*Proc) Action

```go
func (sip *Proc) Action(state *process.State)
```
Function that runs every tick in the sysinfo process.

Simply prints out the Go runtime stats via the process's logger.

#### type SysInfo

```go
type SysInfo struct {
	sync.Mutex
}
```

System information poller.

#### func  NewSysInfo

```go
func NewSysInfo() *SysInfo
```
Create a new System Information instance.

#### func (*SysInfo) Allocated

```go
func (si *SysInfo) Allocated() uint64
```
Return number of MiB currently allocated.

#### func (*SysInfo) GC

```go
func (si *SysInfo) GC() uint32
```
Return the number of collections performed.

#### func (*SysInfo) GoRoutines

```go
func (si *SysInfo) GoRoutines() int
```
Return the number of Go routines.

#### func (*SysInfo) Heap

```go
func (si *SysInfo) Heap() uint64
```
Return number of MiB used by the heap.

#### func (*SysInfo) Hostname

```go
func (si *SysInfo) Hostname() string
```
Return this system's hostname.

#### func (*SysInfo) RunTime

```go
func (si *SysInfo) RunTime() time.Duration
```
Return the time running.

#### func (*SysInfo) System

```go
func (si *SysInfo) System() uint64
```
Return number of MiB allocated from the system.

#### func (*SysInfo) UpdateStats

```go
func (si *SysInfo) UpdateStats()
```
Update runtime statistics.
