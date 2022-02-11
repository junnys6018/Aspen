package main

import "strings"

type ErrorReporter interface {
	Error() string
	Push(line int, col int, message string)
	HadError() bool
}

type ErrorData struct {
	line    int
	col     int
	message string
}

type AspenError struct {
	source []rune
	data   []ErrorData
}

func NewErrorReporter(source []rune) *AspenError {
	return &AspenError{source: source}
}

func (e *AspenError) Push(line int, col int, message string) {
	e.data = append(e.data, ErrorData{line, col, message})
}

func (e *AspenError) HadError() bool {
	return len(e.data) != 0
}

func (e *AspenError) Error() string {
	builder := strings.Builder{}

	for i, datum := range e.data {
		builder.WriteString(ErrorString(e.source, datum.message, datum.line, datum.col))
		if i != len(e.data)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}
