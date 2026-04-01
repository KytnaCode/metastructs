# metastructs

Metastructs is a code generation tool that automatically generates common boilerplate methods for Go structs. It helps reduce repetitive code by generating struct-to-map converters, name methods, and partial struct definitions.

## Features

- **ToMap**: Generate methods to convert structs to `map[string]any` with customizable tags and `omitempty` support
- **Name**: Generate methods that return the struct's name as a string
- **Partial**: Generate partial struct definitions with pointer fields for optional updates

## Installation

### As a Go tool

Install metastructs in your project:

```bash
go install github.com/kytnacode/metastructs@latest
```

## Usage

Metastructs provides three main commands: `to-map`, `name`, and `partial`.

### Global Flags

- `-p, --pkg`: Package name (defaults to `$GOPACKAGE` environment variable)
- `-t, --target`: Target struct type name (required)

### to-map Command

Generate a method to convert a struct into a map.

```bash
metastructs to-map --target MyStruct --pkg mypackage
```

**Flags:**
- `-m, --method`: Output method name (default: `ToMap`)
- `-s, --tag`: Struct tag name for metadata (default: `to-map`)

**Tag Syntax:**

The `to-map` tag follows the same conventions as `encoding/json`:

```go
type User struct {
    Name        string                      // map key: "Name"
    Email       string   `to-map:"email"`   // map key: "email"
    Password    string   `to-map:"-"`       // excluded from map
    Bio         string   `to-map:",omitempty"` // only included if non-empty
    Avatar      *string  `to-map:"avatar,omitempty"` // only included if non-nil
}
```

**Generated code:**

```go
func (u User) ToMap() map[string]any {
    structMap := map[string]any{
        "Name": u.Name,
        "email": u.Email,
    }
    if u.Bio != "" {
        structMap["Bio"] = u.Bio
    }
    if u.Avatar != nil {
        structMap["avatar"] = u.Avatar
    }
    return structMap
}
```

### name Command

Generate a method that returns the struct's name.

```bash
metastructs name --target MyStruct --pkg mypackage
```

**Flags:**
- `-m, --method`: Output method name (default: `StructName`)

**Generated code:**

```go
func (m MyStruct) StructName() string {
    return "MyStruct"
}
```

### partial Command

Generate a partial struct with all fields as pointers, useful for update operations where you only want to set specific fields.

```bash
metastructs partial --target User --pkg mypackage
```

**Flags:**
- `-e, --prefix`: Struct name prefix (default: `Partial`)
- `-s, --suffix`: Struct name suffix (default: empty)
- `-n, --structname`: Custom struct name (overrides prefix/suffix)

**Example:**

Given this struct:

```go
type User struct {
    Name  string
    Email string
    Age   int
}
```

**Generated code:**

```go
type PartialUser struct {
    Name  *string
    Email *string
    Age   *int
}
```

## Usage with go generate

Add `go generate` directives to your code:

```go
package mypackage

//go:generate metastructs to-map -t User -p mypackage
//go:generate metastructs name -t User -p mypackage
//go:generate metastructs partial -t User -p mypackage

type User struct {
    Name     string `to-map:"name"`
    Email    string `to-map:"email"`
    Password string `to-map:"-"`
}
```

Then run:

```bash
go generate ./...
```

This will generate three files:
- `user_map.go` - ToMap method
- `user_name.go` - StructName method
- `user_partial.go` - PartialUser struct

## Pointer Receivers

To generate methods with pointer receivers, prefix your target with `*`:

```bash
metastructs to-map --target "*User" --pkg mypackage
```

**Generated code:**

```go
func (u *User) ToMap() map[string]any {
    // ...
}
```

## Development

### Prerequisites

- Go 1.25.4 or later
- [just](https://github.com/casey/just) (optional, for running commands)

### Available Commands

```bash
just lint      # Run linter
just format    # Format code
just check     # Run vulnerability checks
just precommit # Run all checks before committing
```

Or use standard Go commands:

```bash
go test ./...
go build
```

## License

MIT License - see [LICENSE.txt](LICENSE.txt) for details.

Copyright (c) 2026 Alejandro Paz
