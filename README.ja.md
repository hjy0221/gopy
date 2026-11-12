# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

`gopy` は Go で書かれた独立した Python 実装です。CPython を呼び出さずに Python ソースを直接実行することを目指します。

実行の流れは、字句解析 → 構文解析と AST → バイトコードコンパイラ → スタック型 VM → Python オブジェクトランタイムと組み込み関数です。

## 対応範囲

字句解析器は、識別子、数値、引用符で囲まれた文字列、演算子と区切り記号、コメント、論理改行、インデントとその解除、ファイル終端を処理します。CLI の `-tokens` でトークン列を表示できます。構文解析器は、単純な代入、関数定義、return 文、if/else ブロック、関数呼び出し、リテラル、算術演算と比較を扱います。`--ast` で構文木を JSON として表示できます。オプションなしでソースファイルを指定すると、バイトコードにコンパイルして Go VM で実行します。フィボナッチの例は `55` を出力します。現在は 64 ビット整数の演算と比較、基本的な条件分岐、位置引数を取るトップレベルの関数、`print` を実行できます。文字列・浮動小数点リテラルと連鎖比較は構文解析で確認できます。

## ビルドと実行

```sh
go test ./...
go run ./cmd/gopy examples/fibonacci.py
go run ./cmd/gopy --ast examples/fibonacci.py
go run ./cmd/gopy -tokens example.py
cat example.py | go run ./cmd/gopy -tokens
```

## プロジェクト構成

```text
cmd/gopy/              CLI
internal/lexer/        字句解析器
internal/parser/       構文解析器
internal/ast/          構文木
internal/compiler/     バイトコードコンパイラ
internal/bytecode/     命令の定義
internal/vm/           スタック型 VM
internal/runtime/      Python オブジェクトモデル (予定)
internal/builtins/     組み込み関数 (予定)
tests/                 統合テスト (予定)
```

サンプルは `print(fibonacci(10))` を呼び出します。
