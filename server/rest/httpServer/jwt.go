// JWTInfo is a struct that contains the information for JWT authentication.
//
// if you want to use JWT authentication, you can set this field in the ServerOptions.
package httpserver

import (
	"time"
)

type JwtInfo struct {
	Realm      string        // default is "JWT"
	Key        string        // user should provide a key for JWT signing and verifying
	Timeout    time.Duration // default is 5 minutes
	MaxRefresh time.Duration // default is 10 minutes
}

func NewJwtInfo(opts ...JWTOptions) *JwtInfo {
	jwtInfo := &JwtInfo{
		Realm:      "JWT",
		Key:        defaultJWTKey,
		Timeout:    time.Minute * 5,
		MaxRefresh: time.Minute * 10,
	}
	for _, opt := range opts {
		opt(jwtInfo)
	}
	return jwtInfo
}

type JWTOptions func(*JwtInfo)

func WithRealm(realm string) JWTOptions {
	return func(o *JwtInfo) {
		o.Realm = realm
	}
}

func WithKey(key string) JWTOptions {
	return func(o *JwtInfo) {
		o.Key = key
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
