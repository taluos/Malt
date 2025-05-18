package auth

import (
	"strings"
	"time"

	"github.com/taluos/Malt/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT 生成一个有效期为 tokenExpiretime 的 JWT Token
func GenerateJWT(PrivateKey string, userID string, name string, role string, expireTime time.Duration) (string, error) {

	if PrivateKey == "" || userID == "" || name == "" || role == "" {
		return "", errors.New("invalid input.")
	}

	// 1. 解析私钥
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(PrivateKey))
	if err != nil {
		return "", errors.Wrapf(err, "parse private key failed.")
	}

	// 2. 创建 claims
	claims := NewCustomClaims(userID, name, role, expireTime)

	// 3. 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// 4. 签名并返回 token 字符串
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrapf(err, "sign token failed.")
	}

	return signedToken, nil
}

func ParseTokenFromContext(c *gin.Context) (string, error) {
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

func ParseRoleFromContext(c *gin.Context, publiKey string) (string, error) {
	token, err := ParseTokenFromContext(c)
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed.")
	}

	role, err := parseRoleFormToken(publiKey, token)
	if err != nil {
		return "", errors.Wrapf(err, "parse token failed.")
	}

	return role, nil
}

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
