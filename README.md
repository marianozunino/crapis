# Crapis - RESP (REdis Serialization Protocol) Implementation

This project is an implementation of the RESP (REdis Serialization Protocol) in Go, based on the tutorial from [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/resp-writer). It extends the basic functionality with additional features like TTL (Time To Live) support and key eviction strategies.

## Overview

RESP is a simple, binary-safe protocol used by Redis for client-server communication. This implementation provides a structure for reading and writing RESP-encoded data, handling basic Redis commands, and managing data expiration.

## Features

- Parsing of RESP data types: Simple Strings, Errors, Integers, Bulk Strings, Arrays
- RESP writer functionality for encoding responses
- Error handling for invalid inputs
- CLI flags for server configuration
- TTL support for key expiration
- Active and passive eviction strategies
- Basic Redis commands:
  - GET: Retrieve the value of a key
  - SET: Set the value of a key
  - SETEX: Set the value and expiration of a key
  - DEL: Delete one or more keys
  - PING: Test if the server is responsive
  - EXPIRE: Set the expiration time of a key

## Eviction Strategies

Crapis implements two eviction strategies for handling key expiration:

1. **Passive Eviction**: When a GET operation is performed on a key, the system checks if the key has expired. If it has, the key is deleted, and a nil value is returned.

2. **Active Eviction**: A background goroutine runs periodically to check for and remove expired keys. This process helps to free up memory proactively, rather than waiting for keys to be accessed.

### How Eviction Works

- **TTL Registry**: The system maintains a separate map (`ttlKeys`) to keep track of keys with TTL. This optimization allows for efficient checking of keys that may have expired.

- **Expiration Check**:
  - Passive: Performed during GET operations.
  - Active: A goroutine runs every 250 milliseconds to check for expired keys.

- **Deletion Process**:
  - The active eviction process uses a two-pass approach to minimize lock contention:
    1. It first scans for expired keys using a read lock.
    2. Then it deletes the expired keys, acquiring a write lock only for the actual deletion.

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/marianozunino/crapis.git
   ```
2. Navigate to the project directory:
   ```
   cd crapis
   ```

### Usage

To interact with the server, you can use `redis-cli`:

```bash
redis-cli -p 6379
```

Then you can use the implemented commands:

```
> SET mykey "Hello"
OK
> GET mykey
"Hello"
> SETEX tempkey 10 "This will expire"
OK
> PING
PONG
> DEL mykey tempkey
(integer) 2
```

### Running the Server

To run the server, use the `go run .` command. You can configure the server using the following CLI flags:

```
Usage:
  crapis [flags]

Flags:
  -f, --aof string                 Path to AOF file (default "database.aof")
  -a, --aof-enabled                Enable AOF
  -b, --bind string                Bind address (default "0.0.0.0")
  -d, --debug                      Enable debug mode
  -i, --eviction-interval-ms int   Eviction interval in milliseconds (default 250)
  -t, --eviction-timeout-ms int    Eviction timeout in milliseconds, must be at at most half of eviction-interval-ms (default 10)
  -h, --help                       help for crapis
  -e, --passive-eviction           Enable passive eviction (default true)
  -p, --port string                Port to listen on (default "6379")
```

For example, to run the server on a different port with debug mode enabled:

```bash
go run . -p 7000 -d
```

This will start the server on port 7000 with debug logging enabled.

## Testing

To run the tests for this project, use the following command:

```
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/) for the tutorial and inspiration.

