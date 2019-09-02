package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
)

var (
	rulesFile  string
	targetCode []string
)

func init() {
	flag.StringVar(&rulesFile, "rules", "Glintfile", "the JSON configuation files for glint")
	flag.Parse()
	targetCode = flag.Args()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	rulesIndex, err := NewRulesIndex(rulesFile)

	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	fs := token.NewFileSet()

	conf := loader.Config{
		Fset: fs,
	}

	for _, arg := range targetCode {
		if fileExists(arg) {
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
		} else {
			conf.Import(arg)
		}
	}

	lprog, err := conf.Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, programInfo := range lprog.InitialPackages() {
		for _, file := range programInfo.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				for _, rule := range rulesIndex.Rules {
					rule.Action(fs, n)
				}
				return true
			})
		}
	}
}
