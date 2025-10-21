Go Language Features Used in Security-Scanner-StackGaurd-Assignment
=======================================================================

LIST OF GO LANGUAGE FEATURES WITH EXPLANATIONS
==============================================

1. **Goroutines**
   - Lightweight threads managed by the Go runtime
   - Used for concurrent execution of functions
   - Example: `go func() { /* concurrent work */ }()`

2. **Channels**
   - Typed conduits for communication between goroutines
   - Used for signal handling and data passing
   - Example: `quit := make(chan os.Signal, 1)`

3. **Buffered Channels (Semaphores)**
   - Channels with capacity used as counting semaphores
   - Limit concurrent operations
   - Example: `semaphore := make(chan struct{}, 30)`

4. **sync.WaitGroup**
   - Synchronization primitive to wait for goroutines to complete
   - Tracks number of active goroutines
   - Methods: Add(), Done(), Wait()

5. **sync.Mutex**
   - Mutual exclusion lock for protecting shared data
   - Prevents race conditions
   - Methods: Lock(), Unlock()

6. **defer**
   - Defers execution of a function until surrounding function returns
   - Used for cleanup operations
   - Example: `defer file.Close()`

7. **Context**
   - Carries deadlines, cancellation signals, and request-scoped values
   - Used for timeout handling and cancellation
   - Example: `context.WithTimeout(context.Background(), 30*time.Second)`

8. **Atomic Operations**
   - Thread-safe operations on variables
   - Used for counters without locks
   - Example: `atomic.AddInt32(&counter, 1)`

9. **Structs**
   - Composite data types with named fields
   - Define custom data structures
   - Example: `type Config struct { Port string }`

10. **Maps**
  - Key-value data structures
  - Dynamic associative arrays
  - Example: `map[string]interface{}{}`

11. **Slices**
  - Dynamic arrays that can grow
  - Reference types pointing to underlying arrays
  - Example: `[]models.Finding`

12. **For Loops**
  - Iteration constructs
  - Regular: `for i := 0; i < n; i++`
  - Range: `for i, v := range slice`

13. **If Statements**
  - Conditional execution
  - Often used with error checking: `if err != nil`

14. **Switch Statements**
  - Multi-way branch statements
  - Used for multiple condition checking
  - Example: `switch config.LogLevel`

15. **Short Variable Declarations**
  - `:=` operator for declaration and assignment
  - Type inferred from right-hand side
  - Example: `result := someFunction()`

16. **Regular Functions**
  - Named code blocks that can be called
  - Take parameters and return values
  - Example: `func main() { ... }`

17. **Methods**
  - Functions associated with a type (receivers)
  - Object-oriented style in Go
  - Example: `func (c *Config) Validate() error`

18. **Anonymous Functions**
  - Functions without names, often used as closures
  - Can be assigned to variables or called immediately
  - Example: `go func() { ... }()`

19. **Error Interface**
  - Built-in interface for error handling
  - Convention: `if err != nil { return err }`

20. **Type Assertions**
  - Extract concrete type from interface{}
  - Two-value form: `value, ok := data.(string)`

21. **Constants**
  - Immutable values known at compile time
  - Declared with `const` keyword

22. **Variables**
  - Named storage locations
  - Package-level with `var`, local with `:=`

23. **Pointers**
  - Types that hold memory addresses
  - Dereference with `*`, address with `&`
  - Example: `*Config`, `&config`

24. **Interface Types**
  - Define method sets that types must implement
  - `interface{}` is the empty interface (any type)
  - Example: `io.Writer`

25. **make() Function**
  - Allocates and initializes slices, maps, channels
  - Returns initialized (not zero) values
  - Example: `make([]int, 0, 10)`

26. **append() Function**
  - Adds elements to slices
  - May reallocate underlying array
  - Example: `slice = append(slice, item)`

27. **len() Function**
  - Returns length of strings, arrays, slices, maps, channels
  - Example: `len(slice)`

28. **cap() Function**
  - Returns capacity of slices, channels
  - Example: `cap(slice)`

29. **String Operations**
  - `strings.Split()` - Split string into substrings
  - `strings.ToLower()` - Convert to lowercase
  - `filepath.Ext()` - Get file extension

30. **Regular Expressions**
  - `regexp.MustCompile()` - Compile regex pattern
  - `regexp.MatchString()` - Test if string matches pattern

31. **Time Operations**
  - `time.Now()` - Get current time
  - `time.Format()` - Format time as string
  - Time parsing and duration calculations

32. **UUID Generation**
  - `uuid.New()` - Generate unique identifier
  - `uuid.New().String()` - Convert to string format

33. **Structured Logging (slog)**
  - `slog.Info()`, `slog.Error()`, etc. - Log messages
  - JSON and text format handlers
  - Configurable log levels

34. **File I/O**
  - `os.OpenFile()` - Open file with flags
  - `io.ReadAll()` - Read entire file content
  - `io.MultiWriter()` - Write to multiple destinations

35. **Environment Variables**
  - `os.Getenv()` - Read environment variable
  - Used for configuration

36. **JSON Handling**
  - `json.Unmarshal()` - Parse JSON into Go structs
  - `json.Marshal()` - Convert Go structs to JSON

37. **HTTP Operations**
  - `http.Client` - HTTP client with timeout/transport
  - `http.Get()` - Make GET requests
  - Status code checking

38. **Command Line Arguments**
  - `os.Args` - Access command line arguments
  - Used for program configuration

39. **Template Rendering**
  - Text/template for HTML generation
  - Template execution with data binding

40. **Type Switching**
  - Switch on interface{} types
  - Determine concrete type at runtime

DETAILED EXPLANATIONS
=====================

**Concurrency Features:**
- **Goroutines**: Enable concurrent execution without the complexity of threads. The Go runtime multiplexes goroutines onto OS threads efficiently.
- **Channels**: Provide safe communication between goroutines. Buffered channels can hold multiple values, unbuffered require synchronous communication.
- **Semaphores**: Implemented with buffered channels to limit concurrent access to resources, preventing system overload.
- **WaitGroups**: Ensure all goroutines complete before proceeding, crucial for synchronization in batch operations.
- **Mutexes**: Prevent race conditions when multiple goroutines access shared data simultaneously.
- **defer**: Guarantees cleanup code runs regardless of how a function exits, essential for resource management.
- **Context**: Enables cancellation and timeout propagation through call chains, improving reliability.
- **Atomic Operations**: Provide lock-free thread-safe operations for simple counters and flags.

**Data Structures:**
- **Structs**: Form the backbone of Go programs, grouping related data fields together.
- **Maps**: Provide O(1) average-time complexity lookups, ideal for caching and configuration.
- **Slices**: Flexible arrays that grow automatically, used extensively for collections.

**Control Flow:**
- **For Loops**: Go's only looping construct, versatile for different iteration patterns.
- **If Statements**: Standard conditional execution, often combined with short variable declarations.
- **Switch**: Cleaner than multiple if-else chains, supports type switching for interfaces.

**Functions:**
- **Regular Functions**: Basic building blocks of Go programs.
- **Methods**: Associate functions with types, enabling object-oriented patterns.
- **Anonymous Functions**: Create closures, useful for goroutines and deferred execution.

**Type System:**
- **Constants**: Compile-time evaluated values, improving performance and preventing accidental changes.
- **Pointers**: Allow efficient passing of large structures and enable mutation of function parameters.
- **Interfaces**: Enable polymorphism and dependency injection patterns.

**Built-in Functions:**
- **make()**: Essential for creating reference types (slices, maps, channels) with proper initialization.
- **append()**: Core slice manipulation, handles growth and reallocation automatically.
- **len()**/**cap()**: Provide metadata about data structures for bounds checking and optimization.

**Standard Library Usage:**
- **String Operations**: Essential for text processing and file path manipulation.
- **Regular Expressions**: Powerful pattern matching for secret detection and validation.
- **Time Operations**: Critical for logging, timeouts, and scheduling.
- **UUID Generation**: Creates unique identifiers for tracking requests and resources.
- **Logging**: Structured logging with levels and JSON output for production monitoring.
- **File I/O**: Basic file operations for configuration and data persistence.
- **Environment Variables**: Configuration without hardcoding sensitive values.
- **JSON**: Standard data interchange format, ubiquitous in web APIs.
- **HTTP**: Foundation of web services and API communication.
- **Templates**: Server-side HTML generation for web interfaces.

This comprehensive list covers all major Go language features used in the Security Scanner project, demonstrating modern Go programming practices with strong emphasis on concurrency, error handling, and clean architecture.

Go Language Features Used in Security-Scanner-StackGaurd-Assignment
=======================================================================

LIST OF GO LANGUAGE FEATURES WITH EXPLANATIONS
==============================================

1. **Goroutines**
   - Lightweight threads managed by the Go runtime
   - Used for concurrent execution of functions
   - Example: `go func() { /* concurrent work */ }()`

2. **Channels**
   - Typed conduits for communication between goroutines
   - Used for signal handling and data passing
   - Example: `quit := make(chan os.Signal, 1)`

3. **Buffered Channels (Semaphores)**
   - Channels with capacity used as counting semaphores
   - Limit concurrent operations
   - Example: `semaphore := make(chan struct{}, 30)`

4. **sync.WaitGroup**
   - Synchronization primitive to wait for goroutines to complete
   - Tracks number of active goroutines
   - Methods: Add(), Done(), Wait()

5. **sync.Mutex**
   - Mutual exclusion lock for protecting shared data
   - Prevents race conditions
   - Methods: Lock(), Unlock()

6. **defer**
   - Defers execution of a function until surrounding function returns
   - Used for cleanup operations
   - Example: `defer file.Close()`

7. **Context**
   - Carries deadlines, cancellation signals, and request-scoped values
   - Used for timeout handling and cancellation
   - Example: `context.WithTimeout(context.Background(), 30*time.Second)`

8. **Atomic Operations**
   - Thread-safe operations on variables
   - Used for counters without locks
   - Example: `atomic.AddInt32(&counter, 1)`

9. **Structs**
   - Composite data types with named fields
   - Define custom data structures
   - Example: `type Config struct { Port string }`

10. **Maps**
    - Key-value data structures
    - Dynamic associative arrays
    - Example: `map[string]interface{}{}`

11. **Slices**
    - Dynamic arrays that can grow
    - Reference types pointing to underlying arrays
    - Example: `[]models.Finding`

12. **For Loops**
    - Iteration constructs
    - Regular: `for i := 0; i < n; i++`
    - Range: `for i, v := range slice`

13. **If Statements**
    - Conditional execution
    - Often used with error checking: `if err != nil`

14. **Switch Statements**
    - Multi-way branch statements
    - Used for multiple condition checking
    - Example: `switch config.LogLevel`

15. **Short Variable Declarations**
    - `:=` operator for declaration and assignment
    - Type inferred from right-hand side
    - Example: `result := someFunction()`

16. **Regular Functions**
    - Named code blocks that can be called
    - Take parameters and return values
    - Example: `func main() { ... }`

17. **Methods**
    - Functions associated with a type (receivers)
    - Object-oriented style in Go
    - Example: `func (c *Config) Validate() error`

18. **Anonymous Functions**
    - Functions without names, often used as closures
    - Can be assigned to variables or called immediately
    - Example: `go func() { ... }()`

19. **Error Interface**
    - Built-in interface for error handling
    - Convention: `if err != nil { return err }`

20. **Type Assertions**
    - Extract concrete type from interface{}
    - Two-value form: `value, ok := data.(string)`

21. **Constants**
    - Immutable values known at compile time
    - Declared with `const` keyword

22. **Variables**
    - Named storage locations
    - Package-level with `var`, local with `:=`

23. **Pointers**
    - Types that hold memory addresses
    - Dereference with `*`, address with `&`
    - Example: `*Config`, `&config`

24. **Interface Types**
    - Define method sets that types must implement
    - `interface{}` is the empty interface (any type)
    - Example: `io.Writer`

25. **make() Function**
    - Allocates and initializes slices, maps, channels
    - Returns initialized (not zero) values
    - Example: `make([]int, 0, 10)`

26. **append() Function**
    - Adds elements to slices
    - May reallocate underlying array
    - Example: `slice = append(slice, item)`

27. **len() Function**
    - Returns length of strings, arrays, slices, maps, channels
    - Example: `len(slice)`

28. **cap() Function**
    - Returns capacity of slices, channels
    - Example: `cap(slice)`

29. **String Operations**
    - `strings.Split()` - Split string into substrings
    - `strings.ToLower()` - Convert to lowercase
    - `filepath.Ext()` - Get file extension

30. **Regular Expressions**
    - `regexp.MustCompile()` - Compile regex pattern
    - `regexp.MatchString()` - Test if string matches pattern

31. **Time Operations**
    - `time.Now()` - Get current time
    - `time.Format()` - Format time as string
    - Time parsing and duration calculations

32. **UUID Generation**
    - `uuid.New()` - Generate unique identifier
    - `uuid.New().String()` - Convert to string format

33. **Structured Logging (slog)**
    - `slog.Info()`, `slog.Error()`, etc. - Log messages
    - JSON and text format handlers
    - Configurable log levels

34. **File I/O**
    - `os.OpenFile()` - Open file with flags
    - `io.ReadAll()` - Read entire file content
    - `io.MultiWriter()` - Write to multiple destinations

35. **Environment Variables**
    - `os.Getenv()` - Read environment variable
    - Used for configuration

36. **JSON Handling**
    - `json.Unmarshal()` - Parse JSON into Go structs
    - `json.Marshal()` - Convert Go structs to JSON

37. **HTTP Operations**
    - `http.Client` - HTTP client with timeout/transport
    - `http.Get()` - Make GET requests
    - Status code checking

38. **Command Line Arguments**
    - `os.Args` - Access command line arguments
    - Used for program configuration

39. **Template Rendering**
    - Text/template for HTML generation
    - Template execution with data binding

40. **Type Switching**
    - Switch on interface{} types
    - Determine concrete type at runtime

DETAILED EXPLANATIONS
=====================

**Concurrency Features:**
- **Goroutines**: Enable concurrent execution without the complexity of threads. The Go runtime multiplexes goroutines onto OS threads efficiently.
- **Channels**: Provide safe communication between goroutines. Buffered channels can hold multiple values, unbuffered require synchronous communication.
- **Semaphores**: Implemented with buffered channels to limit concurrent access to resources, preventing system overload.
- **WaitGroups**: Ensure all goroutines complete before proceeding, crucial for synchronization in batch operations.
- **Mutexes**: Prevent race conditions when multiple goroutines access shared data simultaneously.
- **defer**: Guarantees cleanup code runs regardless of how a function exits, essential for resource management.
- **Context**: Enables cancellation and timeout propagation through call chains, improving reliability.
- **Atomic Operations**: Provide lock-free thread-safe operations for simple counters and flags.

**Data Structures:**
- **Structs**: Form the backbone of Go programs, grouping related data fields together.
- **Maps**: Provide O(1) average-time complexity lookups, ideal for caching and configuration.
- **Slices**: Flexible arrays that grow automatically, used extensively for collections.

**Control Flow:**
- **For Loops**: Go's only looping construct, versatile for different iteration patterns.
- **If Statements**: Standard conditional execution, often combined with short variable declarations.
- **Switch**: Cleaner than multiple if-else chains, supports type switching for interfaces.

**Functions:**
- **Regular Functions**: Basic building blocks of Go programs.
- **Methods**: Associate functions with types, enabling object-oriented patterns.
- **Anonymous Functions**: Create closures, useful for goroutines and deferred execution.

**Type System:**
- **Constants**: Compile-time evaluated values, improving performance and preventing accidental changes.
- **Pointers**: Allow efficient passing of large structures and enable mutation of function parameters.
- **Interfaces**: Enable polymorphism and dependency injection patterns.

**Built-in Functions:**
- **make()**: Essential for creating reference types (slices, maps, channels) with proper initialization.
- **append()**: Core slice manipulation, handles growth and reallocation automatically.
- **len()**/**cap()**: Provide metadata about data structures for bounds checking and optimization.

**Standard Library Usage:**
- **String Operations**: Essential for text processing and file path manipulation.
- **Regular Expressions**: Powerful pattern matching for secret detection and validation.
- **Time Operations**: Critical for logging, timeouts, and scheduling.
- **UUID Generation**: Creates unique identifiers for tracking requests and resources.
- **Logging**: Structured logging with levels and JSON output for production monitoring.
- **File I/O**: Basic file operations for configuration and data persistence.
- **Environment Variables**: Configuration without hardcoding sensitive values.
- **JSON**: Standard data interchange format, ubiquitous in web APIs.
- **HTTP**: Foundation of web services and API communication.
- **Templates**: Server-side HTML generation for web interfaces.        