package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT_Success(t *testing.T) {
	token, err := GenerateJWT(TestPrivateKey, "uid", "uname", "admin", time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_InvalidInput(t *testing.T) {
	_, err := GenerateJWT("", "uid", "uname", "admin", time.Minute)
	assert.Error(t, err)
	_, err = GenerateJWT(TestPrivateKey, "", "uname", "admin", time.Minute)
	assert.Error(t, err)
	_, err = GenerateJWT(TestPrivateKey, "uid", "", "admin", time.Minute)
	assert.Error(t, err)
	_, err = GenerateJWT(TestPrivateKey, "uid", "uname", "", time.Minute)
	assert.Error(t, err)
}

func TestGenerateJWT_InvalidKey(t *testing.T) {
	_, err := GenerateJWT("invalid-key", "uid", "uname", "admin", time.Minute)
	assert.Error(t, err)
}

func TestParseRoleFromContext_Success(t *testing.T) {
	token, err := GenerateJWT(TestPrivateKey, "uid", "uname", "admin", time.Minute)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	role, err := ParseRoleFromHTTPContext(c, TestPubliKey)
	assert.NoError(t, err)
	assert.Equal(t, "admin", role)
}

func TestParseRoleFromContext_HeaderMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	_, err := ParseRoleFromHTTPContext(c, TestPubliKey)
	assert.Error(t, err)
}

func TestParseRoleFromContext_HeaderFormatError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Token abcdefg")

	_, err := ParseRoleFromHTTPContext(c, TestPubliKey)
	assert.Error(t, err)
}

func TestParseRoleFromContext_TokenInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.value")

	_, err := ParseRoleFromHTTPContext(c, TestPubliKey)
	assert.Error(t, err)
}

func TestParseRoleFromContext_TokenExpired(t *testing.T) {
	token, err := GenerateJWT(TestPrivateKey, "uid", "uname", "admin", time.Second)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	_, err = ParseRoleFromHTTPContext(c, TestPubliKey)
	assert.Error(t, err)
}
