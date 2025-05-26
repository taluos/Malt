package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInfo struct {
	Realm   string        `json:"realm"`
	Timeout time.Duration `json:"timeout"`

	token         string
	signingMethod jwt.SigningMethod
	privateKey    string
	keyFunc       jwt.Keyfunc

	// RBAC fields
	userID     string
	fullMethod string
	role       string
}

func NewJwtInfo(privateKey, userID, fullMethod, role string, opts ...JWTOption) (*JwtInfo, error) {
	jwtInfo := &JwtInfo{
		Realm:         "JWT",
		Timeout:       DefaultExpireTime,
		signingMethod: jwt.SigningMethodES256, // 统一使用ECDSA
		privateKey:    privateKey,
		userID:        userID,
		fullMethod:    fullMethod,
		role:          role,
	}

	// 设置默认的密钥函数
	jwtInfo.keyFunc = func(token *jwt.Token) (any, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(privateKey))
	}

	for _, opt := range opts {
		opt(jwtInfo)
	}

	var err error
	jwtInfo.token, err = GenerateJWT(*jwtInfo)
	return jwtInfo, err
}

type JWTOption func(*JwtInfo)

func WithRealm(realm string) JWTOption {
	return func(o *JwtInfo) {
		o.Realm = realm
	}
}

func WithTimeout(timeout time.Duration) JWTOption {
	return func(o *JwtInfo) {
		o.Timeout = timeout
	}
}

func WithKeyFunc(keyFunc jwt.Keyfunc) JWTOption {
	return func(o *JwtInfo) {
		o.keyFunc = keyFunc
	}
}

func WithSigningMethod(method jwt.SigningMethod) JWTOption {
	return func(o *JwtInfo) {
		o.signingMethod = method
	}
}

// Getter methods
func (j *JwtInfo) Keyfunc() jwt.Keyfunc {
	return j.keyFunc
}

func (j *JwtInfo) Token() string {
	return j.token
}

func (j *JwtInfo) GetUserID() string {
	return j.userID
}

func (j *JwtInfo) GetFullMethod() string {
	return j.fullMethod
}

func (j *JwtInfo) GetRole() string {
	return j.role
}

func (j *JwtInfo) GetSigningMethod() jwt.SigningMethod {
	return j.signingMethod
}
