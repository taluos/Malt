package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"uid"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewCustomClaims(userID, name, role string, expireTime time.Duration) *CustomClaims {
	now := time.Now()
	customClaims := &CustomClaims{
		UserID: userID,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpiretime)),
		},
	}

	if expireTime != 0 {
		customClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(expireTime))
	}

	return customClaims
}

type CustomClaimsMethod interface {
	GetUserID() string
	GetName() string
	GetRole() string
}

func (c *CustomClaims) GetUserID() string {
	return c.UserID
}

func (c *CustomClaims) GetName() string {
	return c.Name
}

func (c *CustomClaims) GetRole() string {
	return c.Role
}
