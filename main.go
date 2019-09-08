package main

import (
	"errors"
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
	flag.ErrHelp = errors.New("usage:\n\tglint [options] [files or packages]\n\n  -rules string\n\tthe JSON configuation files for glint (default \"Glintfile\")\n")

	flag.StringVar(&rulesFile, "rules", "Glintfile", "the JSON configuation file for glint")
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
		fmt.Println("error parsing rules file", rulesFile, ":", err)
		os.Exit(1)
	}

	if len(targetCode) == 0 {
		fmt.Println("no file names or packages given!\n\n", flag.ErrHelp)
		os.Exit(1)
	}

	fs := token.NewFileSet()

	conf := loader.Config{
		Fset: fs,
	}

	for _, arg := range targetCode {
		if fileExists(arg) {
			f, err := parser.ParseFile(fs, arg, nil, parser.AllErrors)
			f, err = parser.ParseFile(fs, arg, nil, parser.ParseComments)
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

	if err == nil {
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
}
