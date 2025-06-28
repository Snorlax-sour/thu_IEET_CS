package main

import (
	"testing"
)

func TestSplit_string(t *testing.T) {
	got := Split_course_type_and_name("選修- course name")
	for _, content := range got {
		t.Log(content) // need run go test . -v
	}
}
