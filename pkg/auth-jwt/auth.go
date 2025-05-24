// this file is modified from https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/auth/auth.go
package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/cache"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	// keyFunc 用于解析token的密钥
	keyFunc       jwt.Keyfunc
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims

	// 缓存相关
	useCache bool
	cache    cache.CacheMethod
}

func NewAuthenticator(keyfunc jwt.Keyfunc, opts ...AuthOption) (*Authenticator, error) {
	auth := &Authenticator{
		signingMethod: jwt.SigningMethodES256,
		claims:        func() jwt.Claims { return &JWT.CustomClaims{} },
		keyFunc:       keyfunc,
		useCache:      false,
	}

	// 应用选项
	for _, opt := range opts {
		opt(auth)
	}

	// 如果启用缓存但没有提供缓存实例，创建默认缓存
	if auth.useCache && auth.cache == nil {
		cacheInstance, err := cache.NewCache(time.Minute*10, 1000) // 默认10分钟TTL，1000个条目
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create default cache")
		}
		auth.cache = cacheInstance
	}

	return auth, nil
}

type AuthMethod interface {
	RPCAuthenticate(ctx context.Context, fullMethod string) error
	HTTPAuthenticate(authHeader, fullMethod string) error
}

// RPCAuthenticate 从context中提取token字符串，然后进行验证
func (auth *Authenticator) RPCAuthenticate(ctx context.Context, userID, fullMethod, role string) error {
	tokenString, err := JWT.ParseTokenFromRPCContext(ctx)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "missing or invalid authorization token")
	}
	return auth.validateToken(tokenString, userID, fullMethod, role)
}

// HTTPAuthenticate HTTP认证
func (auth *Authenticator) HTTPAuthenticate(authHeader, userID, fullMethod, role string) error {
	tokenString, err := JWT.ParseTokenFromHTTPContext(authHeader)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "missing or invalid authorization token")
	}
	return auth.validateToken(tokenString, userID, fullMethod, role)
}

// validateToken 验证token
func (auth *Authenticator) validateToken(tokenString, userID, fullMethod, role string) error {
	// 如果启用缓存，先检查缓存
	if auth.useCache && auth.cache != nil {
		cacheKey := auth.generateCacheKey(tokenString, fullMethod)
		if cached, found := auth.cache.Get(cacheKey); found {
			if result, ok := cached.(bool); ok && result {
				return nil // 缓存中存在且验证通过
			}
			// 如果缓存中存在但验证失败，继续进行完整验证
		}
	}

	// 校验token
	token, err := jwt.ParseWithClaims(tokenString, auth.claims(),
		auth.keyFunc,
		jwt.WithValidMethods([]string{auth.signingMethod.Alg()}))
	if err != nil || !token.Valid {
		// 验证失败，如果使用缓存则缓存失败结果
		if auth.useCache && auth.cache != nil {
			cacheKey := auth.generateCacheKey(tokenString, fullMethod)
			auth.cache.Set(cacheKey, false)
		}
		return errors.WithCode(code.UserNoAuthority, "Invalid JWT token")
	}

	claims, ok := token.Claims.(*JWT.CustomClaims)
	if !ok {
		return errors.WithCode(code.UserNoAuthority, "invalid claims")
	}

	// FullMethod校验
	if fullMethod != "" && claims.GetFullMethod() != fullMethod {
		return errors.WithCode(code.UserNoAuthority, "Not permitted for this method")
	}
	// UserID校验
	if userID != "" && claims.GetUserID() != userID {
		return errors.WithCode(code.UserNoAuthority, "Not permitted for this user")
	}
	// Role校验
	if role != "" && claims.GetRole() != role {
		return errors.WithCode(code.UserNoAuthority, "Not permitted for this role")
	}

	// 检查token是否过期
	now := time.Now()

	// 优先判断exp
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time) {
		return errors.WithCode(code.UserNoAuthority, "the token expires (via exp)")
	}

	// 再判断nbf
	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
		return errors.WithCode(code.UserNoAuthority, "the token is not yet valid (via nbf)")
	}

	// 验证成功，如果使用缓存则缓存成功结果
	if auth.useCache && auth.cache != nil {
		cacheKey := auth.generateCacheKey(tokenString, fullMethod)
		auth.cache.Set(cacheKey, true)
	}

	return nil
}

// generateCacheKey 生成缓存键
func (auth *Authenticator) generateCacheKey(tokenString, fullMethod string) string {
	hash := sha256.Sum256([]byte(tokenString + ":" + fullMethod))
	return fmt.Sprintf("auth:%x", hash)
}
