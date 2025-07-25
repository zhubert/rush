# VM Logging System

The Rush Virtual Machine includes a comprehensive logging system for debugging and monitoring execution.

## Log Levels

- **none**: No logging (default)
- **error**: Only errors and critical issues
- **warn**: Warnings and errors
- **info**: General information, initialization, and execution statistics
- **debug**: Detailed execution flow including operations and values
- **trace**: Extremely verbose output including every instruction and stack operation

## Usage

### Command Line

Use the `-log-level` flag when running Rush in bytecode mode:

```bash
# Basic info logging
rush -bytecode -log-level=info program.rush

# Detailed debugging
rush -bytecode -log-level=debug program.rush

# Full execution tracing
rush -bytecode -log-level=trace program.rush
```

### Programmatic Usage

```go
// Create VM with logging
vm := vm.NewWithLogger(bytecode, vm.LogDebug)

// Change log level during execution
vm.SetLogLevel(vm.LogTrace)

// Get execution statistics
stats := vm.GetStats()
```

## Log Output

Logs are written to stderr with timestamps and include:

- **INFO**: VM initialization, execution start/end, performance statistics
- **DEBUG**: Instruction execution, constants loaded, operations performed
- **TRACE**: Every stack push/pop, instruction pointer movement, frame operations
- **ERROR**: Runtime errors, stack overflow/underflow, operation failures

## Performance Statistics

When using info level or higher, the VM automatically reports:

- Execution time
- Instructions executed
- Stack operations performed
- Function calls made
- Memory allocations
- Errors encountered
- Instructions per second

## Example Output

```
[VM] INFO: VM initialized with 7 constants, 2048 stack size, 65536 globals size
[VM] DEBUG: Main function has 28 instructions
[VM] DEBUG: Loading constant[0]: 5
[VM] DEBUG: Executing binary operation: OpAdd
[VM] TRACE: IP:6 OP:OpAdd SP:2 Frame:0
[VM] INFO: === VM EXECUTION STATS ===
[VM] INFO: Execution time: 173.084µs
[VM] INFO: Instructions executed: 14
[VM] INFO: Instructions per second: 80885.58
```

## Implementation Details

The logging system uses Go's standard `log` package with configurable levels. All logging is thread-safe and includes microsecond timestamps for performance analysis.