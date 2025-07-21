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

	ErrListingInvalidTitle       = errors.New("invalid title format, expected 3-100 chars")
	ErrListingInvalidDescription = errors.New("invalid description format, expected be 10-5000 chars")
	ErrListingInvalidImageURL    = errors.New("invalid image URL, expected valid image valid URL")
	ErrListingInvalidPrice       = errors.New("invalid price, expected decimal with up to 2 decimal places between 0 and 1_000_000_000")

	ErrPriceSorting = errors.New("Max price must be greater then min pirce")
)
