# Azin

A modern, high-performance systems programming language engineered for structural clarity, low-level execution control, and human readability.

[![CI](https://github.com/azin-lang/azin/actions/workflows/build.yml/badge.svg)](https://github.com/azin-lang/azin/actions)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Fazin-lang.github.io%2Fazin%2F&label=website&color=7289da)](https://azin-lang.org/)
[![Documentation](https://img.shields.io/website?url=https%3A%2F%2Fazin-lang.github.io%2Fazin%2F&label=documentation&color=7289da)](https://docs.azin-lang.org/)

---

# Key design philosophies

Azin bridges the gap between low-level, performance-critical systems programming and a modern, expressive syntax designed for readability and maintainability.

- **Explicit block scoping:** Replaces traditional brace nesting (`{}`) with clean `do` / `end` blocks.
- **Static typing:** Compiler-enforced type safety with no runtime type checks.
- **Minimalist punctuation:** Eliminates unnecessary syntax where program structure is already explicit.
- **Systems-first:** Designed for direct native compilation without relying on a heavyweight runtime or garbage collector.

---

## Language Showcase

Here is a quick glance at writing code in Azin:

### Structs and Pure Functions
```az
importc "stdio"

struct Point is
    mut x: int
    mut y: int
end

fn main do
    var mut p1: Point
    p1.x = 10; p1.y = 20
    
    var mut p2: Point
    if p1.x <= p1.y then
        p2.x = 30
    else
        p2.x = 40
    end
    p2.y = 5
       
    printf("(%d, %d) (%d, %d)\n", p1.x, p1.y, p2.x, p2.y);
end
```

## Documentation

The complete language reference, compiler internals, and API documentation are available [here](https://docs.azin-lang.org/).

## Building

Azin requires Go. See go.mod for the minimum supported Go version.

Clone the repository and build the compiler:

```
git clone https://github.com/azin-lang/Azin.git
cd Azin
go build -o azc ./cmd/azc
```

Alternatively, install the latest development version:
```
go install github.com/azin-lang/Azin/cmd/azc@latest
```

## Usage

Compile an Azin source file:

```
azc hello.az
```

Display the available commands:
```
azc --help
```


## Development

Run the test suite:
```
go test ./...
```
Run the linter:
```
golangci-lint run
```
Format all Go source files:
```
gofmt -w .
```

### Contributing

Before contributing, configure Git to use the repository's shared hooks:
```
git config core.hooksPath .githooks
```
The pre-commit hook formats staged Go files automatically using goimports.

Install goimports:
```
go install golang.org/x/tools/cmd/goimports@latest
```
Before opening a pull request, ensure the project builds cleanly:
```
go test ./...
golangci-lint run
```
For more information about the language, compiler internals, and APIs, visit https://docs.azin-lang.org/. Also see our versioning policy over [here](./VERSIONING.md).
