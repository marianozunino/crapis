# Crapis - RESP (REdis Serialization Protocol) Implementation

This project is a basic implementation of the RESP (REdis Serialization Protocol) in Go.
It's based on the tutorial from [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/resp-writer).

## Overview

RESP is a simple, binary-safe protocol used by Redis for client-server communication. This implementation provides a basic structure for reading RESP-encoded data.

## Features

- Parsing of RESP data types:
  - Simple Strings
  - Errors
  - Integers
  - Bulk Strings
  - Arrays
- Error handling for invalid inputs
- CLI flags for server configuration

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

Here's a basic example of how to use the RESP reader:

```go
package main

import (
	"fmt"
	"strings"
	"github.com/marianozunino/crapis/internal/resp"
)

func main() {
	// Create a new RESP reader
	reader := resp.NewReader(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))

	// Read RESP data
	value, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading RESP data: %v\n", err)
		return
	}

	// Print the parsed data
	fmt.Printf("Parsed RESP data: %+v\n", value)
}
```

Or better yet, use `redis-cli` in conjunction with the `go run .` command:

```bash
redis-cli -r "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
```

or use it in interactive mode.

### Running the Server

To run the server, use the `go run .` command. You can configure the server using the following CLI flags:

```
Usage:
  crapis [flags]

Flags:
  -b, --bind string   Bind address (default "0.0.0.0")
  -d, --debug         Enable debug mode
  -h, --help          help for crapis
  -p, --port string   Port to listen on (default "6379")
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

## Future Improvements

- Implement RESP writer functionality
- Add support for more complex RESP data structures
- Improve error handling and edge case coverage
- Benchmark performance and optimize where necessary
