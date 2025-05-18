package code

// Common: basic errors.
// Code must start with 1xxxxx.

const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknow - 500: Internal server error.
	ErrUnknow

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrTokenInvalid - 401: Token invalid
	ErrTokenInvalid

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound
)

// common: database errors
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101

	// ErrRecordNotFound - 404: Record not found.
	ErrRecordNotFound

	// ErrRedis - 500: Redis error.
	ErrRedis

	// ErrCacheNotFound - 404: Cache not found.
	ErrCacheNotFound
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was missed or empty.
	ErrMissingHeader

	// ErrPasswordIncorrect - 401: Password was incorrect.
	ErrPasswordIncorrect

	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied
)

// common encode/decode errors.
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYAML - 500: Data is not valid YAML.
	ErrInvalidYAML

	// ErrEncodingYAML - 500: YAML data could not be encoded.
	ErrEncodingYAML

	// ErrDecodingYAML - 500: YAML data could not be decoded.
	ErrDecodingYAML
)
