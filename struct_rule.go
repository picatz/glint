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
	e, ok := node.(*ast.CompositeLit)

	if ok {
		var isRightType = func(e ast.Expr) bool {
			if rule.name == "" {
				// allow any type of struct field to be checked
				// if no struct name (go type) is specified
				return true
			}
			switch t := e.(type) {
			case *ast.SelectorExpr:
				strct := fmt.Sprintf("%v.%v", t.X, t.Sel.Name)
				return strct == rule.name
			case *ast.Ident:
				return t.Name == rule.name
			default:
				// TODO: handle more possibilites?
				return false
			}
		}

		// handle stringified struct value assingment, the final check in the
		// ast parsing strategy for this rule type
		var handleStrValue = func(v string, snode ast.Node) {
			if matchAny(v, rule.cannotMatch) {
				fmt.Println(rule.LintMessage(fs, snode))
			}
		}

		if isRightType(e.Type) {
			if len(e.Elts) == 0 { // no struct fileds given
				handleStrValue("nil", e)
			}
			for _, el := range e.Elts {
				kve, ok := el.(*ast.KeyValueExpr)
				if ok {
					kstr := fmt.Sprintf("%v", kve.Key)
					if kstr == rule.field { // TODO: allow check for any field?
						switch v := kve.Value.(type) {
						case *ast.CompositeLit:
							for _, fel := range v.Elts {
								switch finalT := fel.(type) {
								case *ast.SelectorExpr:
									selector := fmt.Sprintf("%v.%v", finalT.X, finalT.Sel.Name)
									handleStrValue(selector, fel)
								case *ast.Ident:
									strData := fmt.Sprintf("%v", finalT.Obj.Data)
									handleStrValue(strData, fel)
								default:
									// TODO: handle more types?
								}
							}
						case *ast.Ident:
							strData := fmt.Sprintf("%v", v.Obj)
							handleStrValue(strData, kve)
						case *ast.SelectorExpr:
							selector := fmt.Sprintf("%v.%v", v.X, v.Sel.Name)
							handleStrValue(selector, v)
						case *ast.BasicLit:
							strData := fmt.Sprintf("%v", v.Value)
							handleStrValue(strData, kve)
						default:
							// TODO: handle more types?
						}
					}
				}
			}
		}
	}
}
