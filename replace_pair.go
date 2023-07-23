package main

import "strings"

type ReplacePair struct {
	before string
	after  string
}

func replace(str string, pairs ...ReplacePair) string {
	var result = str
	for _, pair := range pairs {
		result = strings.ReplaceAll(result, pair.before, pair.after)
	}
	return result
}
