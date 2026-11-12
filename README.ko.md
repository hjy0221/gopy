# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

`gopy`는 Go로 작성한 독립 Python 구현체입니다. CPython을 호출하지 않고 Python 소스를 직접 실행하는 것을 목표로 합니다.

실행 흐름은 Lexer → Parser/AST → 바이트코드 컴파일러 → 스택 기반 VM → Python 객체 런타임 및 내장 함수 순서입니다.

## 현재 개발 단계

현재 Lexer는 식별자, 숫자와 따옴표로 감싼 문자열, 연산자와 구분 기호, 주석, 논리적 줄바꿈, 들여쓰기와 들여쓰기 해제, 파일 끝을 처리합니다. CLI에서 `-tokens`로 토큰을 출력할 수 있습니다. Parser는 단순 대입, 함수 정의, return 문, if/else 블록, 함수 호출, 리터럴, 산술과 비교를 지원합니다. `--ast`는 구문 트리를 JSON으로 출력합니다. 옵션 없이 소스 파일을 지정하면 바이트코드로 컴파일해 Go VM에서 실행합니다. 피보나치 예제는 `55`를 출력합니다. 현재 실행 범위는 64비트 정수 연산과 비교, 기본 조건문, 위치 인수를 받는 최상위 함수, `print`입니다. 문자열·실수 리터럴과 연쇄 비교는 구문 분석 단계에서 확인할 수 있습니다.

## 빌드 및 실행

```sh
go test ./...
go run ./cmd/gopy examples/fibonacci.py
go run ./cmd/gopy --ast examples/fibonacci.py
go run ./cmd/gopy -tokens example.py
cat example.py | go run ./cmd/gopy -tokens
```

## 프로젝트 구조

```text
cmd/gopy/              CLI
internal/lexer/        Lexer
internal/parser/       파서
internal/ast/          구문 트리
internal/compiler/     바이트코드 컴파일러
internal/bytecode/     명령어 정의
internal/vm/           스택 기반 VM
internal/runtime/      Python 객체 모델 (예정)
internal/builtins/     내장 함수 (예정)
tests/                 통합 테스트 (예정)
```

예제는 `print(fibonacci(10))`을 호출합니다.
