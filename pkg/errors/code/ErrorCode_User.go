package code

// Common: user errors.
// Code must start with 1xxxxx.

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists

	// ErrUserPasswordIncorrect - 401: Password was incorrect.
	ErrUserPasswordIncorrect

	// ErrUserPasswordTooShort - 400: Password is too short.
	ErrUserPasswordTooShort

	// ErrUserPasswordTooLong - 400: Password is too long.
	ErrUserPasswordTooLong

	// ErrUserPasswordInvalid - 400: Password is invalid.
	ErrUserPasswordInvalid

	// ErrUserPasswordNotMatch - 400: Password not match.
	ErrUserPasswordNotMatch

	//UserNoAuthority - 403: User no authority.
	UserNoAuthority
)
