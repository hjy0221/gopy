package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"gopy/internal/compiler"
	"gopy/internal/lexer"
	"gopy/internal/parser"
	"gopy/internal/vm"
)

func main() {
	tokens := flag.Bool("tokens", false, "print lexer tokens")
	showAST := flag.Bool("ast", false, "print the syntax tree as JSON")
	flag.Parse()
	if *tokens && *showAST || flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}
	var source string
	if flag.NArg() > 0 {
		data, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		source = string(data)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		source = string(data)
	}
	if *showAST {
		tree, err := parser.Parse(source)
		if err == nil {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(tree)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if !*tokens {
		tree, err := parser.Parse(source)
		if err == nil {
			code, compileErr := compiler.Compile(tree)
			err = compileErr
			if err == nil {
				err = vm.Run(code, os.Stdout)
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	for _, token := range lexer.New(source).All() {
		fmt.Printf("%d:%d\t%-12s %s\n", token.Line, token.Column, token.Type, strings.TrimSpace(token.Literal))
		if token.Type == lexer.Illegal {
			os.Exit(1)
		}
	}
}
