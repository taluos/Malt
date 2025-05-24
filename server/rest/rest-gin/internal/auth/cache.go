package auth

import (
	"time"

	cosjwt "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	internal "github.com/taluos/Malt/server/rest/rest-gin/internal"

	"github.com/gin-gonic/gin"
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
func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var rawJWT string
		authHeader := c.GetHeader("Authorization")
		rawJWT, err = cosjwt.ParseTokenFromHTTPContext(authHeader)
		if err != nil {
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Token is not validable."), nil)
			c.Abort()
			return
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
			c.Abort()

			return
		}

		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			internal.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Set(UsernameKey, secret.Username)
		c.Next()
	}
}

// KeyExpired checks if a key has expired, if the value of user.SessionState.Expires is 0, it will be ignored.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}
