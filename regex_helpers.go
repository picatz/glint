package main

import (
	"regexp"
)

func matchAny(str string, rs []*regexp.Regexp) bool {
	for _, r := range rs {
		match := r.FindString(str)
		if match != "" {
			return true
		}
	}
	return false
}
