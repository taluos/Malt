package auth

import (
	"time"

	"github.com/taluos/Malt/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT 生成一个有效期为 tokenExpiretime 的 JWT Token
// 其中 userID 是用户的唯一标识符，FullMethod 是用户调用的方法名称，role 是用户的角色
// expireTime 是 token 的有效期
func GenerateJWT(PrivateKey string, userID string, FullMethod string, role string, expireTime time.Duration) (string, error) {

	if PrivateKey == "" || userID == "" || FullMethod == "" || role == "" {
		return "", errors.New("invalid input.")
	}

	// 1. 解析私钥
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(PrivateKey))
	if err != nil {
		return "", errors.Wrapf(err, "parse private key failed.")
	}

	// 2. 创建 claims
	claims := NewCustomClaims(userID, FullMethod, role, expireTime)

	// 3. 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// 4. 签名并返回 token 字符串
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrapf(err, "sign token failed.")
	}

	return signedToken, nil
}
