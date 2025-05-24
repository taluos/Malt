package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Defined errors.
var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

// Secret contains the basic information of the secret key.
type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

// CacheStrategy defines jwt bearer authentication strategy which called `cache strategy`.
// Secrets are obtained through grpc api interface and cached in memory.
type CacheStrategy struct {
	get func(kid string) (Secret, error)
}

var _ AuthStrategy = &CacheStrategy{}

// NewCacheStrategy create cache strategy with function which can list and cache secrets.
func NewCacheStrategy(get func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{get}
}

// AuthFunc defines cache strategy as the gin authentication middleware.
func (cache CacheStrategy) AuthFunc() fiber.Handler {
	return func(c fiber.Ctx) error {
		var err error
		var rawJWT string

		rawJWT, err = ParseTokenFromHTTPContext(c)
		if err != nil {
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Token is not validable."), nil)
			c.Drop()
			return errors.WithCode(code.ErrSignatureInvalid, "Token is not validable.")
		}

		// Use own validation logic, see below
		var secret Secret

		claims := &jwt.MapClaims{}
		// Verify the token
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is HMAC signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}

			secret, err = cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			return []byte(secret.Key), nil
		}, jwt.WithAudience(AuthzAudience))
		if err != nil || !parsedT.Valid {
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Token is not validable."), nil)
			c.Drop()

			return errors.WithCode(code.ErrSignatureInvalid, "Token is not validable.")
		}

		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			internal.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Drop()

			return errors.WithCode(code.ErrExpired, "expired at: %s", tm)
		}

		c.Set(UsernameKey, secret.Username)
		c.Next()

		return nil
	}
}

// KeyExpired checks if a key has expired, if the value of user.SessionState.Expires is 0, it will be ignored.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}

// ParseTokenFromContext 从 gin.Context 中解析 JWT Token
func ParseTokenFromHTTPContext(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Authorization header not found.
		return "", errors.New("Authorization header not found.")
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		// Incorrect Authorization header format.
		return "", errors.New("Incorrect Authorization header format.")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		// JWT not found.
		return "", errors.New("JWT not found.")
	}

	return token, nil
}
