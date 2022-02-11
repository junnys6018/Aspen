package main

import (
	"errors"
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
