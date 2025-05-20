// this file is modified from https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/auth/auth.go
package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	// secretFunc 用于解析 token 的密钥, 如果是非对称加密则返回公钥，否则返回私钥
	secretFunc    jwt.Keyfunc
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
}

func NewAuthenticator(keyfunc jwt.Keyfunc, opts ...AuthOptions) (*Authenticator, error) {
	auth := &Authenticator{
		signingMethod: jwt.SigningMethodES256,
		claims:        func() jwt.Claims { return &JWT.CustomClaims{} },
	}

	auth.secretFunc = keyfunc

	// 应用选项
	for _, opt := range opts {
		opt(auth)
	}

	return auth, nil
}

type AuthenMethod interface {
	RPCAuthenticate(ctx context.Context, FullMethod string) error
	HTTPAuthenticate(ctx *gin.Context, FullMethod string) error
}

// Authenticate 从 context 中提取 token 字符串，然后进行验证。需要传入 FullMethod，例如 /pkg.Service/Method
// 来验证这个 token 是否对应这个方法。
func (auth *Authenticator) RPCAuthenticate(ctx context.Context, FullMethod string) error {

	// 先从metadata提取token字符串
	//tokenString, err := extractToken(ctx)
	tokenString, err := JWT.ParseTokenFromRPCContext(ctx)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "missing or invalid authorization token")
	}
	// 调用业务层token验证
	return auth.validateToken(tokenString, FullMethod)
}

func (auth *Authenticator) HTTPAuthenticate(ctx *gin.Context, FullMethod string) error {

	tokenString, err := JWT.ParseTokenFromHTTPContext(ctx)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "missing or invalid authorization token")
	}
	return auth.validateToken(tokenString, FullMethod)
}

func (auth *Authenticator) validateToken(tokenString string, FullMethod string) error {

	// 校验token
	token, err := jwt.ParseWithClaims(tokenString, auth.claims(),
		auth.secretFunc,
		jwt.WithValidMethods([]string{auth.signingMethod.Alg()}))
	if err != nil || !token.Valid {
		return errors.WithCode(code.UserNoAuthority, "Invalid JWT token")
	}

	claims, ok := token.Claims.(*JWT.CustomClaims)
	if !ok {
		return errors.WithCode(code.UserNoAuthority, "invalid claims")
	}

	// FullMethod 校验
	if claims.FullMethod != FullMethod {
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
