package main

import (
	"fmt"
	"strings"
)

func GetLine(source []rune, line int) string {
	currLine := 0
	start := -1
	end := -1

	for i, r := range source {
		if r == '\n' {
			currLine++
			start = end
			end = i
			if currLine == line {
				break
			}
		}
	}

	if currLine < line {
		start = end
		end = len(source)
	}

	return string(source[start+1 : end])
}

func ErrorString(source []rune, message string, line int, col int) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "error: %s\n\n", message)

	lineNumberString := fmt.Sprintf("%d", line)

	sourceCodeLine := GetLine(source, line)

	fmt.Fprintf(&builder, "    %s | %s\n", lineNumberString, sourceCodeLine)

	for i := 0; i < col+len(lineNumberString)+6; i++ {
		builder.WriteRune(' ')
	}
	builder.WriteString("^-- here.\n")

	return builder.String()
}
