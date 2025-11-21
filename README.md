# minid

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peteraba/minid) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/peteraba/minid/main/LICENSE) 

A reasonably fast, (optionally) somewhat sortable, URL-safe and human-friendly ID generator for Go, which can be used either as a CLI tool or a library.

Generate fully random IDs or time-based IDs with configurable precision (seconds, milliseconds, microseconds, or nanoseconds).

It can generate a number of IDs and the ones generated in one go will be unique, but there's no generic uniquness guarantees.

It's best used for pre-generating a large number of IDs, for example to be used in spreadsheets such as Google Spreadsheets.

That said, usually you want to use other tools such as [xid](https://github.com/rs/xid), [snowflake](https://github.com/bwmarrin/snowflake), [sonyflake](https://github.com/sony/sonyflake) or one of the [hundreds of others](https://github.com/topics/id-generator?l=go)

## Features

- **Multiple ID Types**: Pure random IDs or time-based IDs with different precisions
- **Sortable**: Time-based IDs are lexicographically sortable
- **URL-Safe**: Uses alphanumeric characters only, no special characters, and no confusing characters: 0, O, l, etc.
- **Collision-Free**: Built-in duplicate detection ensures unique IDs in one go.
- **No Leading Numbers**: Random suffixes never start with 0, to prevent unintentional removal in some tools.
- **CLI Tool**: Easy-to-use command-line interface
- **Library**: Use as a Go package in your applications

## ID Format

Minid comes with a highly customizable format. By default it will generate a random 4-letter string. It can also prefix the random letters by a number of strings which represent time elapsed since 2025-01-01. For now this is hard coded, but could be easily modified if needed.

Furthermore you can pick the kind of time you want to use which could have a big impact on sortability.

Options are:
- r: no sortabiliy, no timing.
- s: seconds
- m: milliseconds - This would allow you to generate dozens of IDs in a second and assume them to be more or less fully sortable.
- u: microseconds - This would allow you to generate thousands of IDs in a second and assume them to be fully sortable
- n: nanoseconds - With this option you can assume that all your IDs generated on a node to fully sortable as generating a minid takes much longer.

Note: (Partially) sortable IDs change the default for the number of random characters to optimize for shorter IDs.

## Database storage

Minid IDs can be stored as byte slices or strings, up to you.


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

