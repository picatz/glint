package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	rulesIndex := &RulesIndex{}

	err = json.NewDecoder(f).Decode(&rulesIndex)

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
				if r["cannot_match"] != nil {
					for _, cs := range r["cannot_match"].([]interface{}) {
						str, ok := cs.(string)
						if !ok {
							panic("got unexpected cannot match type")
						}
						rg, err := regexp.Compile(str)
						if err != nil {
							panic(err)
						}
						v.CannotMatch = append(v.CannotMatch, rg)
					}
				}
				if r["must_match"] != nil {
					for _, ms := range r["must_match"].([]interface{}) {
						str, ok := ms.(string)
						if !ok {
							panic("got unexpected must match type")
						}
						rg, err := regexp.Compile(str)
						if err != nil {
							panic(err)
						}
						v.MustMatch = append(v.MustMatch, rg)
					}
				}
				if r["comment"] != nil {
					v.Comment = r["comment"].(string)
				}
				u.Rules = append(u.Rules, v)
			case "method":
				v := &MethodRule{}

				//requiredRuleParams(map[string]string{
				//	"match":   "string",
				//	"comment": "string",
				//	//"argument": "float64", // JSON number
				//}, r)

				if r["match"] != nil {
					match, _ := r["match"].(string)
					v.match = match
				}

				if r["comment"] != nil {
					comment, _ := r["comment"].(string)
					v.comment = comment
				}

				if r["dont_use"] != nil {
					dontUse, _ := r["dont_use"].(bool)
					v.dontUse = dontUse
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
				}

				if r["less_than"] != nil {
					lessThanStr := fmt.Sprint(r["less_than"])
					lessThanInt, err := strconv.Atoi(lessThanStr)
					if err != nil {
						panic("got unexpected greather than type")
					}
					v.lessThan = lessThanInt
				}

				if r["equals"] != nil {
					equalsStr := fmt.Sprint(r["equals"])
					equalsInt, err := strconv.Atoi(equalsStr)
					if err != nil {
						panic("got unexpected greather than type")
					}
					v.equals = equalsInt
				}

				if r["cannot_match"] != nil {
					for _, cs := range r["cannot_match"].([]interface{}) {
						str, ok := cs.(string)
						if !ok {
							panic("got unexpected cannot match type")
						}
						rg, err := regexp.Compile(str)
						if err != nil {
							panic(err)
						}
						v.cannotMatch = append(v.cannotMatch, rg)
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

func requiredRuleParams(paramWithType map[string]string, givenParams map[string]interface{}) {
	for p, t := range paramWithType {
		if givenParams[p] == nil {
			panic(fmt.Sprintf("required rule param not set %v", p))
		}
		typeStr := fmt.Sprintf("%T", givenParams[p])
		if typeStr != t {
			panic(fmt.Sprintf("rule param not proper type %v", p))
		}
	}
}
