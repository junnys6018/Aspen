package main

import (
	"fmt"
	"strings"
)

func GetLocation(source []rune, offset int) (line int, col int) {
	line = 1
	col = 1

	for _, r := range source[:offset] {
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}

	return
}

func GetLine(source []rune, offset int) string {
	start := offset - 1
	end := offset + 1

	for start >= 0 && source[start] != '\n' {
		start--
	}

	for end < len(source) && source[end] != '\n' {
		end++
	}

	return string(source[start+1 : end])
}

func ErrorString(source []rune, message string, location int) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "error: %s\n\n", message)

	lineNumber, col := GetLocation(source, location)

	lineNumberString := fmt.Sprintf("%d", lineNumber)

	line := GetLine(source, location)

	fmt.Fprintf(&builder, "    %s | %s\n", lineNumberString, line)

	for i := 0; i < col+len(lineNumberString)+6; i++ {
		builder.WriteRune(' ')
	}
	builder.WriteString("^-- here.\n")

	return builder.String()
}
