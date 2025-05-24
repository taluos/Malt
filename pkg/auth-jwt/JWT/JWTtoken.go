package jwt

import (
	"github.com/taluos/Malt/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT 生成一个有效期为 tokenExpiretime 的 JWT Token
// 其中 userID 是用户的唯一标识符，FullMethod 是用户调用的方法名称，role 是用户的角色
// expireTime 是 token 的有效期

func GenerateJWT(jwtInfo JwtInfo) (string, error) {
	if jwtInfo.privateKey == "" || jwtInfo.userID == "" || jwtInfo.fullMethod == "" || jwtInfo.role == "" {
		return "", errors.New("invalid input: missing required fields")
	}

	// 解析私钥
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(jwtInfo.privateKey))
	if err != nil {
		return "", errors.Wrapf(err, "parse private key failed")
	}

	// 创建claims
	claims := NewCustomClaims(jwtInfo.userID, jwtInfo.fullMethod, jwtInfo.role, jwtInfo.Timeout, jwtInfo.MaxRefresh)

	// 创建token
	token := jwt.NewWithClaims(jwtInfo.signingMethod, claims)

	// 签名并返回token字符串
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrapf(err, "sign token failed")
	}

	return signedToken, nil
}
