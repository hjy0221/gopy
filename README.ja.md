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

## gopython をインストールして実行する

上で取得した gopy フォルダーで、次の順に実行します。

1. 実行ファイルを作ります。1 行目でコマンドを置くフォルダーを作り、2 行目でそこに gopython を保存します。

   ```sh
   mkdir -p "$HOME/.local/bin"
   go build -o "$HOME/.local/bin/gopython" ./cmd/gopy
   ```

   - `mkdir` はディレクトリを作ります。`-p` は不足している親ディレクトリも作り、対象がすでにあればそのまま使います。
   - `go build` は Go のソースから実行ファイルを作ります。`-o` で出力先とファイル名を指定します。`./cmd/gopy` は現在のディレクトリから見た CLI パッケージの場所です。

2. 現在のターミナルで gopython を見つけられるようにします。

   ```sh
   export PATH="$HOME/.local/bin:$PATH"
   ```

3. サンプルを実行します。55 が表示されればインストールは完了です。

   ```sh
   gopython examples/fibonacci.py
   ```

   - `gopython` はインストールした実行ファイルを起動します。`examples/fibonacci.py` は現在のディレクトリから見た、実行する Python ファイルのパスです。

   ```text
   55
   ```

新しいターミナルでも使うには、手順 2 の export 行をシェルの設定ファイルに追加します。zsh では ~/.zshrc に追加してください。「command not found」が表示されたら、そのターミナルで手順 2 をもう一度実行します。

### PATH コマンドの各部分

- `PATH=` シェルがコマンドを探すディレクトリの一覧を設定します。
- `$` `$変数名` の `$` は、その変数の値に置き換えるようシェルに指示する記号です。たとえば `HOME` の値が `/Users/alex` なら、`$HOME/.local/bin` は `/Users/alex/.local/bin` になります。
- `$HOME` ホームディレクトリに展開されます。`$HOME/.local/bin` は手順 1 で `gopython` を置いた場所です。
- `:` 検索するディレクトリ同士を区切ります。
- `$PATH` 変更前の検索一覧に展開されます。既存の場所も引き続き検索できます。
- `export` 変更した変数をシェルから起動するプログラムにも渡します。
- `"..."` ディレクトリ名に空白があっても、展開された値をひとまとまりにします。
- `~/.zshrc` 対話型 zsh の起動時に読み込まれる設定ファイルです。`~` はホームディレクトリを指します。

## 実行オプション

| インストール後 | コマンド (`go run`) | 動作 |
| --- | --- | --- |
| `gopython examples/fibonacci.py` | `go run ./cmd/gopy examples/fibonacci.py` | Python プログラムを実行し、結果を表示します。 |
| `gopython --tokens examples/fibonacci.py` | `go run ./cmd/gopy --tokens examples/fibonacci.py` | 字句解析で得られたトークンを表示します。各行にはトークンの行番号、列番号、種類、元の文字列が表示されます。たとえば `print` は `IDENTIFIER`、`10` は `NUMBER` です。 |
| `gopython --ast examples/fibonacci.py` | `go run ./cmd/gopy --ast examples/fibonacci.py` | 抽象構文木 (AST) を JSON で表示します。式、関数呼び出し、インデントで区切られたブロックの構造を確認できます。 |
| `gopython --help` | `go run ./cmd/gopy --help` | 利用できるオプションを表示して終了します。 |

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
