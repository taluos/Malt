package auth

import (
	"net/http/httptest"
	"testing"

	Casbin "github.com/taluos/Malt/core/RBAC/Casbin"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func NewEnforcer_test(modelPath, policyPath string) *Casbin.RBACEnforcer {
	tempEnforcer, _ := Casbin.NewAdapter(nil, modelPath, policyPath)
	return tempEnforcer
}

func TestNewAuthenticator(t *testing.T) {
	enforcer := NewEnforcer_test(modelPath, policyPath)
	auth, err := NewAuthenticator(JWT.TestPubliKey, enforcer)
	assert.NoError(t, err)
	assert.NotNil(t, auth)
}

func TestAuthenticator_Authenticate_Success(t *testing.T) {
	// 生成token
	token, err := JWT.GenerateJWT(JWT.TestPrivateKey, "uid", "uname", "admin", JWT.TokenExpiretime)
	assert.NoError(t, err)

	// 构造gin.Context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/user", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	enforcer := NewEnforcer_test(modelPath, policyPath)

	auth, err := NewAuthenticator(JWT.TestPubliKey, enforcer)
	assert.NoError(t, err)

	err = auth.Authenticate(c)
	assert.NoError(t, err)
}

func TestAuthenticator_Authenticate_Fail_ParseRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/user", nil)
	// 没有设置 Authorization header
	enforcer := NewEnforcer_test(modelPath, policyPath)
	auth, err := NewAuthenticator(JWT.TestPubliKey, enforcer)
	assert.NoError(t, err)

	err = auth.Authenticate(c)
	assert.Error(t, err)
}

func TestAuthenticator_Authenticate_Fail_VerifyAuth(t *testing.T) {
	// 生成token
	token, err := JWT.GenerateJWT(JWT.TestPrivateKey, "uid", "uname", "editor", JWT.TokenExpiretime)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/api/user", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	enforcer := NewEnforcer_test(modelPath, policyPath)

	auth, err := NewAuthenticator(JWT.TestPubliKey, enforcer)
	assert.NoError(t, err)

	err = auth.Authenticate(c)
	assert.Error(t, err)
}
