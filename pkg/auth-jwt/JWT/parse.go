package jwt

import (
	"context"
	"strings"

	rpcmetadata "github.com/taluos/Malt/api/rpcmetadata"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"

	"github.com/golang-jwt/jwt/v5"
)

// ParseRoleFromHTTPContext 从HTTP头中解析JWT Token的角色
func ParseRoleFromHTTPContext(token string, publicKey string) (string, error) {
	role, err := parseRoleFromToken(publicKey, token)
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed")
	}
	return role, nil
}

// ParseTokenFromHTTPContext 从HTTP Authorization头中解析JWT Token
func ParseTokenFromHTTPContext(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header not found")
	}
	
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("Incorrect Authorization header format")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("JWT not found")
	}

	return token, nil
}

// ParseTokenFromRPCContext 从RPC context中解析JWT Token
func ParseTokenFromRPCContext(ctx context.Context) (string, error) {
	md, ok := rpcmetadata.FromServerContext(ctx)
	if !ok {
		return "", errors.WithCode(code.ErrInvalidAuthHeader, "missing metadata")
	}

	if tokens, ok := md["authorization"]; ok && len(tokens) > 0 {
		for _, val := range tokens {
			parts := strings.SplitN(val, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				return parts[1], nil
			}
		}
	}

	return "", errors.WithCode(code.ErrInvalidAuthHeader, "missing token in metadata")
}

// parseRoleFromToken 从JWT Token中解析角色
func parseRoleFromToken(publicKey string, tokenString string) (string, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		
		pubKey, err := jwt.ParseECPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return nil, errors.Wrapf(err, "parse public key failed")
		}
		return pubKey, nil
	})
	
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed")
	}

	if !token.Valid {
		return "", errors.New("token is invalid")
	}

	return claims.GetRole(), nil
}