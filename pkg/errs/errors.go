package errs

import "errors"

var (
	ErrInvalidLoginFormat = errors.New("invalid login format, expected 3-32 unicode literals, numbers, '_' and '-'")

	ErrInvalidPasswordLength = errors.New("invalid login length")
	ErrWrongPassword         = errors.New("wrong password")

	ErrWrongPasswordOrLogin = errors.New("wrong login or password")

	ErrUserAlreadyExsist = errors.New("login already exists")
	ErrUserAlreadyAuth   = errors.New("user already auth")
	ErrAlreadyLoggedIn   = errors.New("another user is currently logged in")
	ErrUserNotFound      = errors.New("user not found")

	ErrUnauthorized = errors.New("unauthorized")
)
