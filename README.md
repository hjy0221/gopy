# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

`gopy` is an independent Python implementation written in Go. It will execute Python source directly without calling CPython.

The execution pipeline is Lexer → Parser/AST → Bytecode Compiler → Stack VM → Python object runtime and builtins.

## Current milestone

The lexer scans source text for identifiers, numeric and quoted string literals, operators and delimiters, comments, logical newlines, indentation, dedentation, and EOF. The `-tokens` CLI prints the resulting token stream. The parser supports simple assignments, function definitions, return statements, if/else blocks, calls, literals, arithmetic, and comparisons. The `--ast` option prints the syntax tree as JSON. Run a source file without flags to compile it to bytecode and execute it in the Go VM. The Fibonacci example prints `55`. Execution currently covers 64-bit integer arithmetic, integer comparisons, basic conditionals, top-level functions with positional arguments, and `print`. String and float literals and chained comparisons are available only for syntax inspection.

## Build and run

```sh
go test ./...
go run ./cmd/gopy examples/fibonacci.py
go run ./cmd/gopy --ast examples/fibonacci.py
go run ./cmd/gopy -tokens example.py
cat example.py | go run ./cmd/gopy -tokens
```

Project layout

```text
cmd/gopy/              CLI
internal/lexer/        tokenizer
internal/parser/       parser
internal/ast/          syntax tree
internal/compiler/     bytecode compiler
internal/bytecode/     instruction definitions
internal/vm/           stack-based virtual machine
internal/runtime/      Python object model (planned)
internal/builtins/     builtin functions (planned)
tests/                 integration tests (planned)
```

The included example calls `print(fibonacci(10))`.
