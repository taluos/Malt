package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/taluos/Malt/pkg/cache"
)

type AuthOption func(*Authenticator)

// WithSigningMethod 设置签名方法
func WithSigningMethod(method jwt.SigningMethod) AuthOption {
	return func(auth *Authenticator) {
		auth.signingMethod = method
	}
}

// WithClaims 设置自定义Claims
func WithClaims(f func() jwt.Claims) AuthOption {
	return func(auth *Authenticator) {
		auth.claims = f
	}
}

// WithSecretFunc 设置密钥函数
func WithSecretFunc(f jwt.Keyfunc) AuthOption {
	return func(auth *Authenticator) {
		auth.keyFunc = f
	}
}

// WithUseCache 启用缓存
func WithUseCache(useCache bool) AuthOption {
	return func(auth *Authenticator) {
		auth.useCache = useCache
	}
}

// WithCache 设置自定义缓存实例
func WithCache(cache cache.CacheMethod) AuthOption {
	return func(auth *Authenticator) {
		auth.cache = cache
		auth.useCache = true
	}
}
