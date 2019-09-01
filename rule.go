package main

import (
	"go/ast"
	"go/token"
)

// Rule is the basic unit of glint
type Rule interface {
	Action(fs *token.FileSet, node ast.Node)
}
