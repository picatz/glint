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

				if r["call"] != nil {
					call, _ := r["call"].(string)
					v.call = call
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

				if r["call_match"] != nil {
					for _, cs := range r["call_match"].([]interface{}) {
						str, ok := cs.(string)
						if !ok {
							panic("got unexpected cannot match type")
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
