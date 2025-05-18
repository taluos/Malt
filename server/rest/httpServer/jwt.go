package httpserver

import (
	"time"
)

var (
	defaultJWTKey = ":36#Xb#un-*!SXz4:V<sUbAV|$%d5-X6"
)

type JwtInfo struct {
	// default is "JWT"
	Realm string
	// default is empty
	Key string
	// default is 5 minutes
	Timeout time.Duration
	// default is 10 minutes
	MaxRefresh time.Duration
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
