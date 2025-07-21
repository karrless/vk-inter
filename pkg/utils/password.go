package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidatePassword checks password
func ValidatePassword(password, login string) error {
	if len(password) < 8 {
		return newTooShortError(len(password))
	}

	hasUpper, hasLower := false, false
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		}
	}
	if !hasUpper || !hasLower {
		return newNoMixedCaseError()
	}

	if !regexp.MustCompile(`\d`).MatchString(password) {
		return newNoNumber()
	}

	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return newNoSpecialChar()
	}

	if strings.Contains(strings.ToLower(password), strings.ToLower(login)) {
		return newContainsLogin()
	}

	return nil
}
