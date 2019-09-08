package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

// ImportRule handles import specifications.
type ImportRule struct {
	Comment string

	Match []*regexp.Regexp
}

func (rule *ImportRule) String() string {
	return fmt.Sprintf("import rule: match:%v", rule.Match)
}

// Action is required to establish a Rule
func (rule *ImportRule) Action(fs *token.FileSet, node ast.Node) {
	importSpec, ok := node.(*ast.ImportSpec)

	if ok {
		importPath := strings.Replace(importSpec.Path.Value, "\"", "", -1)

		if matchAny(importPath, rule.Match) {
			fmt.Println(rule.LintMessage(fs, node))
		}
	}
}

func (rule *ImportRule) LintMessage(fs *token.FileSet, node ast.Node) string {
	position := fs.Position(node.Pos())

	return fmt.Sprintf(
		"%s:%d:%d:%s",
		position.Filename,
		position.Line,
		position.Column,
		rule.Comment,
	)
}
