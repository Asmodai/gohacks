-*- Mode: gfm -*-

# sysinfo -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/sysinfo"
```

## Usage

#### func  Spawn

```go
func Spawn(mgr process.Manager, interval int) (*process.Process, error)
```

#### type SysInfo

```go
type SysInfo struct {
	sync.Mutex
}
```


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

#### type SysInfoProc

```go
type SysInfoProc struct {
}
```


#### func  NewSysInfoProc

```go
func NewSysInfoProc() *SysInfoProc
```

#### func (*SysInfoProc) Action

```go
func (sip *SysInfoProc) Action(state **process.State)
```
