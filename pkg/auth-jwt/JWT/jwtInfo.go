// JWTInfo is a struct that contains the information for JWT authentication.
//
// if you want to use JWT authentication, you can set this field in the ServerOptions.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInfo struct {
	Realm      string        // default is "JWT"
	token      string        // jwt token
	keyFunc    jwt.Keyfunc   // default is ParseRSAPublicKeyFromPEM
	Timeout    time.Duration // default is 5 minutes
	MaxRefresh time.Duration // default is 10 minutes

	privateKey string
	userID     string
	FullMethod string
	role       string
	expire     time.Duration
}

func NewJwtInfo(privateKey string, userID string, FullMethod string, role string, expire time.Duration, opts ...JWTOptions) (*JwtInfo, error) {
	var err error
	jwtInfo := &JwtInfo{
		Realm: "JWT",
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM([]byte(privateKey))
		},
		Timeout:    time.Minute * 5,
		MaxRefresh: time.Minute * 10,

		privateKey: privateKey,
		userID:     userID,
		FullMethod: FullMethod,
		role:       role,
		expire:     expire,
	}

	for _, opt := range opts {
		opt(jwtInfo)
	}

	jwtInfo.token, err = GenerateJWT(
		jwtInfo.privateKey,
		jwtInfo.userID,
		jwtInfo.FullMethod,
		jwtInfo.role,
		jwtInfo.expire,
	)

	return jwtInfo, err
}

type JWTOptions func(*JwtInfo)

func WithRealm(realm string) JWTOptions {
	return func(o *JwtInfo) {
		o.Realm = realm
	}
}

func WithTimeout(timeout time.Duration) JWTOptions {
	return func(o *JwtInfo) {
		o.Timeout = timeout
	}
}

func WithMaxRefresh(maxRefresh time.Duration) JWTOptions {
	return func(o *JwtInfo) {
		o.MaxRefresh = maxRefresh
	}
}
func WithKeyFunc(keyFunc jwt.Keyfunc) JWTOptions {
	return func(o *JwtInfo) {
		o.keyFunc = keyFunc
	}
}

func (j *JwtInfo) Keyfunc() jwt.Keyfunc {
	return j.keyFunc
}

func (j *JwtInfo) Token() string {
	return j.token
}
