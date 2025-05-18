package auth

import (
	"context"
	"testing"
	"time"

	rpcmetadata "Malt/api/rpcmetadata"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticator(t *testing.T) {
	auth, err := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return []byte("dummy"), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, auth)
}

func TestAuthenticate_ValidToken(t *testing.T) {
	// 生成合法token
	token, err := GenerateJWT(testPrivateKey, "/test.Service/Method")
	assert.NoError(t, err)

	// 构造metadata
	md := map[string][]string{
		"authorization": {"Bearer " + token},
	}
	ctx := rpcmetadata.NewServerContext(context.Background(), md)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(testPubliKey))
	})

	err = auth.Authenticate(ctx, "/test.Service/Method")
	assert.NoError(t, err)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	md := map[string][]string{
		"authorization": {"Bearer invalidtoken"},
	}
	ctx := rpcmetadata.NewServerContext(context.Background(), md)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(testPubliKey))
	})

	err := auth.Authenticate(ctx, "/test.Service/Method")
	assert.Error(t, err)
}

func TestAuthenticate_MissingToken(t *testing.T) {
	ctx := context.Background()
	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(testPubliKey))
	})
	err := auth.Authenticate(ctx, "/test.Service/Method")
	assert.Error(t, err)
}

func TestValidateToken_MethodNotMatch(t *testing.T) {
	token, err := GenerateJWT(testPrivateKey, "/test.Service/OtherMethod")
	assert.NoError(t, err)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(testPubliKey))
	})

	err = auth.ValidateToken(context.Background(), token, "/test.Service/Method")
	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	// 构造过期token
	privateKey, _ := jwt.ParseECPrivateKeyFromPEM([]byte(testPrivateKey))
	claims := CustomClaims{
		MethodName: "/test.Service/Method",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-5 * time.Minute)),
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token, _ := tokenObj.SignedString(privateKey)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(testPubliKey))
	})

	err := auth.ValidateToken(context.Background(), token, "/test.Service/Method")
	assert.Error(t, err)
}
