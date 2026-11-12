# gopy

[English](README.md) · [한국어](README.ko.md) · [日本語](README.ja.md)

gopy는 Python 문법의 일부를 Go로 직접 실행합니다. Go 1.22 이상이 필요하며 CPython 설치는 필요 없습니다.

## 빠른 시작

아래 명령으로 저장소를 내려받고 예제를 실행하세요.

```sh
git clone https://github.com/hjy0221/gopy.git
cd gopy
go run ./cmd/gopy examples/fibonacci.py
```

예제는 `print(fibonacci(10))`을 호출해 `55`를 출력합니다.

## gopython 설치하고 실행하기

위에서 내려받은 gopy 폴더에서 아래 순서대로 실행하세요.

1. 실행 파일을 만듭니다. 첫 줄은 명령을 보관할 폴더를 만들고, 둘째 줄은 그 안에 gopython을 저장합니다.

   ```sh
   mkdir -p "$HOME/.local/bin"
   go build -o "$HOME/.local/bin/gopython" ./cmd/gopy
   ```

   - `mkdir`는 폴더를 만듭니다. `-p`는 필요한 상위 폴더도 만들고, 대상 폴더가 이미 있으면 그대로 사용합니다.
   - `go build`는 Go 소스를 실행 파일로 빌드합니다. `-o` 뒤에는 저장할 경로와 파일 이름을 지정합니다. `./cmd/gopy`는 현재 폴더를 기준으로 CLI 소스가 있는 곳입니다.

2. 현재 터미널에 gopython이 있는 위치를 알려 줍니다.

   ```sh
   export PATH="$HOME/.local/bin:$PATH"
   ```

3. 예제를 실행합니다. 화면에 55가 나오면 설치가 완료된 것입니다.

   ```sh
   gopython examples/fibonacci.py
   ```

   - `gopython`은 설치한 실행 파일을 시작합니다. `examples/fibonacci.py`는 현재 폴더를 기준으로 실행할 Python 파일의 경로입니다.

   ```text
   55
   ```

새 터미널에서도 쓰려면 2단계의 export 줄을 셸 설정 파일에 추가하세요. zsh를 사용한다면 ~/.zshrc에 추가하면 됩니다. “command not found”가 나오면 해당 터미널에서 2단계를 다시 실행하세요.

### PATH 명령의 각 부분

- `PATH=` 셸이 명령을 찾을 때 사용하는 폴더 목록을 설정합니다.
- `$` `$변수이름`에서 `$`는 셸에 그 변수의 값을 넣으라는 표시입니다. 예를 들어 `HOME`의 값이 `/Users/alex`이면 `$HOME/.local/bin`은 `/Users/alex/.local/bin`으로 바뀝니다.
- `$HOME` 사용자 홈 폴더로 바뀝니다. `$HOME/.local/bin`은 1단계에서 `gopython`을 설치한 폴더입니다.
- `:` 검색 목록에서 폴더와 폴더를 구분합니다.
- `$PATH` 설정하기 전의 검색 목록으로 바뀝니다. 기존 폴더에서도 계속 명령을 찾을 수 있습니다.
- `export` 변경한 변수를 이 셸에서 실행하는 프로그램에도 전달합니다.
- `"..."` 폴더 이름에 공백이 있어도 펼쳐진 값 전체를 하나로 묶습니다.
- `~/.zshrc` 대화형 zsh가 시작할 때 읽는 설정 파일입니다. `~`는 사용자 홈 폴더를 뜻합니다.

## 실행 옵션

| 설치 후 | 명령 (`go run`) | 동작 |
| --- | --- | --- |
| `gopython examples/fibonacci.py` | `go run ./cmd/gopy examples/fibonacci.py` | Python 코드를 실행하고 프로그램의 출력을 보여 줍니다. |
| `gopython --tokens examples/fibonacci.py` | `go run ./cmd/gopy --tokens examples/fibonacci.py` | 소스를 나눈 토큰을 보여 줍니다. 각 줄에는 토큰의 줄·열 번호, 종류, 원문이 나옵니다. 예를 들어 `print`는 `IDENTIFIER`, `10`은 `NUMBER`입니다. |
| `gopython --ast examples/fibonacci.py` | `go run ./cmd/gopy --ast examples/fibonacci.py` | 추상 구문 트리(AST)를 JSON으로 보여 줍니다. 식과 함수 호출, 들여쓰기 블록이 어떻게 묶이는지 확인할 수 있습니다. |
| `gopython --help` | `go run ./cmd/gopy --help` | 사용 가능한 옵션을 표시하고 종료합니다. |

옵션은 파일 이름 앞에 적으세요. `--tokens`와 `--ast`는 한 번에 하나만 사용할 수 있습니다. 두 옵션 모두 프로그램을 실행하지 않고 코드의 구조를 보여 줍니다.

## 표준 입력으로 실행

파일 이름을 생략하면 표준 입력에서 Python 소스를 읽습니다. 아래 예제는 `7`을 출력합니다. `--tokens`와 `--ast`도 표준 입력에서 사용할 수 있습니다.

```sh
printf 'print(1 + 2 * 3)\n' | go run ./cmd/gopy
```

여기서 `printf`는 zsh에 내장된 셸 명령입니다. 따옴표 안의 글자를 출력하고, `\n`은 줄바꿈을 넣습니다. 파이프(`|`)가 그 글자를 gopy에 넘기면, gopy가 Python 코드를 실행해 `7`을 출력합니다.

## 지원하는 Python 문법

단순 대입, 위치 인수를 받는 최상위 함수, 재귀 호출, `return`, `if/else`, 64비트 정수 연산과 비교, `print`를 실행할 수 있습니다. 문자열·실수 리터럴과 연쇄 비교는 AST에서 확인할 수 있습니다. 해당 문법과 일반 나눗셈(`/`), 중첩 함수의 실행은 아직 지원하지 않습니다.

## 동작 원리

Lexer는 소스를 토큰으로 나눕니다. Parser는 토큰을 AST로 묶습니다. 컴파일러는 AST를 바이트코드로 바꾸고, 스택 VM이 바이트코드를 실행합니다.

- `cmd/gopy/` — CLI와 실행 옵션
- `internal/lexer/` — 토큰 분석기
- `internal/parser/ + internal/ast/` — 파서와 구문 트리
- `internal/compiler/ + internal/bytecode/` — 컴파일러와 바이트코드 정의
- `internal/vm/` — 스택 VM
- `examples/fibonacci.py` — 55를 출력하는 예제

## 다른 실행 예제

저장소 폴더에서 아래 명령을 실행하세요. 소스에서 바로 실행하려면 `gopython` 자리에 `go run ./cmd/gopy`를 넣으세요.

| 실행 명령 | 예제 내용 | 출력 |
| --- | --- | --- |
| `gopython examples/arithmetic.py` | 정수 연산, 몫과 나머지 | `22 12 85`<br>`3 2` |
| `gopython examples/conditionals.py` | if/else로 짝수와 홀수 판별 | `True False` |
| `gopython examples/functions.py` | 함수 안에서 다른 함수 호출 | `25` |

## 빌드와 테스트

```sh
go build -o gopy ./cmd/gopy
./gopy examples/fibonacci.py
go test ./...
```
