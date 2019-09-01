package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

// MethodRule handles selector expressions (method calls)
type MethodRule struct {
	comment  string
	match    string
	argument int

	// generic
	dontUse     bool
	cannotMatch []*regexp.Regexp

	// int action specific
	greaterThan int
	lessThan    int
	equals      int
}

func (rule *MethodRule) String() string {
	return fmt.Sprintf("method rule for %v (%v): %v", rule.match, rule.argument, rule.greaterThan)
}

// Action is required to establish a Rule
func (rule *MethodRule) Action(fs *token.FileSet, node ast.Node) {
	ce, ok := node.(*ast.CallExpr)
	if ok {
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if ok {
			methodCall := fmt.Sprintf("%v.%v", se.X, se.Sel.Name)
			if methodCall == rule.match {
				if rule.dontUse {
					position := fs.Position(node.Pos())
					mesg := fmt.Sprintf(
						"%s:%d:%d:%s",
						position.Filename,
						position.Line,
						position.Column,
						rule.comment,
					)
					fmt.Println(mesg)
					return
				}
				if rule.argument == -1 { // all arguments
					// TODO: do this
				} else { // a specific argument
					arg := ce.Args[rule.argument]
					bl, ok := arg.(*ast.BasicLit)
					if ok {
						switch bl.Kind {
						case token.STRING:
							strValue := strings.Replace(bl.Value, "\"", "", -1)
							position := fs.Position(bl.Pos())
							mesg := fmt.Sprintf(
								"%s:%d:%d:%s",
								position.Filename,
								position.Line,
								position.Column,
								rule.comment,
							)
							for _, cm := range rule.cannotMatch {
								match := cm.FindString(strValue)
								if match != "" {
									fmt.Println(mesg)
								}
							}
						case token.INT:
							argInt, err := strconv.Atoi(bl.Value)
							position := fs.Position(bl.Pos())
							mesg := fmt.Sprintf(
								"%s:%d:%d:%s",
								position.Filename,
								position.Line,
								position.Column,
								rule.comment,
							)
							if err == nil {
								if argInt <= rule.greaterThan || !(argInt >= rule.lessThan) || (argInt != rule.equals) {
									fmt.Println(mesg)
								}
							}
						}
						//if bl.Kind == token.INT {
						//	argInt, err := strconv.Atoi(bl.Value)
						//	position := fs.Position(bl.Pos())
						//	mesg := fmt.Sprintf(
						//		"%s:%d:%d:%s",
						//		position.Filename,
						//		position.Line,
						//		position.Column,
						//		rule.comment,
						//	)
						//	if err == nil {
						//		if argInt <= rule.greaterThan || !(argInt >= rule.lessThan) || (argInt != rule.equals) {
						//			fmt.Println(mesg)
						//		}
						//	}
						//}
					}
				}
			}
		}
	}
}
