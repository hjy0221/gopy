# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

gopy runs a small subset of Python directly in Go. Running it requires Go 1.22 or later. CPython is optional.

## Quick start

Run these commands to clone the repository and execute the included example.

```sh
git clone https://github.com/hjy0221/gopy.git
cd gopy
go run ./cmd/gopy examples/fibonacci.py
```

The example calls `print(fibonacci(10))` and prints `55`.

## Command options

| Command | What it does |
| --- | --- |
| `go run ./cmd/gopy examples/fibonacci.py` | Runs the Python program and shows its output. |
| `go run ./cmd/gopy --tokens examples/fibonacci.py` | Shows the lexer output. Each row gives the token’s line and column, type, and text. For example, `print` is an `IDENTIFIER` and `10` is a `NUMBER`. |
| `go run ./cmd/gopy --ast examples/fibonacci.py` | Shows the abstract syntax tree (AST) as JSON. Use it to see how expressions, calls, and blocks fit together. |
| `go run ./cmd/gopy --help` | Shows the available options and exits. |

Put options before the filename. Use either `--tokens` or `--ast` at a time. Both inspection modes display source structure without running the program.

## Read from standard input

Omit the filename to read Python source from standard input. This example prints `7`. You can also use standard input with either inspection option.

```sh
printf 'print(1 + 2 * 3)\n' | go run ./cmd/gopy
```

## Supported Python features

Execution supports simple assignments, top-level functions with positional arguments, recursion, `return`, `if/else`, 64-bit integer arithmetic and comparisons, and `print`. The parser can also show string and floating-point literals and chained comparisons in the AST. These features, true division (`/`), and nested functions are outside the current execution subset.

## How it works

The lexer splits source into tokens. The parser groups those tokens into an AST. The compiler turns the AST into bytecode, and the stack VM runs it.

- `cmd/gopy/` — CLI and options
- `internal/lexer/` — Lexer
- `internal/parser/ + internal/ast/` — Parser and syntax tree
- `internal/compiler/ + internal/bytecode/` — Compiler and bytecode definitions
- `internal/vm/` — Stack VM
- `examples/fibonacci.py` — Example that prints 55

## Build and test

```sh
go build -o gopy ./cmd/gopy
./gopy examples/fibonacci.py
go test ./...
```
