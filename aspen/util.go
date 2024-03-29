package main

import (
	"errors"
	"fmt"
	"strings"
)

/**
 * Takes a quoted string then unquotes, and unescapes it
 */
func UnescapeString(s string) (string, error) {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return "", errors.New("string is not quoted")
	}

	builder := strings.Builder{}

	for i := 1; i < len(s)-1; i++ {
		b := s[i]
		if b == '\\' {
			i++
			if i < len(s)-1 {
				switch s[i] {
				case 'n':
					builder.WriteByte('\n')
				default:
					return "", errors.New("bad escape sequence")
				}
			} else {
				return "", errors.New("bad escape sequence")
			}
		} else {
			builder.WriteByte(b)
		}
	}

	return builder.String(), nil
}

func Unreachable(message string) {
	panic(fmt.Sprintf("%s: unreachable code", message))
}

func AddString(lhs, rhs []rune) []rune {
	new := make([]rune, 0, len(lhs)+len(rhs))
	new = append(new, lhs...)
	new = append(new, rhs...)
	return new
}

func OrdinalSuffixOf(i int) string {
	j := i % 10
	k := i % 100
	if j == 1 && k != 11 {
		return fmt.Sprintf("%dst", i)
	}
	if j == 2 && k != 12 {
		return fmt.Sprintf("%dnd", i)
	}
	if j == 3 && k != 13 {
		return fmt.Sprintf("%drd", i)
	}
	return fmt.Sprintf("%dth", i)
}
