package auth

import (
	casbin "github.com/taluos/Malt/core/RBAC/Casbin"
	cosjwt "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secretFunc   jwt.Keyfunc
	RBACEnforcer *casbin.RBACEnforcer

	publicKey     string
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
}

func NewAuthenticator(publiKey string, enforcer *casbin.RBACEnforcer, opts ...AuthOptions) (*Authenticator, error) {
	auth := &Authenticator{
		publicKey:     publiKey,
		signingMethod: jwt.SigningMethodES256,
		claims:        func() jwt.Claims { return &cosjwt.CustomClaims{} },
	}

	// Init keyFunc
	keyFunc := func(token *jwt.Token) (any, error) {
		// 确保算法匹配
		if token.Method != auth.signingMethod {
			return nil, errors.WithCode(code.ErrInvalidAuthHeader, "unexpected signing method: %v, expected: %v", token.Header["alg"], auth.signingMethod.Alg())
		}

		pubKey, err := jwt.ParseECPublicKeyFromPEM([]byte(publiKey))
		if err != nil {
			return nil, errors.WithCode(code.ErrInvalidAuthHeader, "failed to parse public key: %v", err)
		}

		return pubKey, nil
	}

	auth.secretFunc = keyFunc

	auth.RBACEnforcer = enforcer

	// 应用选项
	for _, opt := range opts {
		opt(auth)
	}

	return auth, nil
}

func (auth *Authenticator) Authenticate(token, path, method string) error {
	role, err := cosjwt.ParseRoleFromHTTPContext(token, auth.publicKey)
	if err != nil {
		return errors.WithCode(code.ErrInvalidAuthHeader, "failed to parse role from context: %v", err)
	}
	ok := auth.validateAuth(role, path, method)
	if !ok {
		return errors.WithCode(code.ErrInvalidAuthHeader, "failed to verify auth")
	}

	return nil
}

func (auth *Authenticator) validateAuth(role, path, method string) bool {
	ok, err := auth.RBACEnforcer.VerifyAuth(role, path, method)
	if err != nil {
		return false
	}
	return ok
}
