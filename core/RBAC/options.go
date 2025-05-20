package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthOptions func(*Authenticator)

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) AuthOptions {
	return func(auth *Authenticator) {
		auth.signingMethod = method
	}
}

// WithClaims with customer claim
// If you use it in Server, f needs to return a new jwt.Claims object each time to avoid concurrent write problems
// If you use it in Client, f only needs to return a single object to provide performance
func WithClaims(f func() jwt.Claims) AuthOptions {
	return func(auth *Authenticator) {
		auth.claims = f
	}
}

// WithPublicKey with public key option.
func WithPublicKey(key string) AuthOptions {
	return func(auth *Authenticator) {
		auth.publicKey = key
	}
}
