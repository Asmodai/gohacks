<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# fileio -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/fileio"
```

## Usage

```go
var (

	// Signalled when a file is not a regular file.
	//
	// That is not a symlink, pipe, socket, device, etc.
	ErrNotRegular = errors.Base("not a regular file")

	// Signalled when an attempt is made to process a file that is
	// too large.
	//
	// The size limit is configurable via `Options`.
	ErrTooLarge = errors.Base("file exceeds size limit")

	// Signalled if a size option is invalid.
	ErrInvalidSize = errors.Base("invalid size")

	// Signalled if we're trying to operate on a symbolic link without
	// `FollowSymlinks` enabled.
	ErrSymlinkDenied = errors.Base("symlink not allowed")
)
```

#### type Chunk

```go
type Chunk struct {
	Offset int64
	Data   []byte
}
```


#### type Options

```go
type Options struct {
	// Maximum number of bytes to read.
	//
	// If this value is set to a negative number then it is interpreted
	// as "read an unlimited number of bytes".
	//
	// However "unlimited", in this case, equates to `math.MaxInt` bytes.
	//
	// The reason for this is that we return a slice of bytes, and the
	// maximum number of elements in a Go slice is `math.MaxInt`.
	//
	// The default value is 0.
	MaxReadBytes int64

	// Should symbolic links be followed?
	//
	// If false, then symbolic links are not followed.
	//
	// The default value is false.
	FollowSymlinks bool
}
```

File I/O options.

#### type Reader

```go
type Reader interface {
	// The file name that we wish to load.
	Filename() string

	// Check whether the file exists.
	//
	// If the file exists, then `true` is returned along with no error.
	//
	// If the file does not exist, then `false` is returned along with
	// no error.
	//
	// If the file exists but is not a regular file, then false is
	// returned along with `ErrNotRegular`.
	//
	// If we are following symbolic links and the file exists and is
	// a symbolic link then it is resolved to the symbolic link's target
	// and, if that exists, `true` is returned along with no error.
	Exists() (bool, error)

	// Check if the file is a symbolic link.
	IsSymlink() (bool, error)

	// Read the file and return a byte array of its content.
	//
	// If `MaxReadBytes` in the options is zero then the number of bytes
	// read will be at most `math.MaxInt`.
	//
	// If `MaxReadBytes` in the options is higher than zero then the
	// number of bytes read will be at most `MaxReadBytes`.
	//
	// `MaxReadBytes` can never be negative and can never exceed
	// `math.MaxInt`.
	Load() ([]byte, error)

	// Open the file and return an `io.ReadCloser`.
	//
	// `MaxReadBytes` is ignored.
	Open() (io.ReadCloser, error)

	// Open the file and stream its contents to the specified writer.
	//
	// If `limit` is zero then the entire contents shall be copied.
	//
	// If `limit` is greater than zero, then at most `limit` bytes will
	// be copied.
	CopyTo(writer io.Writer, limit int64) (int64, error)

	// Stream chunks of up to `chunkSize` bytes.
	//
	// `bufSize` will utilise a readahead buffer of the given size.
	//
	// If `limit` is zero then the entirety of the content will be
	// streamed.
	//
	// if `chunkSize` is zero or lower then a default chunk size of
	// 64 * 1024 shall be used.
	Stream(ctx context.Context, chunkSize, bufSize int, limit int64) StreamResult
}
```

File reader.

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

#### func  NewReaderWithFile

```go
func NewReaderWithFile(filename string) (Reader, error)
```
Create a new FileReader with the given file name.

#### func  NewReaderWithFileAndOptions

```go
func NewReaderWithFileAndOptions(filename string, opts Options) (Reader, error)
```
Create a new FileReader with the given file name and options.

#### type StreamResult

```go
type StreamResult struct {
	ChunkCh <-chan Chunk // Channel for chunks.
	ErrorCh <-chan error // Channel for errors.
	Close   func()       // Cancel the stream.
	Wait    func() error // Wait for completion.
}
```
