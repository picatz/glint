package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

// AssignmentRule handles import specifications.
type AssignmentRule struct {
	comment string

	// matchers
	is    []*regexp.Regexp
	match []*regexp.Regexp
}

func (rule *AssignmentRule) String() string {
	return fmt.Sprintf("assignment rule: match:%v", rule.match)
}

// Action is required to establish a Rule
func (rule *AssignmentRule) Action(fs *token.FileSet, node ast.Node) {
	assignSpec, ok := node.(*ast.AssignStmt)

	if ok {
		for _, rh := range assignSpec.Rhs {
			switch ce := rh.(type) {
			case *ast.CallExpr:
				switch ident := ce.Fun.(type) {
				case *ast.Ident:
					switch decl := ident.Obj.Decl.(type) {
					case *ast.FuncDecl:
						for i, field := range decl.Type.Results.List {
							switch t := field.Type.(type) {
							case *ast.Ident:
								// if node is of rule "is" type
								if matchAny(t.Name, rule.is) {
									// then check the left hand side of the assignment
									switch lh := assignSpec.Lhs[i].(type) {
									case *ast.Ident:
										// and check the var names
										if matchAny(lh.Name, rule.match) {
											fmt.Println(rule.LintMessage(fs, lh))
										}
									default:
										// TODO: figure out more types?
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (rule *AssignmentRule) LintMessage(fs *token.FileSet, node ast.Node) string {
	position := fs.Position(node.Pos())

	return fmt.Sprintf(
		"%s:%d:%d:%s",
		position.Filename,
		position.Line,
		position.Column,
		rule.comment,
	)
}
