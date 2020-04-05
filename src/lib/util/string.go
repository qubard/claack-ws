package util

import "strings"

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if (r == ' ') || (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && !(r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func ValidFormString(s string) bool {
	return len(s) > 0 && IsAlphanumeric(s)
}
