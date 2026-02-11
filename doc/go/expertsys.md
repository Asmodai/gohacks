<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# expertsys -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/expertsys"
```

## Usage

```go
var (
	ErrNotStable = errx.Base("engine did not stabilise")
)
```

#### type Engine

```go
type Engine struct {
}
```

Expert system engine.

#### func (*Engine) RunToFixpoint

```go
func (e *Engine) RunToFixpoint(wmem WorkingMemory, maxIters int) (int, error)
```

#### type WorkingMemory

```go
type WorkingMemory interface {
	dag.Filterable

	// Return the current version of the working memory.
	Version() uint64
}
```

Working memory interface.

This interface defines the fact storage system for the expert system. It is
derived from `dag.Filterable`.

#### func  NewWorkingMemory

```go
func NewWorkingMemory() WorkingMemory
```
Create a new working memory instance.
