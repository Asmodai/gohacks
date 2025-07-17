-*- Mode: gfm -*-

# fileloader -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/fileloader"
```

## Usage

#### type FileLoader

```go
type FileLoader interface {
	// The file path that we wish to load.
	Filename() string

	// Check if the file exists.
	Exists() (bool, error)

	// Read the file and return a byte array of its content.
	Load() ([]byte, error)
}
```

File loader.

A utility that provides file opening functionality wrapped in a mockable
interface.

To use:

    1. Create an instance with the file path you wish to load:

```go

    load := fileloader.NewWithFile("/path/to/file")

```

    2. Check it exists (optional):

```go

    found, err := load.Exists()

    if err != nil {
    	panic("File does not exist: " + err.Error())
    }

```

    3. Load your file:

```go

    data, err := load.Load()
    if err != nil {
    	panic("Could not load file: " + err.Error())
    }

```

The `Load` method returns the file content as a byte array.

#### func  NewWithFile

```go
func NewWithFile(filename string) FileLoader
```
Create a new FileLoader object with the given file name.
