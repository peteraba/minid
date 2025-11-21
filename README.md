# minid

A reasonably fast, (optionally) somewhat sortable, URL-safe and human-friendly ID generator for Go. Generate fully random IDs or time-based IDs with configurable precision (seconds, milliseconds, microseconds, or nanoseconds).

It can generate a number of IDs and the ones generated in one go will be unique, but there's no generic uniquness guarantees.

## Features

- **Multiple ID Types**: Pure random IDs or time-based IDs with different precisions
- **Sortable**: Time-based IDs are lexicographically sortable
- **URL-Safe**: Uses alphanumeric characters only, no special characters, and no confusing characters: 0, O, l, etc.
- **Collision-Free**: Built-in duplicate detection ensures unique IDs in one go.
- **No Leading Numbers**: Random suffixes never start with 0, to prevent unintentional removal in some tools.
- **CLI Tool**: Easy-to-use command-line interface
- **Library**: Use as a Go package in your applications

## Installation

### CLI Tool

```bash
go install github.com/peteraba/minid/cmd@latest
```

Or build from source:

```bash
make build
```

## Usage

### Command-Line Interface

Generate a single random ID (default):
```bash
minid
```

Generate multiple random IDs:
```bash
minid 10
```

Generate time-based IDs with different precisions:
```bash
minid s      # Unix seconds (default randLength: 3)
minid ms     # Unix milliseconds
minid us     # Unix microseconds
minid ns     # Unix nanoseconds
```

Generate multiple time-based IDs:
```bash
minid s 5    # 5 IDs with Unix seconds precision
minid ms 10  # 10 IDs with Unix milliseconds precision
```

Customize random suffix length:
```bash
minid -randLength 5    # or -rl 5
minid s -rl 6          # Unix seconds with 6-character random suffix
```

**Examples:**
```bash
$ minid
aB3x

$ minid 3
xY9z
mK2p
nL8q

$ minid s
132f3bXSZ

$ minid ms 5
11YFAUybhw7
11YFAUybzgu
11YFAUybEWK
11YFAUybF58
11YFAUybRd3
```

## Development

### Running Tests

```bash
make test
# or
go test ./...
```

### Building

```bash
make build
```

## License

See [LICENSE](LICENSE) file for details.

