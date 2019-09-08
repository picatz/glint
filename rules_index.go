package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// RulesIndex is a collection of Rule objects
type RulesIndex struct {
	Rules []Rule
}

// NewRulesIndex creates a RulesIndex from a given file.
func NewRulesIndex(filePath string) (*RulesIndex, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	configBuffer := &bytes.Buffer{}

	// strip out line-level (not inline) comments from the json config
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue // skip
		}
		configBuffer.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading rules file:", err)
		os.Exit(1)
	}

	rulesIndex := &RulesIndex{}

	err = json.NewDecoder(configBuffer).Decode(&rulesIndex)

	if err != nil {
		return nil, err
	}

	return rulesIndex, nil
}

// UnmarshalJSON is used when decoding JSON
func (u *RulesIndex) UnmarshalJSON(b []byte) error {
	holder := map[string][]interface{}{}

	err := json.Unmarshal(b, &holder)

	// abort if we have an error other than the wrong type
	if _, ok := err.(*json.UnmarshalTypeError); err != nil && !ok {
		return err
	}

	var extractImportRule = func(ruleSpec map[string]interface{}) {
		v := &ImportRule{}

		if ruleSpec["comment"] != nil {
			v.Comment = ruleSpec["comment"].(string)
		}

		if ruleSpec["match"] != nil {
			for _, cs := range ruleSpec["match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.Match = append(v.Match, rg)
			}
		}

		u.Rules = append(u.Rules, v)
	}

	var extractMethodRule = func(ruleSpec map[string]interface{}) {
		v := &MethodRule{}

		if ruleSpec["call"] != nil {
			call, _ := ruleSpec["call"].(string)
			v.call = call
		}

		if ruleSpec["comment"] != nil {
			comment, _ := ruleSpec["comment"].(string)
			v.comment = comment
		}

		if ruleSpec["avoid"] != nil {
			avoid, _ := ruleSpec["avoid"].(bool)
			v.avoid = avoid
		}

		if ruleSpec["argument"] != nil {
			argumentStr := fmt.Sprint(ruleSpec["argument"])
			argumentStrInt, err := strconv.Atoi(argumentStr)
			if err != nil {
				panic("got unexpected argument type")
			}
			v.argument = argumentStrInt
		}

		if ruleSpec["greater_than"] != nil {
			greaterThanStr := fmt.Sprint(ruleSpec["greater_than"])
			greaterThanInt, err := strconv.Atoi(greaterThanStr)
			if err != nil {
				panic("got unexpected greather than type")
			}
			v.greaterThan = greaterThanInt
		} else {
			v.ignoreGreaterThan = true
		}

		if ruleSpec["less_than"] != nil {
			lessThanStr := fmt.Sprint(ruleSpec["less_than"])
			lessThanInt, err := strconv.Atoi(lessThanStr)
			if err != nil {
				panic("got unexpected greather than type")
			}
			v.lessThan = lessThanInt
		} else {
			v.ignoreLessThan = true
		}

		if ruleSpec["equals"] != nil {
			equalsStr := fmt.Sprint(ruleSpec["equals"])
			equalsInt, err := strconv.Atoi(equalsStr)
			if err != nil {
				panic("got unexpected greather than type")
			}
			v.equals = equalsInt
		} else {
			v.ignoreEquals = true
		}

		if ruleSpec["match"] != nil {
			for _, cs := range ruleSpec["match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.match = append(v.match, rg)
			}
		}

		if ruleSpec["call_match"] != nil {
			for _, cs := range ruleSpec["call_match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected call match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.callMatch = append(v.callMatch, rg)
			}
		}

		u.Rules = append(u.Rules, v)
	}

	var extractStructRule = func(ruleSpec map[string]interface{}) {
		v := &StructRule{}

		if ruleSpec["comment"] != nil {
			comment, _ := ruleSpec["comment"].(string)
			v.comment = comment
		}

		if ruleSpec["name"] != nil {
			name, _ := ruleSpec["name"].(string)
			v.name = name
		}

		if ruleSpec["field"] != nil {
			field, _ := ruleSpec["field"].(string)
			v.field = field
		}

		if ruleSpec["match"] != nil {
			for _, cs := range ruleSpec["match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.match = append(v.match, rg)
			}
		}

		u.Rules = append(u.Rules, v)
	}

	var extractCommentRule = func(ruleSpec map[string]interface{}) {
		v := &CommentRule{}

		if ruleSpec["comment"] != nil {
			comment, _ := ruleSpec["comment"].(string)
			v.comment = comment
		}

		if ruleSpec["match"] != nil {
			for _, cs := range ruleSpec["match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected cannot match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.Match = append(v.Match, rg)
			}
		}

		u.Rules = append(u.Rules, v)
	}

	var extractAssignmentRule = func(ruleSpec map[string]interface{}) {
		v := &AssignmentRule{}

		if ruleSpec["comment"] != nil {
			comment, _ := ruleSpec["comment"].(string)
			v.comment = comment
		}

		if ruleSpec["match"] != nil {
			for _, cs := range ruleSpec["match"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected match type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.match = append(v.match, rg)
			}
		}

		if ruleSpec["is"] != nil {
			for _, cs := range ruleSpec["is"].([]interface{}) {
				str, ok := cs.(string)
				if !ok {
					panic("got unexpected is type")
				}
				rg, err := regexp.Compile(str)
				if err != nil {
					panic(err)
				}
				v.is = append(v.is, rg)
			}
		}

		u.Rules = append(u.Rules, v)
	}

	for _, rule := range holder["rules"] {
		switch r := rule.(type) {
		case map[string]interface{}:
			switch r["type"].(string) {
			case "import":
				extractImportRule(r)
			case "method":
				extractMethodRule(r)
			case "struct":
				extractStructRule(r)
			case "comment":
				extractCommentRule(r)
			case "assignment":
				extractAssignmentRule(r)
			}
		default:
			fmt.Println(r)
		}
	}

	return nil
}
