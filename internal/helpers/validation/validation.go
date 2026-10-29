package validation

import (
	"strings"
	"unicode/utf8"
)

func ValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}

	for _, r := range username {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_') {
			return false
		}
	}

	return true
}

func ValidEmail(email string) bool {
	return len(email) <= 254 &&
		strings.Contains(email, "@") &&
		!strings.ContainsAny(email, " \t\r\n") &&
		utf8.ValidString(email)
}

func ValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 128
}
