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

	CannotMatch []*regexp.Regexp
	MustMatch   []*regexp.Regexp
}

func (rule *ImportRule) String() string {
	return fmt.Sprintf("import rule: cannot_match:%v must_match:%v", rule.CannotMatch, rule.MustMatch)
}

// Action is required to establish a Rule
func (rule *ImportRule) Action(fs *token.FileSet, node ast.Node) {
	importSpec, ok := node.(*ast.ImportSpec)

	if ok {
		importPath := strings.Replace(importSpec.Path.Value, "\"", "", -1)

		for _, re := range rule.CannotMatch {
			match := re.FindString(importPath)
			if match != "" {
				position := fs.Position(node.Pos())
				mesg := fmt.Sprintf(
					"%s:%d:%d:%s",
					position.Filename,
					position.Line,
					position.Column,
					rule.Comment,
				)
				fmt.Println(mesg)
			}
		}

		for _, re := range rule.MustMatch {
			match := re.FindString(importPath)
			if match == "" {
				position := fs.Position(node.Pos())
				mesg := fmt.Sprintf(
					"%s:%d:%d:%s",
					position.Filename,
					position.Line,
					position.Column,
					rule.Comment,
				)
				fmt.Println(mesg)
			}
		}
	}
}
