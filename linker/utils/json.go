package utils

import (
	"strings"
)

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type jsonTagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseJsonTag(tag string) (string, jsonTagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], jsonTagOptions(tag[idx+1:])
	}
	return tag, jsonTagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o jsonTagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
