package auth

import (
	"context"
	"strings"

	rpcmetadata "github.com/taluos/Malt/api/rpcmetadata"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ParseRoleFromContext 从 gin.Context 中解析 JWT Token 的角色
func ParseRoleFromHTTPContext(c *gin.Context, publiKey string) (string, error) {
	token, err := ParseTokenFromHTTPContext(c)
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed.")
	}

	role, err := parseRoleFormToken(publiKey, token)
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed.")
	}

	return role, nil
}

// ParseTokenFromContext 从 gin.Context 中解析 JWT Token
func ParseTokenFromHTTPContext(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// Authorization header not found.
		return "", errors.New("Authorization header not found.")
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		// Incorrect Authorization header format.
		return "", errors.New("Incorrect Authorization header format.")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		// JWT not found.
		return "", errors.New("JWT not found.")
	}

	return token, nil
}

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

	return "", errors.WithCode(code.ErrInvalidAuthHeader, "missing app or token in metadata")
}

// parseRoleFormToken 从 JWT Token 中解析角色
func parseRoleFormToken(publiKey string, tokenString string) (string, error) {
	claims := CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method.")
		}
		// parse public key
		pubKey, err := jwt.ParseECPublicKeyFromPEM([]byte(publiKey))
		if err != nil {
			return nil, errors.Wrapf(err, "Parse public key failed.")
		}
		// return public key for verification
		return pubKey, nil
	})
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed.")
	}

	if !token.Valid {
		return "", errors.New("token is invalid")
	}

	role := claims.GetRole()
	return role, nil
}
