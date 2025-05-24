package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims

	UserID     string `json:"uid"`
	FullMethod string `json:"method"`
	Role       string `json:"role"`
}

func NewCustomClaims(userID, fullMethod, role string, expireTime time.Duration, refreshTime time.Duration) *CustomClaims {
	now := time.Now()
	customClaims := &CustomClaims{
		UserID:     userID,
		FullMethod: fullMethod,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(DefaultExpireTime)),
			NotBefore: jwt.NewNumericDate(now.Add(-DefaultMaxRefresh)),
		},
	}

	if expireTime != 0 {
		customClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(expireTime))
	}
	if refreshTime != 0 {
		customClaims.RegisteredClaims.NotBefore = jwt.NewNumericDate(now.Add(-refreshTime))
	}

	return customClaims
}

type ClaimsGetter interface {
	GetUserID() string
	GetFullMethod() string
	GetRole() string
}

func (c *CustomClaims) GetUserID() string {
	return c.UserID
}

func (c *CustomClaims) GetFullMethod() string {
	return c.FullMethod
}

func (c *CustomClaims) GetRole() string {
	return c.Role
}
