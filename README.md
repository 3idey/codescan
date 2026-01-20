# codescan

[![Go Reference](https://pkg.go.dev/badge/github.com/3idey/codescan.svg)](https://pkg.go.dev/github.com/3idey/codescan)
[![Go Report Card](https://goreportcard.com/badge/github.com/3idey/codescan)](https://goreportcard.com/report/github.com/3idey/codescan)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go library and CLI tool that scans annotated Go source code and generates OpenAPI 2.0 (Swagger) specifications.

Originally part of [go-swagger](https://github.com/go-swagger/go-swagger), now available as a standalone library for easier integration and faster release cycles.

## Features

- 📝 Scan Go code annotations to generate Swagger specs
- 🔍 Support for Go modules
- 🏷️ Build tags support
- 📦 Available as both library and CLI
- 🐳 Docker image available

## Installation

### As a Library

```bash
go get github.com/3idey/codescan
```

### As a CLI Tool

```bash
go install github.com/3idey/codescan/cmd/codescan@latest
```

### Using Docker

```bash
docker pull ghcr.io/3idey/codescan:latest

# Scan current directory
docker run --rm -v $(pwd):/workspace ghcr.io/3idey/codescan generate ./...
```

## Usage

### Library Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/3idey/codescan/codescan"
)

func main() {
    opts := &codescan.Options{
        Packages:   []string{"./..."},
        ScanModels: true,
    }

    spec, err := codescan.Run(opts)
    if err != nil {
        log.Fatal(err)
    }

    output, _ := json.MarshalIndent(spec, "", "  ")
    fmt.Println(string(output))
}
```

### CLI Usage

```bash
# Generate spec for current package (JSON to stdout)
codescan generate ./...

# Generate spec with YAML output to file
codescan generate -o api.yaml --format yaml ./...

# Generate spec with build tags
codescan generate --tags=integration ./cmd/server

# Include only specific patterns
codescan generate --include="api/*" ./...

# Merge with existing spec
codescan generate -i base-spec.json ./...
```

### CLI Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file (default: stdout) |
| `--format` | Output format: `json` or `yaml` (default: json) |
| `-w, --work-dir` | Working directory for package resolution |
| `--tags` | Build tags to use when scanning |
| `--scan-models` | Include models not referenced by operations |
| `--exclude-deps` | Exclude dependencies from scanning |
| `--include` | Patterns to include |
| `--exclude` | Patterns to exclude |
| `--include-tags` | Tags to include |
| `--exclude-tags` | Tags to exclude |
| `-i, --input` | Input swagger spec to merge with |
| `--x-nullable-pointers` | Set x-nullable for pointer types |
| `--ref-aliases` | Use $ref for type aliases |
| `--transparent-aliases` | Make type aliases completely transparent |
| `--desc-with-ref` | Allow descriptions together with $ref |
| `--compact` | Produce compact JSON output |

## Configuration Options

The `codescan.Options` struct provides the following configuration:

```go
type Options struct {
    // Packages to scan (e.g., "./...", "./cmd/api")
    Packages []string
    
    // InputSpec is an existing spec to merge with
    InputSpec *spec.Swagger
    
    // ScanModels includes models not referenced by operations
    ScanModels bool
    
    // WorkDir is the working directory for package resolution
    WorkDir string
    
    // BuildTags specifies build tags to use
    BuildTags string
    
    // ExcludeDeps excludes dependencies from scanning
    ExcludeDeps bool
    
    // Include patterns
    Include []string
    
    // Exclude patterns
    Exclude []string
    
    // IncludeTags filters operations by tags
    IncludeTags []string
    
    // ExcludeTags excludes operations by tags
    ExcludeTags []string
    
    // SetXNullableForPointers adds x-nullable for pointer types
    SetXNullableForPointers bool
    
    // RefAliases uses $ref for type aliases
    RefAliases bool
    
    // TransparentAliases makes aliases completely transparent
    TransparentAliases bool
    
    // DescWithRef allows descriptions with $ref
    DescWithRef bool
}
```

## Annotations

codescan recognizes swagger annotations in Go comments. See the [go-swagger documentation](https://goswagger.io/use/spec.html) for a complete guide on annotation syntax.

### Quick Examples

#### API Info (in main package)

```go
// Package classification API.
//
// Documentation of the API.
//
//	Schemes: https
//	Host: localhost
//	BasePath: /api/v1
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main
```

#### Route/Operation

```go
// swagger:route GET /users users listUsers
//
// Lists all users.
//
// This will show all available users.
//
//	Responses:
//	  200: usersResponse
//	  default: errorResponse
```

#### Model

```go
// User represents a user in the system.
//
// swagger:model
type User struct {
    // The unique identifier
    // required: true
    // example: 123
    ID int64 `json:"id"`
    
    // The user's name
    // required: true
    // min length: 1
    // max length: 100
    Name string `json:"name"`
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This package was extracted from [go-swagger/go-swagger](https://github.com/go-swagger/go-swagger) to provide a standalone spec generation library.
