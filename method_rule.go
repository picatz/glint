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
	matchAny []*regexp.Regexp
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

func (rule *MethodRule) LintMessage(fs *token.FileSet, node ast.Node) string {
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
func (rule *MethodRule) Action(fs *token.FileSet, node ast.Node) {
	ce, ok := node.(*ast.CallExpr)
	if ok {
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if ok {
			methodCall := fmt.Sprintf("%v.%v", se.X, se.Sel.Name)

			if methodCall == rule.match {
				if rule.dontUse {
					fmt.Println(rule.LintMessage(fs, node))
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
							for _, cm := range rule.cannotMatch {
								match := cm.FindString(strValue)
								if match != "" {
									fmt.Println(rule.LintMessage(fs, bl))
								}
							}
						case token.INT:
							argInt, err := strconv.Atoi(bl.Value)
							if err == nil {
								if argInt <= rule.greaterThan || !(argInt >= rule.lessThan) || (argInt != rule.equals) {
									fmt.Println(rule.LintMessage(fs, bl))
								}
							}
						}
					}
				}
			} else if matchAny(methodCall, rule.matchAny) {
				// TODO deal with match any cases
			}
		}
	}
}
