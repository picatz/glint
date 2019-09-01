package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	rulesIndex, err := NewRulesIndex("rules.json")

	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	fs := token.NewFileSet()

	for _, arg := range os.Args[1:] {
		f, err := parser.ParseFile(fs, arg, nil, parser.AllErrors)
		if err != nil {
			log.Printf("could not parse %s: %v", arg, err)
			continue
		}

		ast.Inspect(f, func(n ast.Node) bool {
			for _, rule := range rulesIndex.Rules {
				rule.Action(fs, n)
			}
			return true
		})
	}
}
