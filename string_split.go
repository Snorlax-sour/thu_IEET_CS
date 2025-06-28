package main

import (
	"strings"
)

func Split_course_type_and_name(name string) []string {
	// this function cut '-' symbol
	cut_symbol := "-"
	full_string := strings.Split(name, cut_symbol)
	return full_string
}
