# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

gopy は Python の一部を Go で直接実行する処理系です。Go 1.22 以降が必要です。CPython のインストールは不要です。

## クイックスタート

次のコマンドでリポジトリを取得し、サンプルを実行します。

```sh
git clone https://github.com/hjy0221/gopy.git
cd gopy
go run ./cmd/gopy examples/fibonacci.py
```

サンプルは `print(fibonacci(10))` を呼び出し、`55` を表示します。

## 実行オプション

| コマンド | 動作 |
| --- | --- |
| `go run ./cmd/gopy examples/fibonacci.py` | Python プログラムを実行し、結果を表示します。 |
| `go run ./cmd/gopy --tokens examples/fibonacci.py` | 字句解析で得られたトークンを表示します。各行にはトークンの行番号、列番号、種類、元の文字列が表示されます。たとえば `print` は `IDENTIFIER`、`10` は `NUMBER` です。 |
| `go run ./cmd/gopy --ast examples/fibonacci.py` | 抽象構文木 (AST) を JSON で表示します。式、関数呼び出し、インデントで区切られたブロックの構造を確認できます。 |
| `go run ./cmd/gopy --help` | 利用できるオプションを表示して終了します。 |

オプションはファイル名の前に指定します。`--tokens` と `--ast` はどちらか一方を指定してください。どちらもプログラムを実行せずにソースの構造を表示します。

## 標準入力からの実行

ファイル名を省略すると、標準入力から Python ソースを読み込みます。次の例は `7` を表示します。`--tokens` と `--ast` でも標準入力を利用できます。

```sh
printf 'print(1 + 2 * 3)\n' | go run ./cmd/gopy
```

## 対応する Python の機能

単純な代入、位置引数を取るトップレベルの関数、再帰呼び出し、`return`、`if/else`、64 ビット整数の演算と比較、`print` を実行できます。文字列・浮動小数点リテラルと連鎖比較は AST で確認できます。これらの構文、通常の除算 (`/`)、入れ子の関数の実行にはまだ対応していません。

## 仕組み

字句解析器がソースをトークンに分割し、構文解析器がトークンから AST を作ります。コンパイラが AST をバイトコードに変換し、スタック型 VM が実行します。

- `cmd/gopy/` — CLI と実行オプション
- `internal/lexer/` — 字句解析器
- `internal/parser/ + internal/ast/` — 構文解析器と構文木
- `internal/compiler/ + internal/bytecode/` — コンパイラとバイトコードの定義
- `internal/vm/` — スタック型 VM
- `examples/fibonacci.py` — 55 を表示するサンプル

## ビルドとテスト

```sh
go build -o gopy ./cmd/gopy
./gopy examples/fibonacci.py
go test ./...
```
