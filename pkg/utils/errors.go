package utils

import "fmt"

type PasswordError struct {
	Code    ErrorCode
	Message string
}

type ImageError struct {
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

const (
	ErrInvalidImageURL ErrorCode = iota
	ErrImageTooLarge
	ErrInvalidImageMimeType
)

func (e PasswordError) Error() string {
	return e.Message
}
func (e ImageError) Error() string {
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

func newImageTooLarge() ImageError {
	return ImageError{
		Code:    ErrImageTooLarge,
		Message: "image is too large (max 50MB)",
	}
}

func newInvalidImageURL() ImageError {
	return ImageError{
		Code:    ErrInvalidImageURL,
		Message: "URL does not point to an image",
	}
}

func newInvalidImageMimeType() ImageError {
	return ImageError{
		Code:    ErrInvalidImageMimeType,
		Message: "invalid image mime type",
	}
}
