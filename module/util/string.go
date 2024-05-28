package util

import (
	"strings"
	"unicode"
)

func GetSecureString(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
