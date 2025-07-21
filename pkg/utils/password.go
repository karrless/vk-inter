package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type PasswordError struct {
	Code    ErrorCode
	Message string
}

type ErrorCode int

const (
	ErrTooShort ErrorCode = iota
	ErrNoMixedCase
	ErrNoNumber
	ErrNoSpecialChar
	ErrContainsLogin
)

func (e PasswordError) Error() string {
	return e.Message
}

func newTooShortError(gotLength int) PasswordError {
	return PasswordError{
		Code:    ErrTooShort,
		Message: fmt.Sprintf("password must contain at least 8 characters (got %d)", gotLength),
	}
}

func newNoMixedCaseError() PasswordError {
	return PasswordError{
		Code:    ErrNoMixedCase,
		Message: "password must contain both uppercase and lowercase letters",
	}
}
func newNoNumber() PasswordError {
	return PasswordError{
		Code:    ErrNoNumber,
		Message: "password must contain at least one number",
	}
}

func newNoSpecialChar() PasswordError {
	return PasswordError{
		Code:    ErrNoSpecialChar,
		Message: "password must contain at least one special character (!@#$%^&* etc.)",
	}
}

func newContainsLogin() PasswordError {
	return PasswordError{
		Code:    ErrContainsLogin,
		Message: "password cannot contain your username",
	}
}

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
