package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

// CommentRule helps lint a program's comments.
type CommentRule struct {
	comment string

	// Matchers
	Match []*regexp.Regexp
}

func (rule *CommentRule) String() string {
	return fmt.Sprintf("comment rule: match:%v", rule.Match)
}

// Action is required to establish a Rule
func (rule *CommentRule) Action(fs *token.FileSet, node ast.Node) {
	f, ok := node.(*ast.File)
	if ok {
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				if matchAny(c.Text, rule.Match) {
					position := fs.Position(c.Pos())
					mesg := fmt.Sprintf(
						"%s:%d:%d:%s",
						position.Filename,
						position.Line,
						position.Column,
						rule.comment,
					)
					fmt.Println(mesg)
				}
			}
		}
	}
}
