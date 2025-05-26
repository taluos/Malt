package auth

import (
	"context"
	"fmt"
	"testing"

	rpcmetadata "github.com/taluos/Malt/api/rpcmetadata"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticator(t *testing.T) {
	auth, err := NewAuthenticator(func(token *jwt.Token) (any, error) {
		return []byte("dummy"), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, auth)
}

func TestAuthenticate_ValidToken(t *testing.T) {
	// 生成 jwtInfo
	jwtInfo, err := JWT.NewJwtInfo(JWT.TestPrivateKey, "uuid", "/test.Service/Method", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, jwtInfo)
	// 生成合法token
	token, err := JWT.GenerateJWT(*jwtInfo)
	fmt.Println(token)
	assert.NoError(t, err)

	// 构造metadata
	md := map[string][]string{
		"authorization": {"Bearer " + token},
	}
	ctx := rpcmetadata.NewServerContext(context.Background(), md)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (any, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(JWT.TestPublicKey))
	})

	err = auth.RPCAuthenticate(ctx, "uuid", "/test.Service/Method", "admin")
	assert.NoError(t, err)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	md := map[string][]string{
		"authorization": {"Bearer invalidtoken"},
	}
	ctx := rpcmetadata.NewServerContext(context.Background(), md)

	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(JWT.TestPublicKey))
	})

	err := auth.RPCAuthenticate(ctx, "", "/test.Service/Method", "")
	assert.Error(t, err)
}

func TestAuthenticate_MissingToken(t *testing.T) {
	ctx := context.Background()
	auth, _ := NewAuthenticator(func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(JWT.TestPublicKey))
	})
	err := auth.RPCAuthenticate(ctx, "", "/test.Service/Method", "")
	assert.Error(t, err)
}
