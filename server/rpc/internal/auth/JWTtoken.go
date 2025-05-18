package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID     string `json:"uid"`
	MethodName string `json:"method"` // gRPC FullMethod，例如 /pkg.Service/Method
	jwt.RegisteredClaims
}

// GenerateJWT 生成一个有效期为 tokenExpiretime 的 JWT Token
func GenerateJWT(PrivateKey string, fullMethod string) (string, error) {
	// 1. 解析私钥
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(PrivateKey))
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %w", err)
	}

	// 2. 创建 claims
	now := time.Now()

	claims := CustomClaims{
		// UserID:     userID,
		MethodName: fullMethod,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiretime)),
		},
	}

	// 3. 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// 4. 签名并返回 token 字符串
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return signedToken, nil
}
