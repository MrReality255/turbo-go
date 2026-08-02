# TurboGo Reference

**Module:** `github.com/MrReality255/turbo-go`  
**Go Version:** 1.24+

TurboGo is a Go utility library providing generic collection helpers, concurrency primitives, a message broker, typed TCP sockets, and logging. It is organized into five packages.

---

## Package `tg/utils`

Core utility functions and types: generics-based collection operations, I/O helpers, math, concurrency primitives, and string/JSON tools.

---

### Generics & Core (`abstract.go`)

#### Type Constraints

| Constraint | Definition |
|---|---|
| `Numeric` | `constraints.Integer \| constraints.Float` |

#### `Coalesce[T any](values ...T) T`
Returns the first non-nil/non-zero value from the arguments. Uses reflection-based nil/zero detection.

#### `Distinct[T comparable](items []T) []T`
Returns a new slice with duplicate values removed, preserving order.

#### `DistinctSort[T comparable](items []T, compareFct func(T, T) bool) []T`
Returns distinct items, sorted by the provided comparison function.

#### `Flatten[T any](src [][]T) []T`
Flattens a two-dimensional slice into a single slice.

#### `IfThen[C any](cond bool, valTrue C, valFalse C) C`
Ternary-style conditional — returns `valTrue` if `cond` is true, else `valFalse`.

#### `IfThenFct[C any](cond bool, fctTrue func() C, fctFalse func() C) C`
Lazy ternary — only evaluates the selected branch function.

#### `IfThenAny(cond bool, valTrue any, valFalse any) interface{}`
Ternary returning an untyped `interface{}`.

#### `IfPass[C any](cond bool, fct func() (*C, error)) (*C, error)`
Executes `fct` only when `cond` is true; returns `(nil, nil)` otherwise.

#### `IsNil[T any](val T) bool`
Reflection-aware nil check that handles typed nils (e.g., an interface holding a nil pointer), zero-value numerics, and empty strings.

#### `Sort[T any](src []T, lessFct func(T, T) bool)`
In-place sort using a custom less function (wraps `sort.Sort`).

#### `RingBuffer[C any]`
A fixed-capacity circular buffer.

| Method | Description |
|---|---|
| `NewRingBuffer[C](length int) *RingBuffer[C]` | Create a new ring buffer of given capacity |
| `Add(items ...C)` | Append items, overwriting oldest when full |
| `Content() []C` | Return contents in insertion order |

---

### Arrays (`arrays.go`)

#### `ArrayCast[S Numeric, T Numeric](src []S) []T`
Converts a numeric slice from type `S` to type `T`.

#### `ArrayClone[T any](src []T) []T`
Returns a shallow copy of the slice.

#### `ArrayFilter[T any](src []T, checkFct func(T) bool) []T`
Returns elements for which `checkFct` returns true. Pass `nil` to copy all.

#### `ArrayFilterIdx[T any](src []T, checkFct func(T, int) bool) []T`
Filter with access to the element index.

#### `ArrayReduce[T any, S any](src []T, initValue S, cb func(S, T) S) S`
Reduces the slice to a single accumulated value.

#### `ArrayGroupBy[T any, C comparable](items []T, keyFct func(T) C) map[C][]T`
Groups elements into a map by a key derived from each element.

#### `ArrayHasAny[T any](arr []T, fct func(T) bool) bool`
Returns true if any element satisfies the predicate.

#### `ArrayMap[T any, V any](src []T, fct func(T) V) []V`
Transforms each element using `fct`.

#### `ArrayMapIdx[T any, V any](src []T, fct func(T, int) V) []V`
Map with access to the element index.

#### `ArrayMapEx[T any, V any](src []T, fct func(T) (V, bool)) []V`
Map with optional inclusion — only items where `bool` is true are kept.

#### `ArrayMapErr[T any, V any](src []T, fct func(T) (V, error)) ([]V, error)`
Map that propagates errors; stops on the first error.

#### `ArrayMapExErr[T any, V any](src []T, fct func(T) (V, bool, error)) ([]V, error)`
Combined filter+map with error propagation.

#### `ArrayMapErrIdx[T any, V any](src []T, fct func(T, int) (V, bool, error)) ([]V, error)`
Full-featured map: index access, optional inclusion, and error propagation.

#### `ArraySelect[T any](src []*T, checkFct func(*T, *T) bool) *T`
Selects a single pointer element by pairwise comparison (e.g., max/min).

#### `ArrayChoose[T any](src []*T, checkFct func(*T, *T) bool) *T`
Alias for `ArraySelect` (same behavior).

#### `ArrayChooseValue[T any](src []T, checkFct func(T, T) bool, emptyValue T) T`
Value-type variant — returns `emptyValue` for empty slices.

#### `ArrayFind[T any](src []T, cond func(T) bool) T`
Returns the first element matching the condition, or the zero value.

#### `ArraySort[T any](src []T, lessFct func(T, T) bool) []T`
Returns a sorted copy (non-destructive).

#### `SortArray[T any](src []T, lessFct func(T, T) bool)`
In-place sort of the provided slice.

#### `ArrayToMap[T comparable, S comparable, V any](ar []T, keyFct func(T) S, valueFct func(T) V) map[S]V`
Converts a slice to a map using key and value extraction functions.

#### `ArrayToMapEx[T comparable, S comparable, V any](ar []T, keyFct func(T) (S, bool), valueFct func(T) V) map[S]V`
Like `ArrayToMap` but the key function can exclude items by returning `false`.

---

### Maps (`maps.go`)

#### `MapClone[T comparable, S any](m map[T]S) map[T]S`
Returns a shallow copy of the map.

#### `MapKeys[C comparable, T any](src map[C]T, lessFct func(C, C) bool) []C`
Returns all map keys, optionally sorted.

#### `MapKeysIf[C comparable, T any](src map[C]T, condFct func(C, T) bool, lessFct func(C, C) bool) []C`
Returns filtered and optionally sorted keys.

#### `MapMerge[C comparable, T any](src map[C]T, other map[C]T, overwrite bool)`
Merges `other` into `src`. When `overwrite` is false, existing keys are not replaced.

#### `MapValues[C comparable, T any](src map[C]T, lessFct func(T, T) bool) []T`
Returns all values, optionally sorted.

#### `MapValuesIf[C comparable, T any](src map[C]T, selectFct func(C, T) bool, lessFct func(T, T) bool) []T`
Returns filtered and optionally sorted values.

#### `MapMap[K comparable, T any, K2 comparable, T2 any](src map[K]T, keyFct func(K) K2, valueFct func(K, T) T2) map[K2]T2`
Transforms both keys and values into a new map.

#### `MapMapValues[K comparable, T any, T2 any](src map[K]T, valueFct func(K, T) T2) map[K]T2`
Transforms values only, keeping keys unchanged.

#### `MapToArray[T comparable, S any, C any](src map[T]S, fct func(T, S) C, lessFct func(C, C) bool) []C`
Converts a map to a sorted array using a transformation function.

---

### SafeMap (`safe_map.go`)

A generic mutex-protected map.

```go
type SafeMap[K comparable, V any] struct { ... }
```

| Function / Method | Description |
|---|---|
| `NewSafeMap[K, V]() *SafeMap[K, V]` | Create an empty thread-safe map |
| `Set(key K, value V)` | Store a value |
| `Get(key K) V` | Retrieve a value (zero-value if missing) |
| `Delete(k K)` | Remove a key |
| `Prepare(id K, f func() V)` | Set key only if absent (lazy init) |
| `GetValues(lessFct func(V, V) bool) []V` | Return all values, optionally sorted |
| `Clone() map[K]V` | Return a plain map copy |
| `Contains(id K) bool` | Check key existence |

---

### Math (`math.go`)

#### `Clamp[T float64 | int](value, min, max T) T`
Constrains a value between min and max.

#### `Max[T INumeric](values ...T) T`
Returns the maximum of the provided values.

#### `Min[T INumeric](values ...T) T`
Returns the minimum of the provided values.

#### `Sum[T INumeric](values ...T) T`
Returns the sum of all values.

#### `Align[T constraints.Integer](src, step T) T`
Rounds down `src` to the nearest multiple of `step`.

#### `TruncStep[T constraints.Integer](src, step T) T`
Same as `Align` — truncates to the nearest lower multiple.

#### `LeastSquareSlope(values []float64) float64`
Calculates the slope via least-squares linear regression over indexed values.

#### `LinearFunction[S Numeric, T Numeric]`
A linear function `y = kx + q` with generic numeric types.

| Function / Method | Description |
|---|---|
| `NewLinearFunction[S, T](x1, y1, x2, y2 float64)` | Create from two points |
| `Eval(x S) T` | Compute y for given x |
| `Reverse(y T) S` | Compute x for given y |

#### `PermEx(fctResult func(int, []float64), fcts ...func([]float64) (float64, float64, float64))`
Generates all permutations across multiple numeric ranges and invokes `fctResult` for each combination.

#### `Counter`
A thread-safe incrementing counter.

| Method | Description |
|---|---|
| `NewCounter() *Counter` | Create a counter starting at 0 |
| `Next() int` | Increment and return |
| `Next64() int64` | Increment and return as int64 |
| `Inc()` | Increment without returning |
| `Get() int64` | Read current value |

---

### Strings & JSON (`strings.go`)

#### JSON Serialization

| Function | Description |
|---|---|
| `FromJSON(src interface{}, ptr interface{}) error` | Deserialize from `[]byte`, `string`, `io.Reader`, or Iris Context |
| `ParseJSON[T any](src any) (*T, error)` | Generic typed JSON deserialization |
| `ToJSON(obj interface{}) string` | Serialize to JSON string |
| `ToJSONB(obj interface{}) []byte` | Serialize to JSON bytes |
| `SaveToJSON(file string, obj interface{}, readable bool) error` | Write JSON to file (optionally indented) |

#### String Utilities

| Function | Description |
|---|---|
| `GetGUID() string` | Generate a UUID v4 string |
| `GetToken() string` | UUID without dashes (32 hex characters) |
| `Hash(input string, count int) int` | FNV-32a hash modulo `count` |
| `EvaluateAsStr(x interface{}) string` | Convert any value to string (supports `func() string`) |
| `SplitStr(s, by string) []string` | Split and remove empty segments |
| `StrToIntDef(str string, defaultValue int) int` | Parse int with fallback |
| `StrToInt64Def(str string, defaultValue int64) int64` | Parse int64 with fallback |
| `StrToNrDef[A constraints.Integer](str string, def A) A` | Generic integer parse with fallback |
| `StrToWordDef(str string, def uint16) uint16` | Parse uint16 with fallback |

#### Date/Time Parsing

| Function | Description |
|---|---|
| `DateTimeToStr(x time.Time) string` | Format as `"2006-01-02 15:04:05"` |
| `DateToStr(x time.Time) string` | Format as `"2006-01-02"` |
| `StrToTime(str string) time.Time` | Parse multiple formats (local timezone) |
| `StrToTimeUTC(str string) time.Time` | Parse multiple formats (UTC) |
| `StrToTimeLoc(str string, loc *time.Location) time.Time` | Parse with explicit location |

#### `StringList` / `IStrings`

A thread-safe string accumulator with optional distinct-only mode.

```go
func NewStringList(maxCount int, isDistinct bool, initItems ...string) IStrings
```

| Method | Description |
|---|---|
| `Add(items ...string)` | Append strings (skips duplicates in distinct mode) |
| `Addf(item string, args ...any)` | Append a formatted string |
| `Content() []string` | Return all items |
| `Sorted() []string` | Return items sorted alphabetically |
| `Join(sep string) string` | Join items with separator |
| `SortJoin(sep string) string` | Sort then join |

#### `CreateTable[T any](rows []T, cols []string, rowFct func(T) []any) string`
Renders a formatted ASCII table with auto-sized columns.

---

### Dates (`dates.go`)

#### `YearBegin(year int) time.Time`
Returns midnight on January 1 of the given year (local timezone).

#### `YearEnd(year int) time.Time`
Returns one second before midnight on December 31 (local timezone).

---

### Miscellaneous (`misc.go`)

#### `ConstFct[T any](value T) func() T`
Returns a closure that always returns `value`.

#### `IgnoreErr(_ error)`
Explicitly discards an error (avoids linter warnings).

#### `IsErr(err error, options ...error) bool`
Checks if `err` matches any of the provided sentinel errors via `errors.Is`.

#### `ExecErr(f interface{}) error`
Executes a function that is either `func()` or `func() error`.

#### `Must[T any](src T, err error) T`
Panics on error, otherwise returns the value.

#### `MustSucceed(err error)`
Panics if error is non-nil.

#### `StrToFloat64(str string) float64`
Parses a float, normalizing commas to dots. Returns 0 on failure.

#### `StrToFloat64E(str string) float64`
European-format parser: removes dots (thousands separator), converts comma to dot.

#### `FloatToStr(x float64) string`
Formats with 6 decimal places.

#### `FloatToStrP(x float64) string`
Formats with 12 decimal places (precision variant).

#### `CallWith[T any](callback func(func()), what func() T) T`
Executes `what` inside a `callback` wrapper (useful for mutex-protected returns).

#### `DataOrErr[T any](data T, err error) ItemWithErr[T]`
Wraps a value+error pair into an `ItemWithErr` struct.

#### `ItemWithErr[T any]`
```go
type ItemWithErr[T any] struct {
    Data T
    Err  error
}
```

#### `CountMap` / `Count(m CountMap)`
A map of `*int` to `bool` — `Count` increments each pointer where the value is true.

---

### Errors (`errors.go`)

#### `ErrorList`
A thread-safe list of errors that combines into a single semicolon-separated error.

| Function / Method | Description |
|---|---|
| `NewErrorList(maxCount int) *ErrorList` | Create with initial capacity |
| `Add(err error)` | Append (nil is ignored) |
| `Err() error` | Return combined error, or nil if empty |

---

### I/O (`io.go`)

#### File System

| Function | Description |
|---|---|
| `FileExists(file string) bool` | Check if a file exists |
| `DirExists(dir string) bool` | Check if a directory exists |
| `Dir(p string) ([]string, error)` | List all entries (files + dirs) |
| `DirEx(p string, wantDir, wantFiles bool) ([]string, error)` | Filtered directory listing |
| `MkDir(targetDir string) error` | Create directory tree (`MkdirAll`) |
| `MkDirFor(targetFile string) error` | Create parent directory for a file |

#### File Name Utilities

| Function | Description |
|---|---|
| `GetFileNameWithoutExt(filename string) string` | Base name without extension |
| `GetFileNameWithSuffix(path, suffix string) string` | Insert suffix before extension |

#### JSON File Operations

| Function | Description |
|---|---|
| `LoadFromJSON(file string, target interface{}) error` | Read and unmarshal JSON file |
| `NewFromFileJSON[T any](file string) (*T, error)` | Generic load-from-JSON |
| `LoadDirJSON[T any](dir string, prepFct func() *T) ([]*T, error)` | Load all JSON files in a directory |

#### Binary (Gob) Serialization

| Function | Description |
|---|---|
| `SaveToBin(filename string, obj any) error` | Serialize to file via `gob` |
| `LoadFromBin(filename string, ptr any) error` | Deserialize from file via `gob` |

#### CSV Operations

| Function | Description |
|---|---|
| `SaveToCSV[T any](filename string, data []T, header []string, rowCb func(T) []string) error` | Write typed data to CSV |
| `LoadFromCSV[T any](filename string, rowCb func([]string) *T, skipRows int) ([]*T, error)` | Read CSV into typed slice |

#### Binary Read/Write

| Function | Description |
|---|---|
| `ByteArray(items ...interface{}) ([]byte, error)` | Encode items to big-endian bytes |
| `MustByteArray(items ...interface{}) []byte` | Panicking variant |
| `WriteBytes(w io.Writer, items ...interface{}) error` | Write big-endian to writer |
| `FromBytes(b []byte, targets ...interface{}) error` | Decode big-endian from bytes |
| `FromReader(r io.Reader, targets ...interface{}) error` | Decode big-endian from reader |

#### Resource Management

| Function | Description |
|---|---|
| `CloseAfter(x io.Closer, fct func() error) error` | Execute function then close resource |
| `CloseAfterWith[T any](x io.Closer, fct func() (T, error)) (T, error)` | Generic variant with return value |

---

### Sync Helpers (`sync.go`)

#### `ExecLocked(ptr *sync.Mutex, fct func())`
Acquires the mutex, executes `fct`, then releases.

#### `ExecLockedErr(ptr *sync.Mutex, fct func() error) error`
Same as above but returns an error.

---

### Runner (`runner.go`)

A lifecycle manager for goroutine loops with close/wait semantics.

```go
type IRunner interface {
    Close() error
    ExecLocked(fct func())
    Start()
    Wait()
}
```

#### `NewRunner(onNext func() bool, onClose func() error) IRunner`
- `onNext`: called in a loop; return `false` to stop.
- `onClose`: called once when closing.
- `Start()` launches the loop goroutine.
- `Wait()` blocks until closed.
- `ExecLocked(fct)` runs `fct` only if not yet closed.

---

### Perceptron (`perceptron.go`)

A simple single-layer perceptron for binary classification.

```go
type Perceptron struct {
    Weights      []float64
    bias         float64
    learningRate float64
}
```

| Function / Method | Description |
|---|---|
| `NewPerceptron(inputSize int, learningRate float64) *Perceptron` | Initialize with random weights |
| `Train(inputs PerceptronItemList, epochs int)` | Train on labeled data |
| `Predict(input []float64) float64` | Returns 1 or 0 |
| `Save(file string) error` | Save weights to JSON |
| `Load(file string) error` | Load weights from JSON |

```go
type PerceptronItem struct {
    Flags      []float64
    IsSelected bool
}
type PerceptronItemList []*PerceptronItem
```

---

### Web (`web.go`)

An Iris-framework helper for HTTP handlers.

```go
type WebContextHandler struct { ... }
```

| Function / Method | Description |
|---|---|
| `NewContextHandler(ctx iris.Context, onError func(string, error)) *WebContextHandler` | Create handler |
| `RespondJSON(data interface{})` | Write 200 + JSON body |
| `RespondJSONErr(content any, err error)` | Respond or return 500 on error |
| `ClientErr(ctx iris.Context, err error)` | Respond with 400 |
| `SrvErr(err error)` | Respond with 500 |
| `GetParam(key string) string` | Read URL parameter |

---

### Testing Helpers (`testing.go`)

| Function | Description |
|---|---|
| `TestAsString(t, testCase int, hint, required string, actual any)` | Compare `fmt.Sprintf("%v", actual)` to expected |
| `TestString(t, testCase int, hint, required, actual string)` | Direct string equality assertion |
| `TestStringF(t, testCase int, hint, want, got string, args ...any)` | Assert with formatted actual |

---

### App (`app.go`)

#### `PrintErr(msg string, err error)`
Prints message + error to stdout if error is non-nil.

---

## Package `tg/sync`

Concurrency building blocks: batch processing, throttled goroutine launching, and wait groups.

---

### ConcurrentLauncher (`concurrent_launcher.go`)

Controls concurrent goroutine execution with a semaphore-style limit.

```go
type IConcurrentLauncher interface {
    Go(fct interface{})
    Locked(fct interface{}) error
    Wait() error
}
```

#### `NewConcurrentLauncher(maxCount int, limit int) IConcurrentLauncher`
- `maxCount`: expected total operations (sizes error list).
- `limit`: max concurrent goroutines (0 = unlimited).
- `Go(fct)`: launch a goroutine (blocks if at limit).
- `Wait()`: block until all goroutines finish; returns combined errors.
- `Locked(fct)`: execute under a shared mutex (for serialized sections).

---

### WaitGroup (`wait_group.go`)

A higher-level wait group with error capture and concurrency limiting.

```go
type IWaitGroup interface {
    Go(fct any)
    Wait() error
}
```

#### `NewWaitGroup(maxCount int, threadCount int) IWaitGroup`
- `maxCount`: sizes the launcher.
- `threadCount`: max concurrent goroutines (defaults to `maxCount` if 0).
- `Go(fct)`: submit work (`func()` or `func() error`).
- `Wait()`: blocks until all work completes; returns first error encountered.

---

### BatchProcessor (`batch_processor.go`)

Accumulates items and flushes in batches triggered by size or timeout.

```go
type BatchFlushReason uint8
const (
    BatchFlushTimeout  BatchFlushReason = 1  // timer expired
    BatchFlushOverflow BatchFlushReason = 2  // batch full
    BatchFlushExplicit BatchFlushReason = 3  // manual Flush() call
)

type BatchProcessorConfig struct {
    FlushThreads int           // max concurrent flush operations
    Length       int           // batch capacity
    Timeout      time.Duration // max time before auto-flush
}

type IBatchProcessor interface {
    Add(adderFct func() int) error
    Flush() error
}
```

#### `NewBatchProcessor(config BatchProcessorConfig, onSwap func(BatchFlushReason) func() error) IBatchProcessor`
- `onSwap` is called when a batch is ready to flush. It receives the reason and must return a flush function that will be executed (potentially in a separate goroutine).
- `Add(adderFct)`: the adder function adds items and returns the count added.
- `Flush()`: force-flush the current batch synchronously.

---

### BufferProcessor (`buffer_processor.go`)

A typed wrapper around `BatchProcessor` that manages an internal buffer.

```go
type IBufferProcessor[T any] interface {
    Add(items ...T) error
    Flush() error
}
```

#### `NewBufferProcessor[T any](cfg BatchProcessorConfig, onFlush func([]T, BatchFlushReason) error) IBufferProcessor[T]`
- Accumulates items of type `T` internally.
- Flushes via `onFlush` when batch size or timeout is reached.

---

## Package `tg/broker`

A generic in-process message broker for request/response and pub/sub communication between typed members.

---

### Handle (`utils.go`)

A 64-bit identifier composed of a 32-bit type ID (high) and 32-bit sequence ID (low).

```go
type Handle uint64
const HandleAny = 0
```

| Function / Method | Description |
|---|---|
| `NewHandle(typeID, seqID uint32) Handle` | Construct a handle |
| `NewHandleType(typeID uint32) Handle` | Handle with seqID = 0 (type-level) |
| `GetTypeID() uint32` | Extract type portion |
| `GetSeqID() uint32` | Extract sequence portion |

---

### Broker (`broker.go`)

```go
type IBroker[Cmd ICommand] interface {
    AddMember(id Handle, messageHandler MemberMessageHandler[Cmd]) IMember[Cmd]
    Close()
}
```

#### `New[Command ICommand](descriptor CommandDescriptor[Command], timeout time.Duration) IBroker[Command]`
Creates a broker that dispatches messages through an internal queue.

- `CommandDescriptor` provides `GetID` and `GetRef` functions to extract handle/reference from commands.
- Messages are dispatched based on receiver handle:
  - Exact match (non-zero seqID): delivered to that specific member.
  - Type-level (seqID = 0): delivered to one member of that type.
  - Subscribers: delivered to all members subscribed to the command's type.

---

### Member (`member.go`)

```go
type IMember[Command ICommand] interface {
    Request(receiver Handle, cmd Command) (Command, error)
    RequestMultiple(receiver Handle, cmd Command, handler RequestHandler[Command])
    Send(receiver Handle, cmd Command)
    Subscribe(cmdType ...uint32)
    Close()
}
```

| Method | Description |
|---|---|
| `Request` | Send a command and wait for a single response |
| `RequestMultiple` | Send and receive multiple responses via callback |
| `Send` | Fire-and-forget message to a receiver |
| `Subscribe` | Register interest in command types (pub/sub) |
| `Close` | Remove this member from the broker |

---

### RequestManager (`request_manager.go`)

Manages outstanding requests with timeout tracking.

```go
type IRequestManager[Command ICommand] interface {
    Abort()
    Accept(msg Command)
    Request(req Command) (Command, error)
    RequestMultiple(req Command, handler func(Command, error) bool)
}
```

#### `NewRequestManager[Command](descriptor, senderFct, receiverFct, timeout) IRequestManager[Command]`

- `senderFct`: called to transmit the request command.
- `receiverFct`: called for incoming messages that don't match an active request.
- `timeout`: per-request deadline.
- `Abort()`: cancels all pending requests with `ErrRequestAborted`.

**Sentinel Errors:**
- `ErrHandleConflict` — duplicate request ID
- `ErrRequestTimeout` — deadline exceeded
- `ErrRequestAborted` — broker or member was aborted

---

## Package `tg/comm`

TCP networking with typed serialization.

---

### TCP Utilities (`utils.go`)

| Function | Description |
|---|---|
| `GetServerAddr(port int) string` | Returns `"0.0.0.0:<port>"` |
| `Serve(addr string, handler func(net.Conn) error, errHandler func(error)) error` | Accept TCP connections in a loop |
| `ServeLocal(port int, handler, errHandler) error` | Shortcut for `Serve` on all interfaces |
| `NewTcpClient(addr string, port int) (net.Conn, error)` | Dial a TCP connection |
| `NewTcpServer(port int, handler func(net.Conn), errHandler func(error)) (io.Closer, error)` | Non-blocking TCP server (returns listener) |
| `GetHttpFrom(url string, target interface{}) error` | HTTP GET + JSON decode into target |

---

### TypedSocket (`typed_socket.go`)

A generic typed read/write abstraction over raw TCP connections.

```go
type ITypedSocket[T any] interface {
    io.Closer
    Read() (T, error)
    Wait()
    Write(data T) error
}
```

#### `TypedSocketFactory[T any]`

```go
func NewTypedSocketFactory[T any](
    onRead func(conn IAbstractSocket) (T, error),
    onWrite func(conn IAbstractSocket, data T) error,
) *TypedSocketFactory[T]
```

| Method | Description |
|---|---|
| `New(src IAbstractSocket) ITypedSocket[T]` | Wrap an existing connection |
| `NewTcpClient(addr string, port int) (ITypedSocket[T], error)` | Connect and wrap |

---

## Package `tg/log`

Structured logging facade built on [logrus](https://github.com/sirupsen/logrus) with optional file output.

### Setup

```go
func Setup(setup *Config)
```

`Config.Filename` — if set, logs to file (and stdout). Otherwise stdout only. Sets log level to Debug.

### Logging Functions

| Function | Description |
|---|---|
| `Debug(msg string, params ...any)` | Debug-level log |
| `Info(msg string, params ...any)` | Info-level log |
| `Warn(msg string, params ...any)` | Warning-level log |
| `LogError(msg string, err error)` | Error-level log |
| `IfError(msg string, err error)` | Log error only if err != nil |
| `LogIfWarnErr(msg string, err error)` | Log warning only if err != nil |
| `LogWarnErr(msg string, err error)` | Always log warning with error |
| `LogIfErrCtx(msg string, err error, ctx Context)` | Log error with structured context |

### Context

```go
type Context map[string]any
```

Implements `String()` returning `{key=value, ...}` format, sorted by key.

---

## Dependencies

| Dependency | Purpose |
|---|---|
| `github.com/google/uuid` | UUID generation |
| `github.com/kataras/iris/v12` | Web framework (context parsing) |
| `github.com/sirupsen/logrus` | Structured logging |
| `golang.org/x/exp` | Generic constraints (`constraints.Integer`, etc.) |
