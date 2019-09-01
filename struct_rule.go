package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

type StructRule struct {
	comment string
	name    string
	field   string

	// inspecting selectors, strings
	cannotMatch []*regexp.Regexp
}

func (rule *StructRule) String() string {
	return fmt.Sprintf("struct rule: %s", rule.comment)
}

func (rule *StructRule) LintMessage(fs *token.FileSet, node ast.Node) string {
	position := fs.Position(node.Pos())

	return fmt.Sprintf(
		"%s:%d:%d:%s",
		position.Filename,
		position.Line,
		position.Column,
		rule.comment,
	)
}

// Action is required to establish a Rule
func (rule *StructRule) Action(fs *token.FileSet, node ast.Node) {
	assignment, ok := node.(*ast.AssignStmt)

	if ok {
		// TODO: inspect Lhs ( such as finding _ assignments which are ignored)
		var isRightType = func(e ast.Expr) bool {
			se, ok := e.(*ast.SelectorExpr)
			if ok {
				strct := fmt.Sprintf("%v.%v", se.X, se.Sel.Name)
				return strct == rule.name
			}
			return false
		}

		// handle stringified struct value assingment, the final check in the
		// ast parsing strategy for this rule type
		var handleStrValue = func(v string, snode ast.Node) {
			if matchAny(v, rule.cannotMatch) {
				fmt.Println(rule.LintMessage(fs, snode))
			}
		}

		// loop over potential multile var assignments
		// and become one with the AST
		for _, expr := range assignment.Rhs {
			switch e := expr.(type) {
			case *ast.CompositeLit:
				if isRightType(e.Type) {
					for _, el := range e.Elts {
						kve, ok := el.(*ast.KeyValueExpr)
						if ok {
							kstr := fmt.Sprintf("%v", kve.Key)
							if kstr == rule.field {
								switch v := kve.Value.(type) {
								case *ast.CompositeLit:
									for _, fel := range v.Elts {
										switch finalT := fel.(type) {
										case *ast.SelectorExpr:
											selector := fmt.Sprintf("%v.%v", finalT.X, finalT.Sel.Name)
											handleStrValue(selector, fel)
										default:
											// TODO: handle more types
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
}
