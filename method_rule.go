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
	comment   string
	call      string
	callMatch []*regexp.Regexp
	argument  int

	// generic
	dontUse     bool
	cannotMatch []*regexp.Regexp

	// int action specific
	greaterThan       int
	lessThan          int
	equals            int
	ignoreGreaterThan bool
	ignoreLessThan    bool
	ignoreEquals      bool
}

func (rule *MethodRule) String() string {
	return fmt.Sprintf("method rule for %v (%v): %v", rule.call, rule.argument, rule.greaterThan)
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

func (rule *MethodRule) ProcessMethodCall(methodCall string, fs *token.FileSet, node ast.Node, ce *ast.CallExpr) {

	// digs into the ast node from found function arguments
	var handleArgument func(ast.Node)

	handleArgument = func(arg ast.Node) {
		switch a := arg.(type) {
		case *ast.BasicLit:
			bl := a
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
					// this is trickier than I thought it would be...

					if argInt == rule.equals && !rule.ignoreEquals {
						fmt.Println(rule.LintMessage(fs, bl))
						return
					}

					if argInt < rule.lessThan && !rule.ignoreLessThan {
						fmt.Println(rule.LintMessage(fs, bl))
						return
					}

					if argInt > rule.greaterThan && !rule.ignoreGreaterThan {
						fmt.Println(rule.LintMessage(fs, bl))
						return
					}
				}
			default:
				fmt.Println(a)
			}
		case *ast.Ident: // initial hint there's a variable being used
			switch v := a.Obj.Decl.(type) {
			case *ast.ValueSpec:
				handleArgument(v)
			default:
				fmt.Println(a)
			}
		case *ast.ValueSpec: // handle variables
			for _, v := range a.Values {
				handleArgument(v)
			}
		case *ast.BinaryExpr:
			handleArgument(a.X)
			handleArgument(a.Y)
		default:
			fmt.Println(a)
		}
	}

	if methodCall == rule.call || matchAny(methodCall, rule.callMatch) {
		if rule.dontUse {
			fmt.Println(rule.LintMessage(fs, node))
			return
		}
		if rule.argument == -1 { // all arguments
			for _, arg := range ce.Args {
				handleArgument(arg)
			}
		} else { // a specific argument
			arg := ce.Args[rule.argument]
			handleArgument(arg)
		}
	} else if matchAny(methodCall, rule.cannotMatch) { // TODO: consider removal because call_match
		fmt.Println(rule.LintMessage(fs, node))
	}
}

// Action is required to establish a Rule
func (rule *MethodRule) Action(fs *token.FileSet, node ast.Node) {
	ce, ok := node.(*ast.CallExpr)
	if ok {
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if ok {
			methodCall := fmt.Sprintf("%v.%v", se.X, se.Sel.Name)
			rule.ProcessMethodCall(methodCall, fs, node, ce)
		}
	}
}
