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

	for _, rule := range holder["rules"] {
		switch r := rule.(type) {
		case map[string]interface{}:
			switch r["type"].(string) {
			case "import":
				v := &ImportRule{}

				if r["comment"] != nil {
					v.Comment = r["comment"].(string)
				}

				if r["match"] != nil {
					for _, cs := range r["match"].([]interface{}) {
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
			case "method":
				v := &MethodRule{}

				if r["call"] != nil {
					call, _ := r["call"].(string)
					v.call = call
				}

				if r["comment"] != nil {
					comment, _ := r["comment"].(string)
					v.comment = comment
				}

				if r["avoid"] != nil {
					avoid, _ := r["avoid"].(bool)
					v.avoid = avoid
				}

				if r["argument"] != nil {
					argumentStr := fmt.Sprint(r["argument"])
					argumentStrInt, err := strconv.Atoi(argumentStr)
					if err != nil {
						panic("got unexpected argument type")
					}
					v.argument = argumentStrInt
				}

				if r["greater_than"] != nil {
					greaterThanStr := fmt.Sprint(r["greater_than"])
					greaterThanInt, err := strconv.Atoi(greaterThanStr)
					if err != nil {
						panic("got unexpected greather than type")
					}
					v.greaterThan = greaterThanInt
				} else {
					v.ignoreGreaterThan = true
				}

				if r["less_than"] != nil {
					lessThanStr := fmt.Sprint(r["less_than"])
					lessThanInt, err := strconv.Atoi(lessThanStr)
					if err != nil {
						panic("got unexpected greather than type")
					}
					v.lessThan = lessThanInt
				} else {
					v.ignoreLessThan = true
				}

				if r["equals"] != nil {
					equalsStr := fmt.Sprint(r["equals"])
					equalsInt, err := strconv.Atoi(equalsStr)
					if err != nil {
						panic("got unexpected greather than type")
					}
					v.equals = equalsInt
				} else {
					v.ignoreEquals = true
				}

				if r["match"] != nil {
					for _, cs := range r["match"].([]interface{}) {
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

				if r["call_match"] != nil {
					for _, cs := range r["call_match"].([]interface{}) {
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
			case "struct":
				v := &StructRule{}

				if r["comment"] != nil {
					comment, _ := r["comment"].(string)
					v.comment = comment
				}

				if r["name"] != nil {
					name, _ := r["name"].(string)
					v.name = name
				}

				if r["field"] != nil {
					field, _ := r["field"].(string)
					v.field = field
				}

				if r["match"] != nil {
					for _, cs := range r["match"].([]interface{}) {
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
			case "comment":
				v := &CommentRule{}

				if r["comment"] != nil {
					comment, _ := r["comment"].(string)
					v.comment = comment
				}

				if r["match"] != nil {
					for _, cs := range r["match"].([]interface{}) {
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
		default:
			fmt.Println(r)
		}
	}

	return nil
}
