// this file is modified from https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/auth/auth.go
package auth

import (
	rpcmetadata "Malt/api/rpcmetadata"
	"Malt/pkg/errors"
	"Malt/pkg/errors/code"

	"strings"
	"time"

	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secretFunc jwt.Keyfunc

	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
}

func NewAuthenticator(keyfunc jwt.Keyfunc, opts ...AuthOptions) (*Authenticator, error) {
	auth := &Authenticator{
		signingMethod: jwt.SigningMethodES256,
		claims:        func() jwt.Claims { return &CustomClaims{} },
	}

	auth.secretFunc = keyfunc

	// 应用选项
	for _, opt := range opts {
		opt(auth)
	}

	return auth, nil
}

type AuthenMethod interface {
	Authenticate(ctx context.Context, FullMethod string) error
	ValidateToken(ctx context.Context, tokenString string, FullMethod string) error
}

// Authenticate
func (auth *Authenticator) Authenticate(ctx context.Context, FullMethod string) error {

	// 先从metadata提取token字符串
	tokenString, err := extractToken(ctx)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "missing or invalid authorization token")
	}
	// 调用业务层token验证
	return auth.ValidateToken(ctx, tokenString, FullMethod)
}

// Server should put token in redis and metadata both
func (auth *Authenticator) ValidateToken(ctx context.Context, tokenString string, FullMethod string) error {

	// 校验token
	token, err := jwt.ParseWithClaims(tokenString, auth.claims(), auth.secretFunc, jwt.WithValidMethods([]string{auth.signingMethod.Alg()}))
	if err != nil || !token.Valid {
		return errors.WithCode(code.UserNoAuthority, "Invalid JWT token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return errors.WithCode(code.UserNoAuthority, "invalid claims")
	}

	// FullMethod 校验
	if claims.MethodName != FullMethod {
		return errors.WithCode(code.UserNoAuthority, "token not permitted for this method")
	}

	// 检查token是否过期
	now := time.Now()

	// 优先判断 exp
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time) {
		return errors.WithCode(code.UserNoAuthority, "the token expires (via exp)")
	}

	// 或者使用 iat + tokenExpiretime
	if claims.IssuedAt != nil {
		expireAt := claims.IssuedAt.Time.Add(tokenExpiretime)
		if now.After(expireAt) {
			return errors.WithCode(code.UserNoAuthority, "the token expires (via iat + tokenExpiretime)")
		}
	}

	return nil
}

// extractToken reads Authorization Bearer or "token" metadata from context.
// It handles metadata.Get returning a string or a slice of strings.
func extractToken(ctx context.Context) (string, error) {
	md, ok := rpcmetadata.FromServerContext(ctx)
	if !ok {
		return "", errors.WithCode(code.ErrInvalidAuthHeader, "missing metadata")
	}

	// 1. 优先检查 Authorization 头
	if tokens, ok := md["authorization"]; ok && len(tokens) > 0 {
		for _, val := range tokens {
			parts := strings.SplitN(val, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				return parts[1], nil
			}
		}
	}

	// 2. 检查metadata中的token字段
	if tokenVals, ok := md[strings.ToLower(JWTTokenKey)]; ok && len(tokenVals) > 0 {
		return tokenVals[0], nil
	}

	return "", errors.WithCode(code.ErrInvalidAuthHeader, "missing app or token in metadata")
}
