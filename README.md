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

## Install and run with gopython

Run these commands from the gopy folder you cloned above.

1. Build the command once. The first line creates a folder for your commands. The second puts gopython there.

   ```sh
   mkdir -p "$HOME/.local/bin"
   go build -o "$HOME/.local/bin/gopython" ./cmd/gopy
   ```

   - `mkdir` creates a directory. `-p` creates any missing parent directories and accepts an existing destination.
   - `go build` builds an executable from Go source. `-o` sets the output path and filename. `./cmd/gopy` points to the CLI package relative to the current directory.

2. Tell your current terminal where to find gopython.

   ```sh
   export PATH="$HOME/.local/bin:$PATH"
   ```

3. Run the example. You should see 55.

   ```sh
   gopython examples/fibonacci.py
   ```

   - `gopython` starts the installed executable. `examples/fibonacci.py` is the path to the Python file, relative to the current directory.

   ```text
   55
   ```

To use the command in new terminals, add the export line from step 2 to your shell settings file. For zsh, use ~/.zshrc. If you see “command not found”, run step 2 again in that terminal.

### What the PATH command means

- `PATH=` Sets the search path the shell uses to find commands.
- `$` In `$NAME`, the dollar sign tells the shell to substitute the value of the variable named `NAME`. For example, if `HOME` is `/Users/alex`, `$HOME/.local/bin` becomes `/Users/alex/.local/bin`.
- `$HOME` Expands to your home directory. `$HOME/.local/bin` is where step 1 installs `gopython`.
- `:` Separates directories in the search list.
- `$PATH` Expands to the existing search list, so the shell keeps looking in those directories too.
- `export` Exports the updated variable to programs launched by the shell.
- `"..."` Keeps the expanded value together even if a directory name contains spaces.
- `~/.zshrc` A configuration file read when interactive zsh starts. `~` means your home directory.

## Command options

| After installation | Command (`go run`) | What it does |
| --- | --- | --- |
| `gopython examples/fibonacci.py` | `go run ./cmd/gopy examples/fibonacci.py` | Runs the Python program and shows its output. |
| `gopython --tokens examples/fibonacci.py` | `go run ./cmd/gopy --tokens examples/fibonacci.py` | Shows the lexer output. Each row gives the token’s line and column, type, and text. For example, `print` is an `IDENTIFIER` and `10` is a `NUMBER`. |
| `gopython --ast examples/fibonacci.py` | `go run ./cmd/gopy --ast examples/fibonacci.py` | Shows the abstract syntax tree (AST) as JSON. Use it to see how expressions, calls, and blocks fit together. |
| `gopython --help` | `go run ./cmd/gopy --help` | Shows the available options and exits. |

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
