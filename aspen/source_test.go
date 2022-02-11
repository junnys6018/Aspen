package main

import "testing"

func TestGetLine(t *testing.T) {
	testCase := func(source string, expect string, line int) {
		got := GetLine([]rune(source), line)
		if expect != got {
			t.Errorf("expected \"%s\", got \"%s\"", expect, got)
		}
	}

	testCase("", "", 1)

	testCase("line1\nline2\nline3", "line1", 1)

	testCase("line1\nline2\nline3", "line2", 2)

	testCase("line1\nline2\nline3", "line3", 3)

	testCase("line1\nline2\nline3\n", "", 4)
}
